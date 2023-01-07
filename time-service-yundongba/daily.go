package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

var waitGroup sync.WaitGroup

func initDaily() {
	c := cron.New(cron.WithSeconds())
	// “2006-01-02 15:04:05”是Go语言的创建时间，且必须为这几个准确的数字。
	// 指定时间执行，cron格式（秒，分，时，天，月，周）	spec := fmt.Sprintf("00 00 %v %v %v ?", period, day, month)
	spec := fmt.Sprintf("0 59 7 * * ? *")
	c.AddFunc(spec, everyDayEightoClock)

	spec = fmt.Sprintf("0 59 11 * * ? *")
	c.AddFunc(spec, everyDayTwelveoClock)

	spec = fmt.Sprintf("0 59 17 * * ? *")
	c.AddFunc(spec, everyDaySeventeenoClock)
}

func everyDayEightoClock() {

	// 1. 查看map是否为空
	if len(robMap) == 0 {
		return
	}

	title := fmt.Sprintf("现在是%s,有%d人参加抢券", time.Now().Format(Layout), len(robMap))
	var message strings.Builder
	message.WriteString("用户名：")
	for fromUser, robStruct := range robMap {
		fmt.Println(fromUser, robStruct)
		message.WriteString(fmt.Sprint(fromUser, "+", robStruct.amount, ","))
		waitGroup.Add(1)
		go asyncCouponClock(robStruct)
		delete(robMap, fromUser)
	}

	noticeMasterPhone(title, message.String())
	waitGroup.Wait()
}

func asyncCouponClock(robStruct RobStruct) {
	couponClock(robStruct.period, robStruct.amount, robStruct.accessToken)
	waitGroup.Done()
}

func everyDayTwelveoClock() {

	// 1. 查看map是否为空
	if len(robMap) == 0 {
		return
	}

}

func everyDaySeventeenoClock() {

	// 1. 查看map是否为空
	if len(robMap) == 0 {
		return
	}

}
