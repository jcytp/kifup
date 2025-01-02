package api

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/jcytp/kifup-api/service/dao"
)

func SetupTables() {
	if err := dao.CreateAccountTable(); err != nil {
		log.Fatal()
	}
	if err := dao.CreateKifuTable(); err != nil {
		log.Fatal()
	}
	if err := dao.CreateKifuOptionTable(); err != nil {
		log.Fatal()
	}
	if err := dao.CreateKifuTagTable(); err != nil {
		log.Fatal()
	}
	if err := dao.CreateKifuBranchTable(); err != nil {
		log.Fatal()
	}
	if err := dao.CreateKifuMoveTable(); err != nil {
		log.Fatal()
	}
}

type GetServerStatusResponse struct {
	Mode   string `json:"mode"`
	Status string `json:"status"`
}

func GetServerStatus(c *gin.Context) (*GetServerStatusResponse, string, error) {
	resp := &GetServerStatusResponse{
		Mode:   "debug",
		Status: "running",
	}
	return resp, "", nil
}
