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
	log.Printf("[消息接收] - 收到消息, 发送消息方：%s, 接受消息方：%s, 消息类型为: %s, 消息内容为: %s\n", textMsg.FromUserName, textMsg.ToUserName, textMsg.MsgType, content)

	bool, _ := regexp.MatchString("Bearer [a-zA-Z0-9]{8}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{12}\\+[1-8]0$", content)

	if bool {
		// Bearer 702af673-d508-4fa6-9ddd-f93944b07c4b+80
		split := strings.Split(content[7:], "+")
		checkAndRob(split[1], split[0], c, textMsg.ToUserName, textMsg.FromUserName)
	} else {
		returnContent := "您的输入的格式不正确，请重新整理后再试！"
		// 对接收的消息进行被动回复
		WXMsgReply(c, textMsg.ToUserName, textMsg.FromUserName, returnContent)
	}

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

func checkAndRob(amount string, accessToken string, c *gin.Context, fromUser string, toUser string) {
	var returnMsg string
	// 获取抢券时间段
	period := getPeriodTime()
	// 对比当前时间和抢券时间，确保在token的两个小时以内
	location, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")
	day := now.Format("02")

	// 抢券时间
	robTime, _ := time.ParseInLocation(Layout, fmt.Sprintf("%v-%v-%v %v:00:00", year, month, day, period), location)
	sub2hours, _ := time.ParseDuration("-2h")
	// 有效期时间 - 抢券前两小时
	efficientStartTime := robTime.Add(sub2hours)

	if now.Before(efficientStartTime) {
		returnMsg = fmt.Sprintf("将为您抢券的时间是:%v点整\n现在的时间是：%v\n不在token的有效期内\n请详细阅读抢券手册再重新发起抢券请求", period, now.Format(Layout))
		WXMsgReply(c, fromUser, toUser, returnMsg)
		return
	}

	// 有效期时间 - 抢券开始后1分钟
	add1min, _ := time.ParseDuration("1m")
	efficientEndTime := robTime.Add(add1min)
	if now.After(efficientEndTime) {
		returnMsg = fmt.Sprintf("将为您抢券的时间是:%v点整\n现在的时间是：%v\n抢券已经开始超过1分钟\n不提供这种无意义的抢券服务\n请详细阅读抢券手册再重新发起抢券请求", period, now.Format(Layout))
		WXMsgReply(c, fromUser, toUser, returnMsg)
		return
	}

	err := getStock(period, accessToken)
	if err != nil {
		returnMsg = err.Error()
		WXMsgReply(c, fromUser, toUser, returnMsg)
		return
	} else {
		returnMsg = fmt.Sprintf("恭喜您,当前token有效!将于今日%v点整为您抢%v元消费券", period, amount)
	}

	// 金额转换
	int, err := strconv.Atoi(amount)
	if err != nil {
		int = 30
	}

	// 对接收的消息进行被动回复
	WXMsgReply(c, fromUser, toUser, returnMsg)

	robStruct := RobStruct{
		period:      period,
		amount:      int,
		accessToken: accessToken,
	}
	robMap[fromUser] = robStruct

}
