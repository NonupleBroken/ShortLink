package web

import (
	"ShortLink/config"
	"fmt"
	"github.com/gin-gonic/gin"
)


func startHTTPServer(cfg *config.Config) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	err := InitWebLogger(&cfg.Log)
	if err != nil {
		return err
	}
	r.Use(GinLogger(), GinRecovery(true))

	indexGroup := r.Group("/")
	indexGroup.GET("", index)

	linkGroup := r.Group("/")
	linkGroup.GET("/:link_id", getShortLinkByLinkID)

	adminGroup := r.Group("/admin", adminAuth(&cfg.Server))
	adminGroup.GET("/check", checkLink)
	adminGroup.POST("/add", addLink)
	adminGroup.POST("/delete", deleteLink)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Addr, cfg.Server.Port)
	err = r.Run(addr)

	return err
}
