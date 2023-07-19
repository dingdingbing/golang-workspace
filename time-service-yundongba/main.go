package main

import (
	"log"
	"sync"
	"time"
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
	// robMap = make(map[string]RobStruct)
	// waitGroup.Add(1)
	// go initDaily()
	// router := gin.Default()
	// router.SetTrustedProxies([]string{"122.51.126.249"})
	// router.GET("/wx", WXCheckSignature)
	// router.POST("/wx", WXMsgReceive)
	// router.Run(":9000")
	// waitGroup.Wait()
	log.Println("1", time.Now())
	time.Sleep(time.Second * 2)
	log.Println("1", time.Now())
	time.Sleep(time.Second * 3)
	log.Println("1", time.Now())
	time.Sleep(time.Second * 5)
	log.Println("1", time.Now())
}
