package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// WXTextMsg 微信文本消息结构体
type WXTextMsg struct {
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

// WXMsgReceive 微信消息接收
func WXMsgReceive(c *gin.Context) {
	var textMsg WXTextMsg
	err := c.ShouldBindXML(&textMsg)
	if err != nil {
		log.Printf("[消息接收] - XML数据包解析失败: %v\n", err)
		return
	}

	content := textMsg.Content
	log.Printf("[消息接收] - 收到消息, 消息类型为: %s, 消息内容为: %s\n", textMsg.MsgType, content)

	bool, _ := regexp.MatchString("Bearer [a-zA-Z0-9]{8}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{12}\\+[1-8]0$", content)

	var returnContent string
	if bool {
		split := strings.Split(content, "+")
		returnContent = checkAndRob(split[0], split[1])
	} else {
		returnContent = "您的输入的格式不正确，请重新整理后再试！"
	}
	// 对接收的消息进行被动回复
	WXMsgReply(c, textMsg.ToUserName, textMsg.FromUserName, returnContent)
}

// WXMsgReply 微信消息回复
func WXMsgReply(c *gin.Context, fromUser string, toUser string, content string) {
	repTextMsg := WXRepTextMsg{
		ToUserName:   toUser,
		FromUserName: fromUser,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      content,
	}

	msg, err := xml.Marshal(&repTextMsg)
	if err != nil {
		log.Printf("[消息回复] - 将对象进行XML编码出错: %v\n", err)
		return
	}
	_, _ = c.Writer.Write(msg)
}

func checkAndRob(amount string, accessToken string) string {
	var returnMsg string
	// 获取抢券时间段
	period := getPeriodTime()
	periodInt, _ := strconv.Atoi(period)
	// 对比当前时间和抢券时间，确保在token的两个小时以内
	location, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")
	day := now.Format("02")
	// 有效期时间
	efficientTime, _ := time.ParseInLocation(Layout, fmt.Sprintf("%v-%v-%v %v:00:00", year, month, day, periodInt-2), location)
	if now.Before(efficientTime) {
		returnMsg = fmt.Sprintf("将为您抢券的时间是:%v,现在的时间是：%v,不在token的有效期内,请详细阅读抢券手册再重新发起抢券请求", period, now.Format(Layout))
		return returnMsg
	}

	err := getStock(period, accessToken)
	if err == nil {
		return err.Error()
	} else {
		returnMsg = fmt.Sprintf("恭喜您,当前token有效!将于今日%v点整为您抢%v元消费券", period, amount)
	}

	// 金额转换
	int, err := strconv.Atoi(amount)
	if err != nil {
		int = 30
	}

	couponClock(period, int, accessToken)
	return returnMsg
}
