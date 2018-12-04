package conf

import (
	"github.com/gin-gonic/gin"
	. "github.com/sipt/shuttle/constant"
)

func APIRoute(router *gin.RouterGroup, eventChan chan *EventObj) {
	//dns
	router.GET("/dns", GetDNSConfig)
	router.POST("/dns", SetDNSConfig)

	//MitM rules
	router.GET("/mitm/rules", GetMitMRules)
	router.POST("/mitm/rules", AppendMitMRules)
	router.DELETE("/mitm/rules", DelMitMRules)

	//rule
	router.GET("/rules", GetRule)
	router.POST("/rules", SetRule)

	//proxy & proxy group
	router.GET("/proxy", GetProxy)
	router.POST("/proxy", SetProxy)
	router.GET("/proxy/group", GetProxyGroup)
	router.POST("/proxy/group", SetProxyGroup)

	//general
	router.POST("/general", SetGeneralConfig(eventChan))
	router.GET("/general", GetGeneralConfig)

	//http map
	router.POST("/http/map", SetHttpMap)
	router.GET("/http/map", GetHttpMap)
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
