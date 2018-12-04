package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/sipt/shuttle/constant"
	"github.com/sipt/shuttle/controller/api"
	"github.com/sipt/shuttle/controller/api/conf"
	"github.com/sipt/shuttle/controller/web"
	"github.com/sipt/shuttle/log"
	"net"
	"net/http"
	"strings"
)

var server *http.Server

type IControllerConfig interface {
	GetControllerInterface() string
	SetControllerInterface(string)
	GetControllerPort() string
	SetControllerPort(string)
	GetLogLevel() string
}

func StartController(config IControllerConfig, eventChan chan *EventObj) {
	//if level == "info" {
	//gin.SetMode(gin.ReleaseMode)
	//gin.DefaultWriter = ioutil.Discard
	//}
	e := gin.Default()
	e.Use(Cors())
	api.APIRoute(e.Group("/api"), eventChan)
	conf.APIRoute(e.Group("/api/config"), eventChan)
	web.WebRoute(e)
	server = &http.Server{
		Addr:    net.JoinHostPort(config.GetControllerInterface(), config.GetControllerPort()),
		Handler: e,
	}
	log.Logger.Infof("[Controller] listen to:%s", server.Addr)
	server.ListenAndServe()
}

func ShutdownController() {
	s := server
	server = nil
	if s == nil {
		return
	}
	s.RegisterOnShutdown(func() {
		log.Logger.Infof("Stopped Controller goroutine...")
	})
	go func() {
		ctx := context.Background()
		s.Shutdown(ctx)
	}()
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		origin := c.Request.Header.Get("Origin")
		var headerKeys []string
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			//下面的都是乱添加的-_-~
			// c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", headerStr)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}

		c.Next()
	}
}
