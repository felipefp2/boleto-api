package api

import "github.com/gin-gonic/gin"

//InstallV1 instala a api versao 1
func InstallV1(router *gin.Engine) {
	v1 := router.Group("v1")
	v1.Use(timingMetrics())
	v1.Use(ReturnHeaders())
	v1.POST("/boleto/register", ParseBoleto(), registerBoleto)
	v1.POST("/boleto/edit", ParseBoleto(), editBoleto)
	v1.GET("/boleto/:id/:pk", getBoletoByID)
}
