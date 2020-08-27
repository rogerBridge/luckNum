package main

import (
	"bytes"
	"errors"
	"go11x5/play"
	"log"
	"net/http"
	"net/url"
	"time"
)

func main() {
	err := pushMsgToBot("hello")
	if err!=nil {
		log.Println(err)
	}
}

func pushMsgToBot(msg string) error {
	c := http.Client{
		Timeout: 30 * time.Second,
	}
	content := url.Values{
		"chat_id": []string{"615491801"},
		"text": []string{msg},
	}
	host := "https://log.fenr.men"
	// 开始加密
	encBytes := play.Encrypt(content.Encode())
	log.Println("Encrypted data is:", encBytes)
	// 将加密后的请求发送给自由节点
	req, err := http.NewRequest(http.MethodPost, host+"/bot/tg?method=sendMessage", bytes.NewBuffer(encBytes))
	if err!=nil {
		log.Println("try establish conn: ", err)
		return err
	}
	req.Header.Add("Authorization", "c371b934-b18c-4a44-a4a9-4d830fd0a527")
	resp, err := c.Do(req)
	if err!=nil {
		log.Println(err)
		return err
	}
	if resp.Header.Get("Complete") == "true" {
		return nil
	}else {
		log.Println("resp header: ", resp.Header)
		return errors.New("send msg fail")
	}
}