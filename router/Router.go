package router

import (
	"git.in.codoon.com/Overseas/runbox/first-test/controller"
	"github.com/gin-gonic/gin"
	"log"
)

func Router(engine *gin.Engine) {
	log.Printf("StartHttpServer ")

	eg_sports := engine.Group("/v1/")
	{
		eg_sports.POST("/v1/gps", controller.PostGps)
		eg_sports.POST("/v1/gps", controller.PostGps)

	}

}
