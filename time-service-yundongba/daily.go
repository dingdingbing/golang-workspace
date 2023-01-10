package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

func initDaily() {
	log.Println("init start")
	c := cron.New(cron.WithSeconds())

	// 每天抢券前30秒核对是否有人发起了请求
	spec := fmt.Sprintf("30 59 7 * * ?")
	c.AddFunc(spec, everyDayEightoClock)

	spec = fmt.Sprintf("30 59 11 * * ?")
	c.AddFunc(spec, everyDayTwelveoClock)

	spec = fmt.Sprintf("30 59 17 * * ?")
	c.AddFunc(spec, everyDaySeventeenoClock)

	spec = fmt.Sprintf("00 00 * * * ?")
	c.AddFunc(spec,
		func() {
			log.Println("系统功能正常！")
		})
	log.Println("init end")

	c.Start()
	defer c.Stop()

	select {}
}

func everyDayEightoClock() {

	log.Println("everyDayEightoClock check ")

	// 1. 查看map是否为空
	if len(robMap) == 0 {
		log.Println("today no mission for 8 o'clock ")
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
	log.Println("start asyncCouponClock for ", robStruct, time.Now())
	couponClock(robStruct.period, robStruct.amount, robStruct.accessToken)
	waitGroup.Done()
}

func everyDayTwelveoClock() {

	// 1. 查看map是否为空
	if len(robMap) == 0 {
		log.Println("today no mission for 12 o'clock ")
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

func everyDaySeventeenoClock() {

	// 1. 查看map是否为空
	if len(robMap) == 0 {
		log.Println("today no mission for 18 o'clock ")
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
