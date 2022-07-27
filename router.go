package main

import (
	"github.com/gin-gonic/gin"
	"wxcloudrun-golang/service"
)

func myRouter(r *gin.Engine) {
	r.GET("/wx", service.WXCheckSignature)
	r.POST("/wechat/msg", service.WxMsgReceive)
}
