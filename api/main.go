package main

import (
	"log/slog"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/jcytp/kifup-api/common/db"
	"github.com/jcytp/kifup-api/common/env"
	"github.com/jcytp/kifup-api/common/handler"
	"github.com/jcytp/kifup-api/common/log"
	"github.com/jcytp/kifup-api/service/api"
)

func main() {
	log.Setup(slog.LevelDebug)

	// environment values
	env.Initialize()

	// dabase setup
	if db.CheckDBFileExists() {
		db.New() // 既存のDBの使用
	} else {
		if err := db.DownloadDB(); err != nil {
			// ダウンロードできなければ新規作成して使用
			db.New()
			api.SetupTables()
			db.UploadDB() // 作成したDBをS3にアップロード
		} else {
			db.New() // ダウンロードしたDBの使用
		}
	}
	defer db.Close()
	db.StartBackupCycle()    // 定期的なバックアップの作成
	db.ScheduleFinalBackup() // 正常終了時の最終バックアップ

	// gin engine
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = env.AllowedOrigins()
	config.AllowMethods = []string{"GET", "HEAD", "OPTIONS", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(config))
	r.Use(handler.MwSetCacheControl)

	// swagger route
	if env.SwaggerEnable() {
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

	// health-check api
	rPub.GET("/status", handler.HandlerOut(api.GetServerStatus))

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

	r.Run(":80") // default -> localhost:8080
}
