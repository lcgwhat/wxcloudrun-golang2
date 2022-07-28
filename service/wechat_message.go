package service

import (
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"strings"
	"time"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/util"
)

// 与填写的服务器配置中的Token一致
const Token = "clear_love"

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
	err := ctx.BindXML(&textMsg)
	if err != nil {
		log.Printf("[消息接收] - XML数据包解析失败: %v\n", err)
		return
	}
	msg := "【签到失败】"
	if textMsg.MsgType == "text" {
		res := strings.Split(textMsg.Content, "-")
		if len(res) >= 2 {
			stringUId := res[0]
			uid, err := strconv.Atoi(stringUId)
			if err != nil {
				msg = msg + ""
			}
			note := ""
			if len(res) >= 3 {
				note = res[2]
			}

			now := time.Now()
			existCheck := dao.CheckImp.ExistCheck(uid, now)
			if existCheck == true {
				msg = msg + "今天已经签到过了。。。"
			} else {
				check := model.Check{
					UserId:    uid,
					Username:  textMsg.FromUserName,
					CheckDate: now,
					Note:      note,
				}
				err := dao.CheckImp.UpsertCheck(&check)
				if err != nil {
					msg = msg + err.Error()
				} else {
					msg = "签到成功"
				}
			}
		}
	}
	// 对接收的消息进行被动回复
	WXMsgReply(ctx, textMsg.ToUserName, textMsg.FromUserName, msg)
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
func WXMsgReply(c *gin.Context, fromUser, toUser string, msg string) {
	repTextMsg := WXRepTextMsg{
		ToUserName:   toUser,
		FromUserName: fromUser,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      msg,
	}

	textMsg, err := xml.Marshal(&repTextMsg)
	if err != nil {
		log.Printf("[消息回复] - 将对象进行XML编码出错: %v\n", err)
		return
	}
	_, _ = c.Writer.Write(textMsg)
}
