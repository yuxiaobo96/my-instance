package controller

import "github.com/julienschmidt/httprouter"

func InitRouter(router *httprouter.Router) {
	// 健康检测
	//router.GET("/healthy", HealthCheck)
}
