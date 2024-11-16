// common/handler/middleware.go

package handler

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
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

	accountID, ok := claims["sub"].(string)
	if !ok {
		slog.WarnContext(c.Request.Context(), "Invalid token claims", "error", "sub claim is missing or not a string")
		c.Next()
		return
	}

	ctx := context.WithValue(c.Request.Context(), log.AccountIDLogKey, accountID)
	c.Request = c.Request.WithContext(ctx)
	c.Set("accountID", accountID)
	c.Next()
}

func GetAccountID(c *gin.Context) string {
	return c.GetString("accountID")
}

func MwRequireSession(c *gin.Context) {
	aid := GetAccountID(c)
	if aid == "" {
		ResponseUnauthorized(c, "accountID is missing", nil)
		c.Abort()
		return
	}
}

// ----------------------------------------

func MwStorePathIDs(c *gin.Context) {
	params := c.Params
	for _, param := range params {
		value, err := strconv.ParseInt(param.Value, 10, 64)
		if err != nil {
			msg := fmt.Sprintf("Failed to parse parameter %s", param.Key)
			ResponseBadRequest(c, msg, err)
			c.Abort()
			return
		}
		c.Set(param.Key, value)
	}
	c.Next()
}
