package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.GET("/coupon/auth/:token", checkAuth)
	// router.GET("/coupon/auth/:amount", testReturn)
	router.GET("/coupon/rob/:token/:amount", robCoupon)
	router.Run("localhost:8080")
	// couponClock("09", 30, "token")

}
