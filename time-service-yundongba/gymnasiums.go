package main

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

/**
	Host: wx-api.papa.com.cn
	Origin: https://wx.papa.com.cn
	Accept-Encoding: gzip, deflate, br
	Connection: keep-alive
	Accept: application/json, text/plain,
	User-Agent: Mozilla/5.0 (iPhone; CPU iPhone OS 16_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.29(0x18001d35) NetType/4G Language/zh_CN
	Referer: https://wx.papa.com.cn/
	Content-Length: 120
	Accept-Language: zh-CN,zh-Hans;q=0.9
	Content-Type: application/x-www-form-urlencoded; charset=UTF-8
**/

/*
*

	订阅场馆，监听是否有场地
	需要保持连续不中断
	param: month - 月
	param：day - 日
	param: period - 订阅时间段
	return : int 0-失败 1-成功 2-其他情况
	return : error 具体失败信息

*
*/
func subscribeGymnasiums(month string, day string, start string, long int) (int, string) {
	now := time.Now()
	year := now.Year()
	_, err := time.Parse(Layout2, fmt.Sprintf("%v-%v-%v", year, month, day))
	if err != nil {
		fmt.Println(err.Error())
		return 0, "你选择的不是一个正常日期!"
	}
	url := "https://wx-api.papa.com.cn/v2"
	reqBody := []byte("client_type=browser&sport_tag_id=8&date_str=2022-11-21&r=stadia.skuList&access_token_wx=661c3961c7688967c0ce0533926c8535")
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))

	req.Header.Add("Host", "wx-api.papa.com.cn")
	req.Header.Add("Origin", "https://wx.papa.com.cn")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Encoding", "gzip,compress,br,deflate")
	req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.29(0x18001d35) NetType/4G Language/zh_CN")
	req.Header.Add("Referer", "https://wx.papa.com.cn/")
	req.Header.Add("Content-Length", "120")
	req.Header.Add("Accept-Language", "zh-CN,zh-Hans;q=0.9")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	// 校验请求是否成功
	switch res.StatusCode {
	case http.StatusOK:
		break
	case http.StatusUnauthorized:
		return 0, "当前用户token已经过期,请联系提供者"
	default:
		return 0, "还有其他情况的错误,请联系提供者检查代码"
	}

	// 解析选择时间是否有场地，可组合
	result := switchContentEncoding(res)
	return 1, fmt.Sprintf("%v", result)

}

func enough(gym Gym, start int, end int) string {
	if start < 8 || start > 21 || start > end || end < 9 || end > 22 {
		return ""
	}
	// 0~13
	skuList := gym.SkuList
	startInd, endInd := start-8, end-8

	result := make(map[int]bool)
	for i := 0; i < endInd-startInd; i++ {
		result[start+i] = isEmpty(skuList[i])
	}

	var builder strings.Builder
	for key := range result {
		if result[key] {
			fmt.Fprint(&builder, fmt.Sprintf("%v点~%v点  场地空闲~\n", key, key+1))
		}
	}

	// 这里不会报空指针
	return builder.String()
}

/*
*

	默认剔除室内场地 index = 0, 1
	默认剔除培训付费场地 index 6, 7

*
*/
func isEmpty(sku []Sku) bool {
	for i := 2; i < 6; i++ {
		// 场地没有被锁定
		if !sku[i].IsLock {
			return true
		}
	}
	return false
}

/*
*

	预订场馆

*
*/
func bookGymnasiums() {

}

/*
*

	检测返回的body是否经过压缩，并返回解压的内容

*
*/
func switchContentEncoding(res *http.Response) Gym {
	var gymStr Gym
	var bodyReader io.Reader
	var err error
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		bodyReader, err = gzip.NewReader(res.Body)
		break
	case "deflate":
		bodyReader = flate.NewReader(res.Body)
		break
	default:
		bodyReader = res.Body
		break
	}
	body, err := ioutil.ReadAll(bodyReader)
	// fmt.Println("body:" + string(body))
	if err == nil {
		json.Unmarshal([]byte(body), &gymStr)
	}
	return gymStr
}

type Gym struct {
	V            string            `json:"v"`
	R            string            `json:"r"`
	StadiumId    string            `json:"stadium_id"`
	DateStr      string            `json:"date_str"`
	SportTag     string            `json:"sport_tag"`
	SportTagList map[string]string `json:"sport_tag_list"`
	FieldStr     []string          `json:"fieldStr"`
	SkuList      []([]Sku)         `json:"skuList"`
	TimeStr      []string          `json:"timeStr"`
}

type Sku struct {
	Sku         string `json:"sku"`
	Remark      string `json:"remark"`
	SkuName     string `json:"sku_name"`
	FieldName   string `json:"field_name"`
	TimeId      string `json:"time_id"`
	TimeStr     string `json:"time_str"`
	SportTagStr string `json:"sport_tag_str"`
	Price       string `json:"price"`
	IsLock      bool   `json:"is_lock"`
	LockMsg     string `json:"lock_msg"`
	IsGroup     bool   `json:"is_group"`
	GroupId     string `json:"group_id"`
	LockStatus  int    `json:"lock_status"`
}
