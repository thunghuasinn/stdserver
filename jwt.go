package stdserver

import (
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const JwtContextKey = "jwt"

type LoginFunc func(c *fiber.Ctx) (jwt.Claims, error)

func defaultLoginHandler(c *fiber.Ctx) (jwt.Claims, error) {
	iss := defaultAppName
	cfg, ok := c.Locals("config").(*Settings)
	if ok {
		iss = cfg.Name
	}
	id, err := uuid.NewRandom()
	if err != nil {
		id = uuid.Must(uuid.FromBytes(make([]byte, 16)))
	}
	return &jwt.RegisteredClaims{
		Audience:  jwt.ClaimStrings{"dev"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		ID:        id.String(),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    iss,
		Subject:   "anonymous",
	}, nil
}

func JWT(cfg *Settings, claimsType jwt.Claims) fiber.Handler {
	logger := cfg.Logger.With().Str("module", "core").Str("subModule", "jwt").Logger()
	defer func() {
		if r := recover(); r != nil {
			logger.Fatal().Msgf("panic: %+v", r)
		}
	}()
	if cfg.LoginHandler == nil {
		cfg.LoginHandler = defaultLoginHandler
	}

	kt, err := LoadKeyTableFromDir(cfg.KeyTableDir)
	if err != nil {
		logger.Fatal().Err(err).Msg("while loading key table")
	}
	signMap := kt.GetPrivateKeys()
	ware := jwtware.New(jwtware.Config{
		SigningKeys:   kt.GetPublicKeys(),
		SigningMethod: "ES256",
		ContextKey:    JwtContextKey,
		Claims:        claimsType,
		Filter:        cfg.SkipAuth,
	})
	return func(c *fiber.Ctx) error {
		if c.Method() == fiber.MethodPost && c.Path() == cfg.LoginPath {
			claims, err := cfg.LoginHandler(c)
			if err != nil {
				return err
			}
			token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
			kid := ""
			for kid = range signMap {
				break
			}
			token.Header["kid"] = kid
			t, err := token.SignedString(signMap[kid])
			if err != nil {
				return fiber.ErrInternalServerError
			}

			return c.JSON(fiber.Map{"data": fiber.Map{"token": t}})
		} else {
			return ware(c)
		}
	}
}
