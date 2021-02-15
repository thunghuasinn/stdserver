package stdserver

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/google/uuid"
)

const JwtContextKey = "jwt"

type LoginFunc func(c *fiber.Ctx) (jwt.Claims, error)

func defaultLoginHandler(c *fiber.Ctx) (jwt.Claims, error) {
	iss := "stdserver"
	cfg, ok := c.Locals("config").(*Settings)
	if ok {
		iss = cfg.Name
	}
	id, err := uuid.NewRandom()
	if err != nil {
		id = uuid.Must(uuid.FromBytes(make([]byte, 16)))
	}
	return &jwt.StandardClaims{
		Audience:  "dev",
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
		Id:        id.String(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    iss,
		Subject:   "anonymous",
	}, nil
}

func JWT(cfg *Settings, claimsType jwt.Claims) fiber.Handler {
	logger := cfg.Logger.WithField("module", "JWT")
	defer func() {
		if r := recover(); r != nil {
			logger.WithField("panic", r).Fatal("panic")
		}
	}()
	if cfg.LoginHandler == nil {
		cfg.LoginHandler = defaultLoginHandler
	}

	kt, err := LoadKeyTableFromDir(cfg.KeyTableDir)
	if err != nil {
		logger.WithError(err).Fatal("while loading key table")
	}
	signMap := kt.GetPrivateKeys()
	ware := jwtware.New(jwtware.Config{
		SigningKeys:   kt.GetPublicKeys(),
		SigningMethod: "ES256",
		ContextKey:    "jwt",
		Claims:        claimsType,
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
