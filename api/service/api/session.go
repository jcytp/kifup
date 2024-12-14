// service/api/session.go

package api

import (
	"log/slog"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/jcytp/kifup-api/common/env"
	"github.com/jcytp/kifup-api/common/handler"
	"github.com/jcytp/kifup-api/service/dao"
)

type requestLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context, req requestLogin) (string, string, error) {
	account, err := dao.GetAccountByEmail(req.Email)
	if err != nil {
		return "", "UNAUTHORIZED - Invalid email or password", err
	}

	// check password
	if err := bcrypt.CompareHashAndPassword([]byte(account.PassHash), []byte(req.Password)); err != nil {
		return "", "UNAUTHORIZED - Invalid email or password", err
	}

	// generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": account.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(env.SecretKey())
	if err != nil {
		return "", "Failed to generate token", err
	}

	account.LastLoginAt = time.Now()
	if err := dao.UpdateAccount(account); err != nil {
		slog.WarnContext(c, err.Error())
		// not return error
	}

	return tokenString, "", nil
}

func RefreshSession(c *gin.Context) (string, string, error) {
	accountID := handler.GetActorID(c)

	account, err := dao.GetAccountByID(accountID)
	if err != nil {
		return "", "Failed to get Account", err
	}

	// generate JWT token
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": account.ID,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	newTokenString, err := newToken.SignedString(env.SecretKey())
	if err != nil {
		return "", "Failed to generate new token", err
	}

	return newTokenString, "", nil
}
