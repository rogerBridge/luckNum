package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/beevik/etree"
	"go11x5/mysql"
	"io/ioutil"
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
	showThink(prefix, 0, hopeWin) // 首先, 更新luck表
	// 拿到luck表中specific_num, leave_value, stop_probability, hope_income 数值
	luckNumList, err := mysql.GetDataFromLuckTable(prefix)
	if err!=nil {
		log.Println(err)
		return err
	}
	//fmt.Println(luckNumList)
	// 看看最新一期里面有没有什么机会?
	fmt.Println(prefix, "这一期: ", iwant.SelectElement("td").Text())
	nearestOrderNum := iwant.SelectElement("td").Text()
	for k, v := range m{
		for _, luckNum := range luckNumList {
			if k==luckNum.SpecificNum && v==luckNum.LeaveValue {
				fmt.Printf("可选数字: %d, 遗漏数值: %d, 停止概率: %.6f, 数学期望: %.6f\n",luckNum.SpecificNum, luckNum.LeaveValue, luckNum.StopProbability, luckNum.HopeIncome)
				forecastOrderNum, err := nextOne(nearestOrderNum)
				if err!=nil {
					log.Println(err)
				}
				PushMsg(fmt.Sprintf("%s:forecast:%s 可选数字: %d, 遗漏数值: %d, 停止概率: %.6f, 数学期望: %.6f\n", prefix, forecastOrderNum, luckNum.SpecificNum, luckNum.LeaveValue, luckNum.StopProbability, luckNum.HopeIncome))
				// 将预测的数据存入forecast_{jx|gd} 表格
				err = mysql.StoreResultToForecastTable(prefix, forecastOrderNum, luckNum.SpecificNum)
				if err!=nil {
					log.Println(err)
				}
			}
		}
	}
	return nil
}

// 生成下一期orderNum
func nextOne(thisOne string) (string,error) {
	if len(thisOne) != 10 {
		return "", errors.New("error length")
	}
	dateStr := thisOne[:8]
	number := thisOne[8:]
	numberInt, _ := strconv.Atoi(number)

	if numberInt < 42 {
		return dateStr+fmt.Sprintf("%02d", numberInt+1), nil
	}else if numberInt == 42{
		timeThisOne, err := time.Parse("20060102", dateStr)
		if err!=nil {
			return "", err
		}
		return timeThisOne.AddDate(0,0,1).Format("20060102")+"01", nil
	}else {
		return "", errors.New("unknown error")
	}
}

func PushMsg(content string) {
	url := "http://sthink.top:8080/pushMsg"
	reqBody, err := json.Marshal(map[string]string {
		"content": content,
	})
	if err!=nil {
		log.Printf("make reqBody error \n")
		return
	}
	c := http.Client{Timeout: 60*time.Second}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err!=nil {
		log.Fatalln(err)
	}
	req.Header.Set("Authorization", "eb356e0f-2ad8-4352-9c1c-1af1db033f81")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
	if err!=nil {
		log.Printf("%+v\n", err)
		return
	}
	defer resp.Body.Close()

	Bytebody, err := ioutil.ReadAll(resp.Body)
	if err!=nil {
		log.Printf("%+v\n", err)
		return
	}
	log.Println(string(Bytebody))
}