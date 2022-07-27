package main

import (
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/util"
)

// 与填写的服务器配置中的Token一致
const Token = "clear_love"

func main() {
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}
	router := gin.Default()
	router.GET("/wx", WXCheckSignature)
	router.POST("/wechat/msg", WxMsgReceive)

	log.Fatalln(router.Run(":80"))
}

func WXCheckSignature(ctx *gin.Context) {
	signature := ctx.Query("signature")
	timestamp := ctx.Query("timestamp")
	nonce := ctx.Query("nonce")
	echostr := ctx.Query("echostr")
	ok := util.CheckSignature(signature, timestamp, nonce, Token)
	if !ok {
		log.Println("微信公众号接入校验失败!")
		return
	}

	log.Println("微信公众号接入校验成功!")
	_, _ = ctx.Writer.WriteString(echostr)
}

// WxMsgReceive 被动回复用户消息
func WxMsgReceive(ctx *gin.Context) {
	var textMsg WxTextMsg
	err := ctx.ShouldBindXML(&textMsg)
	if err == nil {
		log.Printf("[消息接收] - XML数据包解析失败: %v\n", err)
		return
	}
	log.Printf("[消息接收] - 收到消息, 消息类型为: %s, 消息内容为: %s\n", textMsg.MsgType, textMsg.Content)
	// 对接收的消息进行被动回复
	WXMsgReply(ctx, textMsg.ToUserName, textMsg.FromUserName)
}

type WxTextMsg struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string
	MsgId        int64
}

// WXRepTextMsg 微信回复文本消息结构体
type WXRepTextMsg struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string
	// 若不标记XMLName, 则解析后的xml名为该结构体的名称
	XMLName xml.Name `xml:"xml"`
}

// WXMsgReply 微信消息回复
func WXMsgReply(c *gin.Context, fromUser, toUser string) {
	repTextMsg := WXRepTextMsg{
		ToUserName:   toUser,
		FromUserName: fromUser,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      fmt.Sprintf("[消息回复] - %s--%s", time.Now().Format("2006-01-02 15:04:05"), "签到成功"),
	}

	msg, err := xml.Marshal(&repTextMsg)
	if err != nil {
		log.Printf("[消息回复] - 将对象进行XML编码出错: %v\n", err)
		return
	}
	_, _ = c.Writer.Write(msg)
}
