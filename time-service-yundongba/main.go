package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"sync"
)

// key-> time value -> RobStruct
var robMap map[string]RobStruct
var waitTime int64 = 800

type RobStruct struct {
	period      string
	amount      int
	accessToken string
}

var waitGroup sync.WaitGroup

func main() {
	robMap = make(map[string]RobStruct)
	waitGroup.Add(1)
	go initDaily()
	router := gin.Default()
	err := router.SetTrustedProxies([]string{"122.51.126.249"})
	if err != nil {
		log.Println("应用启动失败！")
		return
	}
	router.GET("/wx", WXCheckSignature)
	router.POST("/wx", WXMsgReceive)
	err = router.Run(":9000")
	if err != nil {
		log.Println("应用启动失败！")
		return
	}
	waitGroup.Wait()

	//log.Println(getStock("0e2afab6-f40b-45be-9995-bbe8008570ce"))
}
