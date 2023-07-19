package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/robfig/cron/v3"
)

const Layout = "2006-01-02 15:04:05"
const Layout2 = "2006-01-02"

/*
*

	couponClock: 定时任务
	description: 会在指定时间抢券,每次仅抢券一张,有重试机制，过了抢券时间的一分钟后，程序会结束
	period 即将抢券的时间段
	amount 抢券的金额
	accessToken 密钥

*
*/
func couponClock(period string, amount int, accessToken string) {
	c := cron.New(cron.WithSeconds())

	location, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")
	day := now.Format("02")
	// 关闭程序时间，以免造成资源浪费
	closeTime, _ := time.ParseInLocation(Layout, fmt.Sprintf("%v-%v-%v %v:01:00", year, month, day, period), location)

	p, _ := strconv.Atoi(period)
	// “2006-01-02 15:04:05”是Go语言的创建时间，且必须为这几个准确的数字。
	// 指定时间执行，cron格式（秒，分，时，天，月，周）	spec := fmt.Sprintf("00 00 %v %v %v ?", period, day, month)
	// 提前一秒钟执行，定时任务会在前500ms内发送请求
	spec1 := fmt.Sprintf("00 00 %v %v %v ?", p, day, month)
	spec2 := fmt.Sprintf("01 00 %v %v %v ?", p, day, month)
	spec3 := fmt.Sprintf("02 00 %v %v %v ?", p, day, month)
	log.Println(spec1, spec2, spec3)

	title, message := "恭喜你，抢券成功", "请前往健身地图核验是否到账~"
	// 消费券code 不变
	stockId, err := getStockId(amount)
	if err != nil {
		title, message = "很遗憾！-1", err.Error()
		// send bark to phone
		noticeMasterPhone(title, message)
		return
	}

	log.Println("what time is it now ? ", time.Now())

	// 80 50的券限制了每周一次，隔一秒抢一次，防止因为服务器挂了的缘故导致00发起的请求失败
	if amount == 80 || amount == 50 {
		c.AddFunc(spec1, func() {
			send(period, stockId, accessToken)
		})
	} else {

		c.AddFunc(spec3, func() {
			for i := 0; i < 8; i++ {
				log.Println("before send time : ", time.Now(), accessToken)
				flag := send(period, stockId, accessToken)
				time.Sleep(time.Second * 2)
				log.Println("after send time : ", time.Now(), accessToken)
				if flag {
					break
				}
			}
		})
	}

	c.Start()
	defer c.Stop()

	/**
		主线程一直睡眠
		select {}
		主线程睡眠2小时后关闭,token最长有效期为2h, 但是一直开启线程有种浪费的感觉
		time.Sleep(time.Hour * 2)
	**/
	// 在任务结束的5分钟后结束定时任务 测试通过
	for {
		if time.Now().After(closeTime) {
			log.Println("结束！")
			break
		}
	}
}

/*
*

	随机返回true or false

*
*/
func sendtest(hour string, price int, accessToken string) bool {
	num := rand.Float32()
	flag := num < 0.5
	fmt.Printf("now: %v, hour: %v,flag: %v, num:%v", time.Now(), hour, flag, num)
	noticeMasterPhone("测试定时任务", fmt.Sprintf("准点测试时间%v", time.Now()))
	return flag
}
