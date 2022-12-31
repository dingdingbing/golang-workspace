package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

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
}
