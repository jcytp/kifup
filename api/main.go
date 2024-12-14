package main

import (
	"log/slog"
	"net/http"

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
	ra := r.Group("/api")
	ra.Use(handler.MwStorePathIDs)
	{
		ra.GET("/test", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"result": "OK"}) })
		ra.POST("/account", handler.HandlerIn(api.CreateAccount))
		ra.POST("/session/login", handler.HandlerInOut(api.Login))

		ra.GET("/account/:accountID", handler.HandlerOut(api.GetAccountByID))
	}

	// public endpoints (optional authentication)
	rp := ra.Group("/")
	rp.Use(handler.MwCheckSession)
	{
		rp.GET("/kifu", handler.HandlerInPagination(api.ListKifus))
		rp.GET("/kifu/:kifuID", handler.HandlerOut(api.GetKifu))
	}

	// require session
	rs := ra.Group("/")
	rs.Use(handler.MwRequireSession)
	{
		rs.GET("/account", handler.HandlerOut(api.GetAccount))
		rs.DELETE("/account", handler.Handler(api.DeleteAccount))
		rs.PUT("/account/password", handler.HandlerIn(api.ChangePassword))
		rs.POST("/session/refresh", handler.HandlerOut(api.RefreshSession))

		rs.POST("/kifu", handler.HandlerInOut(api.CreateKifu))
		rs.DELETE("/kifu/:kifuID", handler.Handler(api.DeleteKifu))
		rs.PUT("/kifu/:kifuID", handler.HandlerIn(api.UpdateKifuInfo))
		rs.PUT("/kifu/:kifuID/moves", handler.HandlerIn(api.UpdateKifuMoves))
	}

	r.Run() // default -> localhost:8080
}
