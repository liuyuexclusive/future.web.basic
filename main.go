package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/liuyuexclusive/utils/web"

	_ "github.com/liuyuexclusive/future.web.basic/docs"
	"github.com/liuyuexclusive/future.web.basic/handler"
)

type start struct{}

func (s *start) Start(engine *gin.Engine) {
	basic := engine.Group("/basic")
	authorized := basic.Group("/")
	authorized.Use(handler.Validate())
	{
		authorized.POST("/role/add", handler.RoleAddOrUpdate)

		authorized.GET("/get_info", handler.CurrentUser)

		authorized.GET("/message/count", handler.MessageCount)

		authorized.GET("/message/init", handler.MessageInit)

		authorized.GET("/message/content", handler.MessageContent)

		authorized.POST("/message/has_read", handler.HasRead)

		authorized.POST("/message/remove_readed", handler.RemoveReaded)

		authorized.POST("/message/restore", handler.Restore)
	}
	basic.POST("/save_error_logger", handler.AddErrorLog)

	basic.POST("/login", handler.Login)

	basic.POST("/logout", handler.Logout)

	// basic.GET("/metrics", func(c *gin.Context) {
	// 	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	// })

	basic.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "test ok")
	})

	// engine.LoadHTMLGlob("dist/*")
	// basic.GET("/test2", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.html", nil)
	// })
}

// @title Future对外开放API
// @version 1.0
// @description XXXXXXXXXXX
// @host
// @BasePath
func main() {
	if err := web.Startup("go.micro.api.basic", new(start), func(options *web.Options) {
		options.IsLogToES = false
		options.IsTrace = false
		options.IsMonitor = false
		options.IsRateLimite = false
	}); err != nil {
		logrus.Fatal(err)
	}
}
