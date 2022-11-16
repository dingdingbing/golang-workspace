package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.GET("/coupon/auth/:token", checkAuth)
	router.GET("/coupon/test/:token/:amount", testReturn)
	router.GET("/coupon/rob/:token/:amount", robCoupon)
	router.Run("localhost:8080")
	// couponClock("09", 30, "token")
	// int, str := subscribeGymnasiums("11", "13", 14, 18)
	// fmt.Printf("code: %v, msg: %v", int, str)
	// var gym gym
	// gymmap := make(map[string]string)
	// str := `{"v":"v2","r":"stadia.skuList"}`
	// json.Unmarshal([]byte(str), &gym)
	// json.Unmarshal([]byte(str), &gymmap)
	// fmt.Println(gym)
	// fmt.Println(gymmap)
	// test 数组添加
	// var groundInfo groundInfo
	// groundInfo.free = false
	// groundInfo.number = append(groundInfo.number, 2)
	// groundInfo.number = append(groundInfo.number, 2)
	// groundInfo.number = append(groundInfo.number, 2)
	// groundInfo.number = append(groundInfo.number, 2)

	// for _, value := range groundInfo.number {
	// 	fmt.Println(value)
	// }
}
