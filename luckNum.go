package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/beevik/etree"
	"go11x5/mysql"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func getLuckNum(prefix string) error {
	// 得到最新的一期的未出数字和遗漏值
	var baseURL, sCookie string
	if prefix == "jx/" {
		baseURL = "https://chart.ydniu.com/trend/syx5jx/"
		sCookie = "sign_56fb72c935c416697d92f8ef2b2f7c47=49686a06834c314ca8ce71a1f0013857; ASP.NET_SessionId=3wghqdrfyvkovt5qqrykhyzf"
	} else if prefix == "gd/" {
		baseURL = "https://chart.ydniu.com/trend/syx5gd/"
		sCookie = "sign_dcd8fbf688814a2d4adb5bba8a184d14=3b3e2e0a980ec17c6eb7a8d6bc820a54; ASP.NET_SessionId=3wghqdrfyvkovt5qqrykhyzf"
	} else {
		log.Fatalln("prefix 不受支持")
	}
	c := http.Client{Timeout: 60 * time.Second}
	form := url.Values{}
	form.Add("method", "CheckUpdate")
	form.Add("index", "0")
	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBufferString(form.Encode()))
	if err != nil {
		log.Println("create http request error!")
		return err
	}
	req.Header.Set("Cookie", sCookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:78.0) Gecko/20100101 Firefox/78.0")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("sign", strconv.Itoa(int(time.Now().UnixNano())/1000000))
	resp, err := c.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	type Res11x5 struct {
		Result  string `json:"result"`
		Success bool   `json:"success"`
	}
	// 开始解析拿到的数据
	r := Res11x5{}
	err = json.NewDecoder(resp.Body).Decode(&r)
	//err = json.Unmarshal(respBytes, &r)
	if err != nil {
		log.Println("解析io.Reader时出错: ", err)
		return err
	}
	//fmt.Println("解析之后的数据是: ", r.Result)
	// 这玩意是xml, 得用xml的方式去解析
	doc := etree.NewDocument()
	if err := doc.ReadFromString(r.Result); err != nil {
		log.Println("解析XML时出错: ", err)
		return err
	}

	root := doc.SelectElement("tbody")
	trList := root.SelectElements("tr")
	// 获取最近的1期数值
	//m := make(map[string][]string, 0)
	iwant := trList[len(trList)-1] // 拿到最新的tr
	//fmt.Printf("最新的一期数据是: %+v\n", iwant)
	numList := make([]string, 11)
	for i, tds := range iwant.SelectElements("td") {
		tdsAttr := tds.SelectAttrValue("class", "unknown")
		if tdsAttr == "y" {
			numList[i-6] = tds.Text()
		}
		if tdsAttr == "lan" {
			numList[i-6] = "N"
		}
	}
	m := make(map[int]int)
	for i, v := range numList {
		if v != "N" {
			m[i+1], _ = strconv.Atoi(v)
		}
	}
	// 得到数据库中遗漏值和期望之间的关系
	showThink(prefix, 0) // 首先, 更新luck表
	// 拿到luck表中specific_num, leave_value, stop_probability, hope_income 数值
	luckNumList, err := mysql.GetDataFromLuckTable(prefix)
	if err!=nil {
		log.Println(err)
		return err
	}
	//fmt.Println(luckNumList)
	// 看看最新一期里面有没有什么机会?
	fmt.Println(prefix, "这一期: ", iwant.SelectElement("td").Text())
	for k, v := range m{
		for _, luckNum := range luckNumList {
			if k==luckNum.SpecificNum && v==luckNum.LeaveValue {
				fmt.Printf("可选数字: %d, 遗漏数值: %d, 停止概率: %.6f, 数学期望: %.6f\n",luckNum.SpecificNum, luckNum.LeaveValue, luckNum.StopProbability, luckNum.HopeIncome)
			}
		}
	}
	return nil
}
