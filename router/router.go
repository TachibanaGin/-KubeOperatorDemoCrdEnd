package router

import (
	"Crd-End/controller"
	"Crd-End/interceptor"
	"Crd-End/ws"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {

	router := gin.Default()

	router.POST("/Crd/Login", controller.Login)
	v1 := router.Group("/Crd/Get")
	{
		v1.GET("/getFront", controller.GetFront)
		//v1.GET("/getFrontList", controller.GetFrontLis)
		v1.GET("/execPod", ws.WsHandler)
	}

	v2 := router.Group("/Crd/Post")
	{
		v2.Use(interceptor.Interceptor())
		v2.POST("/addFront", controller.AddFront)
		v2.POST("/getFrontList", controller.GetFrontList)
		v2.POST("/delFront", controller.DelFront)
		v2.POST("/updateFront", controller.UpdateFront)
		v2.POST("/getFrontOne", controller.GetFrontOne)
		v2.POST("/getPodDetail", controller.GetPodDetail)
		v2.POST("/getDescribePod", controller.GetDescribePod)
	}

	return router
}