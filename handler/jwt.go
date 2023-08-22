package handler

import (
	"crypto/rsa"
	"fmt"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	defaultPrivateKey = "../cert/id_rsa"
	privKey           *rsa.PrivateKey
)

func init() {
	if os.Getenv("PRIVATE_KEY") != "" {
		defaultPrivateKey = os.Getenv("PRIVATE_KEY")
	}
	p, err := os.ReadFile(defaultPrivateKey)
	if nil != err {
		fmt.Printf("failed reading private key. stack trace: %s\n", err)
	}

	privKey, err = jwt.ParseRSAPrivateKeyFromPEM(p)
	if nil != err {
		fmt.Printf("failed parsing private key. stack trace: %s\n", err)
	}
}

func Middleware() echo.MiddlewareFunc {
	var skippers = []string{
		"/login",
		"/register",
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// check skipper path
			for _, skip := range skippers {
				if skip == c.Path() {
					return next(c)
				}
			}

			token := strings.Trim(c.Request().Header.Get("Authorization"), "Bearer ")
			if token == "" {
				return c.JSON(http.StatusForbidden, generated.ErrorResponse{
					Message: "missing or malformed jwt",
				})
			}

			parser, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return &privKey.PublicKey, nil
			})
			if nil != err {
				return c.JSON(http.StatusForbidden, "invalid or expired jwt")
			}

			claims := parser.Claims.(jwt.MapClaims)
			c.Set("user", claims["dat"])

			return next(c)
		}
	}
}

func Create(content any) (string, error) {
	// get token expiry from environment
	// default we will set token for one hour
	expiry, err := time.ParseDuration(os.Getenv("TTL"))
	if nil != err {
		expiry = time.Hour
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["dat"] = content
	claims["exp"] = now.Add(expiry).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	t, err := token.SignedString(privKey)
	if nil != err {
		return "", err
	}

	return t, nil
}
