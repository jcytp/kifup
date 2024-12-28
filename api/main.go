package main

import (
	"log/slog"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/jcytp/kifup-api/common/db"
	"github.com/jcytp/kifup-api/common/handler"
	"github.com/jcytp/kifup-api/common/log"
	"github.com/jcytp/kifup-api/service/api"
)

func main() {
	log.Setup(slog.LevelDebug)

	if db.CheckDBFileExists() {
		db.New()
	} else {
		db.New()
		api.SetupTables()
	}
	defer db.Close()

	r := gin.Default()
	if gin.Mode() == gin.DebugMode {
		config := cors.DefaultConfig()
		config.AllowOrigins = []string{"http://192.168.11.12:8081"} // frontend開発サーバー
		config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
		config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
		r.Use(cors.New(config))

		r.Static("/swagger", "./swagger")
	}

	// public endpoints (no authentication)
	rPub := r.Group("/api")
	rPub.Use(handler.MwStorePathIDs)

	// public endpoints (optional authentication)
	rOpt := rPub.Group("/")
	rOpt.Use(handler.MwCheckSession)

	// require session
	rSes := rOpt.Group("/")
	rSes.Use(handler.MwRequireSession)

	// account api
	rPub.POST("/account", handler.HandlerIn(api.CreateAccount))
	rSes.GET("/account", handler.HandlerOut(api.GetAccount))
	rSes.DELETE("/account", handler.Handler(api.DeleteAccount))
	rSes.PUT("/account/password", handler.HandlerIn(api.ChangePassword))
	rSes.PUT("/account/info", handler.HandlerIn(api.UpdateAccountInfo))
	rPub.GET("/account/:accountID", handler.HandlerOut(api.GetAccountByID))

	// session api
	rPub.POST("/session/login", handler.HandlerInOut(api.Login))
	rSes.POST("/session/refresh", handler.HandlerOut(api.RefreshSession))

	// kifu api
	rSes.POST("/kifu", handler.HandlerInOut(api.CreateKifu))
	rOpt.GET("/kifu", handler.HandlerInPagination(api.ListKifus))
	rOpt.GET("/kifu/:kifuID", handler.HandlerOut(api.GetKifu))
	rSes.DELETE("/kifu/:kifuID", handler.Handler(api.DeleteKifu))
	rSes.PUT("/kifu/:kifuID", handler.HandlerIn(api.UpdateKifuInfo))
	rSes.PUT("/kifu/:kifuID/moves", handler.HandlerIn(api.UpdateKifuMoves))

	r.Run() // default -> localhost:8080
}
