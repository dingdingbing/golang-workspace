package main

import (
	"github.com/gin-gonic/gin"
)

// key-> time value -> RobStruct
var robMap map[string]RobStruct

type RobStruct struct {
	period      string
	amount      int
	accessToken string
}

func main() {
	robMap = make(map[string]RobStruct)
	router := gin.Default()
	router.SetTrustedProxies([]string{"122.51.126.249"})
	router.GET("/wx", WXCheckSignature)
	router.POST("/wx", WXMsgReceive)
	router.Run(":9000")

}
