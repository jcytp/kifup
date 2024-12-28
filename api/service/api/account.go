// service/api/account.go

package api

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/jcytp/kifup-api/common/handler"
	"github.com/jcytp/kifup-api/service/dao"
	"github.com/jcytp/kifup-api/service/model"
)

type requestCreateAccount struct {
	Name     string `json:"name" binding:"required,min=2,max=60"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=6"`
}

func CreateAccount(c *gin.Context, req requestCreateAccount) (string, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "Failed to process password", err
	}
	account := &model.Account{
		Name:     req.Name,
		Email:    req.Email,
		PassHash: string(passHash),
	}

	_, err = dao.InsertAccount(account)
	if err != nil {
		return "Failed to create account", err
	}

	return "", nil
}

func GetAccount(c *gin.Context) (*model.AccountResponse, string, error) {
	aid := handler.GetActorID(c)
	account, err := dao.GetAccountByID(aid)
	if err != nil {
		return nil, "Failed to get account", err
	}
	return account.ToResponse(), "", nil
}

func GetAccountByID(c *gin.Context) (*model.AccountResponse, string, error) {
	aid := c.GetString("accountID")
	account, err := dao.GetAccountByID(aid)
	if err != nil {
		return nil, "Failed to get account", err
	}
	return account.ToResponse(), "", nil
}

func DeleteAccount(c *gin.Context) (string, error) {
	aid := handler.GetActorID(c)

	err := dao.DeleteAccount(aid)
	if err != nil {
		return "Failed to delete account", err
	}

	return "", nil
}

type requestChangePassword struct {
	Password string `json:"password" binding:"required,min=6"`
}

func ChangePassword(c *gin.Context, req requestChangePassword) (string, error) {
	aid := handler.GetActorID(c)

	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "Failed to process password", err
	}

	account, err := dao.GetAccountByID(aid)
	if err != nil {
		return "Failed to get account", err
	}
	account.PassHash = string(passHash)
	err = dao.UpdateAccount(account)
	if err != nil {
		return "Failed to update password", err
	}

	return "", nil
}

type requestUpadateAccountInfo struct {
	Name         string `json:"name" binding:"required,min=2,max=60"`
	IconID       string `json:"icon_id" binding:"max=60"`
	Introduction string `json:"introduction" binding:"max=1000"`
}

func UpdateAccountInfo(c *gin.Context, req requestUpadateAccountInfo) (string, error) {
	aid := handler.GetActorID(c)
	account, err := dao.GetAccountByID(aid)
	if err != nil {
		return "Failed to get account", err
	}
	account.Name = req.Name
	account.IconID = req.IconID
	account.Introduction = req.Introduction
	err = dao.UpdateAccount(account)
	if err != nil {
		return "Failed to update account info", err
	}

	return "", nil
}
