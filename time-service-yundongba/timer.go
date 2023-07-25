package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

/*
*

	6228127c62ff4125c690ea50 5元消费券
	6228149062ff4125c690ea51 10元消费券
	6228153462ff4125c690ea52 20元消费券
	622815b562ff4125c690ea53 30元消费券
	62299598fceddb10cd1cb64d 50元消费券
	6229976ffceddb10cd1cb64f 80元消费券

*
*/
const (
	Coupons5  string = "6228127c62ff4125c690ea50"
	Coupons10        = "6228149062ff4125c690ea51"
	Coupons20        = "6228153462ff4125c690ea52"
	Coupons30        = "622815b562ff4125c690ea53"
	Coupons50        = "62299598fceddb10cd1cb64d"
	Coupons80        = "6229976ffceddb10cd1cb64f"
)

func send(period string, stockId string, accessToken string) bool {
	time.Sleep(time.Duration(waitTime) * time.Millisecond)
	log.Println(time.Now(), "post start")

	url := "https://mapv2.51yundong.me/api/coupon/coupons/send?stockId=" + stockId + "&time=" + period + "%3A00"
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Host", "mapv2.51yundong.me")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Accept-Encoding", "gzip,compress,br,deflate")
	req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.29(0x18001d30) NetType/4G Language/zh_CN")
	req.Header.Add("Referer", "https://servicewechat.com/wx8b97e9b9a6441e29/175/page-frame.html")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	log.Println(time.Now(), " post over", accessToken)
	log.Println("step 1")
	title, message := "恭喜你，抢券成功", "请前往健身地图核验是否到账~"
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)
	log.Println("step 2")
	result := transformation(res)
	fmt.Printf("status: %v, response: %v", res.StatusCode, result["msg"])
	switch res.StatusCode {
	case http.StatusOK:
		log.Println("step 3")
		noticeMasterPhone(title, message)
		return true
	case http.StatusUnauthorized:
		log.Println("step 4")
		title, message = "很遗憾！-2", "当前用户token已经过期"
		noticeMasterPhone(title, message)
		return true
	case http.StatusServiceUnavailable:
		// 503 也不断重试
		log.Println("step 5")
		title, message = "警告警告~", result["msg"]
	default:
		log.Println("step 6")
		title, message = "很遗憾！-3", fmt.Sprintf("错误，请检查代码, status: %d, response: %v", res.StatusCode, result["msg"])
		break
	}
	fmt.Println(title)
	fmt.Println(message)
	noticeMasterPhone(title, message)
	return false
}

func getStock(accessToken string) error {

	url := "https://mapv2.51yundong.me/api/user/users/info?noShowLoading=true"
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Host", "mapv2.51yundong.me")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Accept-Encoding", "gzip,compress,br,deflate")
	req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.29(0x18001d30) NetType/4G Language/zh_CN")
	req.Header.Add("Referer", "https://servicewechat.com/wx8b97e9b9a6441e29/175/page-frame.html")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("请求失败！")
		}
	}(res.Body)
	log.Println(res.StatusCode)

	// body, _ := ioutil.ReadAll(res.Body)
	// fmt.Printf("result res:%v, body:%v", res.StatusCode, string(body))
	switch res.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		return errors.New("当前用户token已经过期")
	default:
		return errors.New("还有其他情况的错误，请检查代码")
	}

}

func getStockId(price int) (string, error) {
	var stockId string
	switch price {
	case 5:
		stockId = Coupons5
		break
	case 10:
		stockId = Coupons10
		break
	case 20:
		stockId = Coupons20
		break
	case 30:
		stockId = Coupons30
		break
	case 50:
		stockId = Coupons50
		break
	case 80:
		stockId = Coupons80
		break
	default:
		return "", errors.New("please choose price")
	}
	return stockId, nil
}

func transformation(response *http.Response) map[string]string {
	result := make(map[string]string)
	body, err := ioutil.ReadAll(response.Body)
	if err == nil {
		err := json.Unmarshal([]byte(string(body)), &result)
		if err != nil {
			return nil
		}
	}
	fmt.Printf("查看http请求返回结果：%v", result)
	return result
}

// “2006-01-02 15:04:05”是Go语言的创建时间，且必须为这几个准确的数字。
// 测试通过
func getPeriodTime() string {

	// 设置时区，否则默认 UTC=美国时间
	location, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")
	day := now.Format("02")

	morning, _ := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%v-%v-%v 08:00:00", year, month, day), location)
	// fmt.Println(morning)
	noon, _ := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%v-%v-%v 12:00:00", year, month, day), location)
	// fmt.Println(noon)
	night, _ := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%v-%v-%v 18:00:00", year, month, day), location)
	// fmt.Println(night)

	if now.Before(morning) {
		return "08"
	} else if now.Before(noon) {
		return "12"
	} else if now.Before(night) {
		return "18"
	} else {
		return "18"
	}

}
