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
	}

	// public endpoints (optional authentication)
	rp := ra.Group("/")
	rp.Use(handler.MwCheckSession)
	{
		// ToDo: search kifu, etc.
	}

	// require session
	rs := rp.Group("/")
	rs.Use(handler.MwRequireSession)
	{
		rs.GET("/account", handler.HandlerOut(api.GetAccount))
		rs.DELETE("/account", handler.Handler(api.DeleteAccount))
		rs.PUT("/account/password", handler.HandlerIn(api.ChangePassword))
		rs.POST("/session/refresh", handler.HandlerOut(api.RefreshSession))

		// ToDo: manage kifu, etc.
	}

	r.Run() // default -> localhost:8080
}
