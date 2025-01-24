package api

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/jcytp/kifup-api/common/env"
	"github.com/jcytp/kifup-api/service/dao"
)

func SetupTables() {
	if err := dao.CreateEmailVerificationTable(); err != nil {
		log.Fatal("failed to create email verification table")
	}
	if err := dao.CreateAccountTable(); err != nil {
		log.Fatal("failed to create account table")
	}
	if err := dao.CreateKifuTable(); err != nil {
		log.Fatal("failed to create kifu table")
	}
	if err := dao.CreateKifuOptionTable(); err != nil {
		log.Fatal("failed to create kifu option table")
	}
	if err := dao.CreateKifuTagTable(); err != nil {
		log.Fatal("failed to create kifu tag table")
	}
	if err := dao.CreateKifuBranchTable(); err != nil {
		log.Fatal("failed to create kifu branch table")
	}
	if err := dao.CreateKifuMoveTable(); err != nil {
		log.Fatal("failed to create kifu move table")
	}
	if err := dao.CreateKifuLikeTable(); err != nil {
		log.Fatal("failed to create kifu like table")
	}
	if err := dao.CreateKifuCommentTable(); err != nil {
		log.Fatal("failed to create kifu comment table")
	}
}

type GetServerStatusResponse struct {
	Mode        string `json:"mode"`
	Environment string `json:"environment"`
	Status      string `json:"status"`
}

func GetServerStatus(c *gin.Context) (*GetServerStatusResponse, string, error) {
	resp := &GetServerStatusResponse{
		Mode:        gin.Mode(),
		Environment: env.Environment(),
		Status:      "running",
	}
	return resp, "", nil
}
