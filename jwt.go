package stdserver

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

/*
var ErrUnauthenticated = fiber.NewError(fiber.StatusUnauthorized, "Unable to authenticate user.")
var ErrUnauthorized = fiber.NewError(fiber.StatusForbidden, "User unauthorized.")
*/
type LoginFunc func(c *fiber.Ctx) (string, error)

func defaultLoginHandler(c *fiber.Ctx) (string, error) {
	return "anonymous", nil
}

func JWT(cfg *Settings) fiber.Handler {
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
	})
	return func(c *fiber.Ctx) error {
		if c.Method() == fiber.MethodPost && c.Path() == cfg.LoginPath {
			sub, err := cfg.LoginHandler(c)
			if err != nil {
				return err
			}
			claims := jwt.StandardClaims{
				Subject:   sub,
				ExpiresAt: time.Now().Add(time.Hour).Unix(),
				Issuer:    "tabApi",
				IssuedAt:  time.Now().Unix(),
			}
			token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
			kid := ""
			for kid = range signMap {
				break
			}
			token.Header["kid"] = kid
			t, err := token.SignedString(signMap[kid])
			if err != nil {
				fmt.Println(err)
				return fiber.ErrInternalServerError
			}

			return c.JSON(fiber.Map{"data": fiber.Map{"token": t}})
		} else {
			return ware(c)
		}
	}
}
