package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/robfig/cron/v3"
)

func cronHere(hour string) {
	c := cron.New(cron.WithSeconds())

	// “2006-01-02 15:04:05”是Go语言的创建时间，且必须为这几个准确的数字。
	// 指定时间执行，cron格式（秒，分，时，天，月，周）
	// spec := fmt.Sprintf("00 00 %v %v %v ?", hour, time.Now().Format("02"), time.Now().Format("01"))
	spec := fmt.Sprintf("00 13 %v %v %v ?", hour, time.Now().Format("02"), time.Now().Format("01"))
	fmt.Println(spec)
	var i int = 1

	c.AddFunc(spec, func() {
		// 尝试三次 测试成功
		for {
			flag := sendtest(hour, 30)
			i++
			if flag || i > 3 {
				break
			}

		}
	})
	c.Start()
	defer c.Stop()
	time.Sleep(time.Second * 60)
}

/*
*

	随机返回true or false

*
*/
func sendtest(hour string, price int) bool {
	num := rand.Float32()
	flag := num > 0.5
	fmt.Printf("now: %v, hour: %v,flag: %v, num:%v", time.Now(), hour, flag, num)
	return flag
}
