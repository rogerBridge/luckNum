package main

import (
	"bytes"
	"log"
	"net/http"
	"net/url"
	"time"
)

func pushMsgToBot(msg string) error {
	c := http.Client{
		Timeout: 60 * time.Second,
	}
	content := url.Values{
		"chat_id": []string{"615491801"},
		"text": []string{msg},
	}
	//fmt.Println(content.Encode())
	req, err := http.NewRequest(http.MethodPost, "http://fenr.men:4000/tg?method=sendMessage", bytes.NewBufferString(content.Encode()))
	req.Header.Add("Authorization", "c371b934-b18c-4a44-a4a9-4d830fd0a527")
	if err!=nil {
		log.Println(err)
		return err
	}
	_, err = c.Do(req)
	if err!=nil {
		log.Println(err)
		return err
	}
	return nil
}