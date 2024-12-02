package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"ppt/login/db"
	"time"
)

const secret = "ppt&w4td%vw*er3r4tfd324sde"

var ctx = context.Background()

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.PostForm("token")
		if token == "" {
			c.JSON(200, gin.H{
				"code":     401,
				"response": "should sync token",
			})
			c.Abort()
			return
		}
		//key := FormatTokenKey(name)
		custClaim, err := ParseToken(token)
		if err != nil {
			c.JSON(200, gin.H{
				"code":     401,
				"response": "invalid token",
			})
			c.Abort()
			return
		}
		if !db.IsTokenOutOfDate(custClaim.ID) {
			c.JSON(200, gin.H{
				"code":     401,
				"response": "token out of date",
			})
			c.Abort()
			return
		}
		//UpdateToken(c, token)
	}
}

func FormatTokenKey(name string) string {
	return fmt.Sprintf("token_%s", name)
}

func UpdateToken(c *gin.Context, token string) {
	db.GetRedis().Set(ctx, token, token, 24*time.Hour)
}

type JwtCustomClaims struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func (j *JwtCustomClaims) GetAudience() (jwt.ClaimStrings, error) {
	return j.Audience, nil
}

func (j *JwtCustomClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return j.ExpiresAt, nil
}

func (j *JwtCustomClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return j.IssuedAt, nil
}

func (j *JwtCustomClaims) GetIssuer() (string, error) {
	return j.Issuer, nil
}

func (j *JwtCustomClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return j.NotBefore, nil
}

func (j *JwtCustomClaims) GetSubject() (string, error) {
	return j.Subject, nil
}

func GenerateToken(id int64, name string) (string, error) {
	jwtClaims := JwtCustomClaims{
		ID:   id,
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "ppt",
			Subject:   "Token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtClaims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString string) (*JwtCustomClaims, error) {
	jwtClaims := &JwtCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, jwtClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil && !token.Valid {
		err = errors.New("invalid token")
	}
	return jwtClaims, err
}
