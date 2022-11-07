package main

import "net/http"

func noticeMasterPhone(title string, content string) {
	http.Get("https://api.day.app/RYXFHftgRhq5BsomYwEb5J/" + title + "/" + content)
}

func noticeConsumerPhone(title string, content string) {
	// 需要替换成 陈雅的通知
	http.Get("https://api.day.app/RYXFHftgRhq5BsomYwEb5J/" + title + "/" + content)
}
