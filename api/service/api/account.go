// service/api/account.go

package api

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/jcytp/kifup-api/common/aws"
	"github.com/jcytp/kifup-api/common/env"
	"github.com/jcytp/kifup-api/common/handler"
	"github.com/jcytp/kifup-api/service/dao"
	"github.com/jcytp/kifup-api/service/model"
)

type requestSendVerificationEmail struct {
	Email string `json:"email" binding:"required,email,max=255"`
}

func SendVerificationEmail(c *gin.Context, req requestSendVerificationEmail) (string, error) {
	ev := model.NewEmailVerification(req.Email)
	if err := dao.InsertEmailVerification(ev); err != nil {
		return "Failed to create verification", err
	}

	subject := "棋譜UP メール認証"
	message := fmt.Sprintf("認証コードを入力して、棋譜UPへの登録を完了してください。\n\n%s\n（有効期限：%s）", ev.VerificationCode, ev.ExpiredAt.Format("2006/01/02 15:04"))
	aws.SesSendEmailOne(env.EmailSender(), ev.Email, subject, message)
	return "", nil
}

type requestCheckVerificationCode struct {
	Email string `json:"email" binding:"required,email,max=255"`
	Code  string `json:"code" binding:"required,len=6"`
}

func CheckVerificationCode(c *gin.Context, req requestCheckVerificationCode) (string, error) {
	ev, err := dao.GetEmailVerification(req.Email)
	if err != nil {
		return "Failed to get verification", err
	}
	if !ev.IsValid(req.Code) {
		return "Invalid verification code", fmt.Errorf("invalid verification code")
	}
	return "", nil
}

type requestCreateAccount struct {
	Name     string `json:"name" binding:"required,min=2,max=60"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=6"`
	Code     string `json:"code" binding:"required,len=6"`
}

func CreateAccount(c *gin.Context, req requestCreateAccount) (string, error) {
	if env.IsProduction() { // Alpha version
		return "", fmt.Errorf("not allowed operation in the alpha version")
	}
	if !env.IsDevelopment() { // 開発環境ではメール認証なしでユーザー登録
		ev, err := dao.GetEmailVerification(req.Email)
		if err != nil {
			return "Failed to get verification", err
		}
		if !ev.IsValid(req.Code) {
			return "Invalid verification code", fmt.Errorf("invalid verification code")
		}
	}

	account_check, err := dao.GetAccountByEmail(req.Email)
	if err != nil && err != sql.ErrNoRows {
		return "Failed to check existing account", err
	}
	if account_check != nil {
		return "Email already registered", fmt.Errorf("email already registered")
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "Failed to process password", err
	}
	account := &model.Account{
		Name:     req.Name,
		Email:    req.Email,
		PassHash: string(passHash),
	}

	if _, err = dao.InsertAccount(account); err != nil {
		return "Failed to create account", err
	}

	if !env.IsDevelopment() {
		if err := dao.MarkEmailVerificationAsUsed(req.Email); err != nil {
			return "Failed to mark verification as used", err
		}
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

	if err := dao.DeleteAccount(aid); err != nil {
		return "Failed to delete account", err
	}

	return "", nil
}

type requestResetPassword struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=6"`
	Code     string `json:"code" binding:"required,len=6"`
}

func ResetPassword(c *gin.Context, req requestResetPassword) (string, error) {
	if !env.IsDevelopment() { // 開発環境ではメール認証なしでパスワード変更
		ev, err := dao.GetEmailVerification(req.Email)
		if err != nil {
			return "Failed to get verification", err
		}
		if !ev.IsValid(req.Code) {
			return "Invalid verification code", fmt.Errorf("invalid verification code")
		}
	}

	account, err := dao.GetAccountByEmail(req.Email)
	if err != nil {
		return "Failed to get account", err
	}
	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "Failed to process password", err
	}
	account.PassHash = string(passHash)
	if err = dao.UpdateAccount(account); err != nil {
		return "Failed to update password", err
	}

	if !env.IsDevelopment() {
		if err := dao.MarkEmailVerificationAsUsed(req.Email); err != nil {
			return "Failed to mark verification as used", err
		}
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
	if err = dao.UpdateAccount(account); err != nil {
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
	if err = dao.UpdateAccount(account); err != nil {
		return "Failed to update account info", err
	}

	return "", nil
}
