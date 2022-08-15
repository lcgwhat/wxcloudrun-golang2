package service

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/dao"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/util"
)

// 与填写的服务器配置中的Token一致
const Token = "clear_love"
const failed = "失败"

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
		res := strings.Split(textMsg.Content, "|")
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
			err = doCheck(uid, textMsg.FromUserName, res[1], note)
			if err != nil {
				msg += err.Error()
			} else {
				msg = "【签到成功】"
			}

		}
		if len(res) == 1 {
			stringTime := res[0]
			loc, _ := time.LoadLocation("Local")
			the_time, err := time.ParseInLocation("2006-01", stringTime, loc)
			if err != nil {
				msg = err.Error()
			}
			failCount := dao.CheckImp.GetCheckFail(textMsg.FromUserName, the_time)
			count := dao.CheckImp.GetCheckMonth(textMsg.FromUserName, the_time)
			res := 0.0
			if count > 0 {
				res = float64(failCount) / float64(count)
				res = res * 100
			}
			msg = fmt.Sprintf("%s共出错了%v, 累计签到次数%v,结果%.2f%%", stringTime, failCount, count, res)
		}
	}
	// 对接收的消息进行被动回复
	WXMsgReply(ctx, textMsg.ToUserName, textMsg.FromUserName, msg)
}

var ErrorChecked = errors.New("今天已经签到过了...")

func doCheck(uid int, openId string, result string, note string) error {
	checkUser, err := dao.CheckUserImp.GetCheckUserById(int32(uid))
	now := time.Now()
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		// 一天前
		d, _ := time.ParseDuration("-24h")
		d1 := now.Add(d)
		user := model.CheckUser{
			OpenId:        openId,
			Username:      fmt.Sprintf("%s%v", "liu", rand.Float32()),
			CreateTime:    time.Now(),
			ContinuousDay: 0,
			LastDay:       d1,
		}
		err := dao.CheckUserImp.UpsertCherUser(&user)
		if err != nil {
			return err
		}
		checkUser = &user
	}

	existCheck := dao.CheckImp.ExistCheck(uid, now)
	if existCheck == true {
		return ErrorChecked
	} else {
		status := model.CheckStatus_1
		if result == failed {
			status = model.CheckStatus_3
		}
		check := model.Check{
			UserId:    uid,
			Username:  openId,
			CheckDate: now,
			Note:      note,
			Status:    status,
		}
		tx := db.Get().Begin()
		err := dao.CheckImp.UpsertCheck(&check)
		if err != nil {
			tx.Rollback()
			return err
		}
		ContinuousDay := 0
		if !util.InSameDay(checkUser.LastDay.Unix(), now.Unix()) {
			ContinuousDay = checkUser.ContinuousDay + 1
		}
		if result == failed {
			ContinuousDay = 0
		}
		err = dao.CheckUserImp.SetCheckDate(checkUser.Id, ContinuousDay, now)
		if err != nil {
			tx.Rollback()
			return err
		}
		return tx.Commit().Error
	}

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
