// common/handler/middleware.go

package handler

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/jcytp/kifup-api/common/env"
	"github.com/jcytp/kifup-api/common/log"
)

func MwCheckSession(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.Next() // no authorization header
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return env.SecretKey(), nil
	})
	if err != nil {
		slog.WarnContext(c.Request.Context(), "Invalid token", "error", "cannot parse authorization header to jwt token")
		c.Next()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		slog.WarnContext(c.Request.Context(), "Invalid token", "error", "token is not valid or claims are incorrect")
		c.Next()
		return
	}
	if err := claims.Valid(); err != nil {
		slog.WarnContext(c.Request.Context(), "Token validation failed", "error", err)
		c.Next()
		return
	}

	actorID, ok := claims["sub"].(string)
	if !ok {
		slog.WarnContext(c.Request.Context(), "Invalid token claims", "error", "sub claim is missing or not a string")
		c.Next()
		return
	}

	ctx := context.WithValue(c.Request.Context(), log.ActorIDLogKey, actorID)
	c.Request = c.Request.WithContext(ctx)
	c.Set("actorID", actorID)
	println("actorID", actorID)
	c.Next()
}

func GetActorID(c *gin.Context) string {
	return c.GetString("actorID")
}

func MwRequireSession(c *gin.Context) {
	aid := GetActorID(c)
	if aid == "" {
		ResponseUnauthorized(c, "actorID is missing", nil)
		c.Abort()
		return
	}
}

// ----------------------------------------

func MwStorePathIDs(c *gin.Context) {
	params := c.Params
	for _, param := range params {
		c.Set(param.Key, param.Value)
	}
	c.Next()
}
