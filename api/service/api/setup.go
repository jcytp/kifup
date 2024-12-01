package api

import (
	"log"

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
