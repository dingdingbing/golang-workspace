package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func checkAuth(c *gin.Context) {
	accessToken := c.Param("token")

	period := getPeriodTime()
	err := getStock(period, accessToken)
	if err == nil {
		c.IndentedJSON(http.StatusOK, "当前token有效!")
	} else {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
	}
}

func robCoupon(c *gin.Context) {
	accessToken := c.Param("token")
	if len(accessToken) == 0 {
		c.IndentedJSON(http.StatusInternalServerError, "你传的token不对!")
		return
	}
	amount := c.Param("amount")
	period := getPeriodTime()
	// todo auth的http接口可以查看过期时间，后期看是否能够给准确提示。但是请求里面的code不好搞
	checkAuth(c)

	int, err := strconv.Atoi(amount)
	if err != nil {
		int = 30
	}

	noticeMasterPhone("恭喜您~", fmt.Sprintf("将于今日%v点整为您抢%v元消费券", period, amount))
	c.IndentedJSON(http.StatusOK, fmt.Sprintf("将于今日%v点整为您抢%v元消费券", period, amount))
	couponClock(period, int, accessToken)

}

/*
*
c.IndentedJSON 并不会打断程序执行
如果有多个，会返回多个json串

*
*/
func testReturn(c *gin.Context) {
	amount := c.Param("amount")
	int, err := strconv.Atoi(amount)
	if err != nil {
		int = 30
	}
	s, err := getStockId(int)
	accessToken := c.Param("token")
	fmt.Println("step finish")
	c.IndentedJSON(http.StatusOK, send("12", s, accessToken))
}
