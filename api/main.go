package main

import (
	"log/slog"
	"net/http"

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
		rs.DELETE("/account", handler.Handler(api.DeleteAccount))
		rs.PUT("/account/password", handler.HandlerIn(api.ChangePassword))
		rs.POST("/session/refresh", handler.HandlerOut(api.RefreshSession))

		// ToDo: manage kifu, manage account, etc.
	}

	r.Run() // default -> localhost:8080
}
