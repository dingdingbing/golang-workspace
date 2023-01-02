package main

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
)

var waitGroup sync.WaitGroup

func main() {
	router := gin.Default()
	router.SetTrustedProxies([]string{"122.51.126.249"})
	router.GET("/coupon/auth/:token", checkAuth)
	router.GET("/coupon/test/:token/:amount", testReturn)
	router.GET("/coupon/rob/:token/:amount", robCoupon)
	router.GET("/", func(c *gin.Context) {
		// If the client is 192.168.1.2, use the X-Forwarded-For
		// header to deduce the original client IP from the trust-
		// worthy parts of that header.
		// Otherwise, simply return the direct client IP
		fmt.Printf("ClientIP: %s\n", c.ClientIP())
	})
	router.GET("/wx", WXCheckSignature)
	router.POST("/wx", WXMsgReceive)
	router.Run(":80")

	// 异步test
	// fmt.Println("work1 return", work1(), time.Now())
	// fmt.Println("main1", time.Now())
	// waitGroup.Wait()
	// fmt.Println("main2", time.Now())
}

// func work1() string {
// 	// sw1 -sw2-ew1-nw2
// 	fmt.Println("start - work1", time.Now())
// 	waitGroup.Add(1)
// 	go work2()
// 	time.Sleep(time.Second * 1)
// 	fmt.Println("end - work1", time.Now())
// 	return "g"

// }

// func work2() {
// 	fmt.Println("start - work2", time.Now())
// 	time.Sleep(time.Second * 10)
// 	fmt.Println("end - work2", time.Now())
// 	waitGroup.Done()
// }
