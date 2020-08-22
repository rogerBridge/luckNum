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

// 得到最新一期的未出数字和遗漏值, 并且将它们的值写入到forecast表中, 并且通过msgbot发送出来
func getLuckNum(prefix string) error {
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
	// m为: 数字:遗漏的数值
	m := make(map[int]int)
	for i, v := range numList {
		if v != "N" {
			m[i+1], _ = strconv.Atoi(v)
		}
	}

	// 更新luck表 && unluck表
	showThink(prefix, 0)
	// 拿到luck表中specific_num, leave_value, stop_probability, hope_income 数值

	luckNumList, err := mysql.GetDataFromLuckTable(prefix)
	if err != nil {
		log.Println(err)
		return err
	}
	// 拿到unluck表中specific_num, leave_value, stop_probability, hope_income 数值
	unluckNumList, err := mysql.GetDataFromUnLuckTable(prefix)
	if err!= nil {
		log.Println(err)
		return err
	}
	// 看看最新一期里面有没有什么机会?
	log.Println(prefix, "这一期: ", iwant.SelectElement("td").Text())
	nearestOrderNum := iwant.SelectElement("td").Text()
	// 这一期所有符合要求的LuckNum都将放入willOrderLuckNum
	var willOrderLuckNum []mysql.LuckNum
	for k, v := range m {
		for _, luckNum := range luckNumList {
			if k == luckNum.SpecificNum && v == luckNum.LeaveValue {
				log.Printf("luck可选数字: %d, 遗漏数值: %d, 停止概率: %.6f, 数学期望: %.6f\n", luckNum.SpecificNum, luckNum.LeaveValue, luckNum.StopProbability, luckNum.HopeIncome)
				willOrderLuckNum = append(willOrderLuckNum, luckNum)
			}
		}
	}
	// 这一期所有符合要求的UnLuckNum都将放入willOrderUnLuckNum
	var willOrderUnLuckNum []mysql.LuckNum
	for k, v := range m {
		for _, unluckNum := range unluckNumList {
			if k == unluckNum.SpecificNum && v == unluckNum.LeaveValue {
				log.Printf("unluck可选数字: %d, 遗漏数值: %d, 停止概率: %.6f, 数学期望: %.6f\n", unluckNum.SpecificNum, unluckNum.LeaveValue, unluckNum.StopProbability, unluckNum.HopeIncome)
				willOrderUnLuckNum = append(willOrderUnLuckNum, unluckNum)
			}
		}
	}
	// 如果最新一期有符合luck表要求的数值, 将它们中hopeIncome数值最大的那个添加进forecast表格中
	if len(willOrderLuckNum) > 0 {
		// 从willOrderLuckNum中找到数学期望最大的那个luckNum, 并将它写入到forecast_{jx/gd}表格中
		vWant := willOrderLuckNum[0]
		for _, v := range willOrderLuckNum {
			if v.HopeIncome > vWant.HopeIncome {
				vWant = v
			}
		}
		// 将预测的数学期望最高的luckNum写入forecast_{jx/gd}中
		forecastOrderNum, err := nextOne(nearestOrderNum)
		if err != nil {
			log.Println(err)
			return err
		}
		// PushMsg(fmt.Sprintf("%s:forecast:%s 可选数字: %d, 遗漏数值: %d, 停止概率: %.6f, 数学期望: %.6f\n", prefix, forecastOrderNum, vWant.SpecificNum, vWant.LeaveValue, vWant.StopProbability, vWant.HopeIncome))
		// 停止推送消息到 Bot
		//err = pushMsgToBot(fmt.Sprintf("%s:forecast:%s 可选数字: %d, 遗漏数值: %d, 停止概率: %.6f, 数学期望: %.6f\n", prefix, forecastOrderNum, vWant.SpecificNum, vWant.LeaveValue, vWant.StopProbability, vWant.HopeIncome))
		//if err!=nil {
		//	log.Println(err)
		//	return err
		//}

		// 首先, 将预测的结果存放到forecast表格中
		err = mysql.StoreResultToForecastTable(prefix, forecastOrderNum, vWant.SpecificNum)
		if err != nil {
			log.Println(err)
			return err
		}
		// 从原始数据表格中找到是否符合
		err = mysql.DetectForecast(prefix)
		if err != nil {
			// 这里有错误不可以直接return err, 如果这样, 新的预测数据就发不出去了
			// 比如: 你预测了一个数字, 然后在统计数据中找不到对应的值, (这种情况很常见), 那你的预测数据发不出去了
			log.Println("DetectForecast: ", err)
			//return err
		}
		msg, err := mysql.StatisticsForecast(prefix)
		if err != nil {
			log.Println("StatisticsForecast: ", err)
			return err
		}
		message := msg + fmt.Sprintf("%s %s %d", prefix, forecastOrderNum, vWant.SpecificNum)
		// 尝试三次, 将消息发送出去
		log.Println("luck: \n" + message)
		//go sendMsgThreeTimes("luck: \n" + message)
		//return errors.New("send msg fail")
	} else {
		err := mysql.DetectForecast(prefix)
		if err != nil {
			log.Println("没有预测数值时, DetectForecast: ", err)
			return err
		}
	}


	// 如果最新一期有符合unluck表要求的数值, 将它们中hopeIncome数值最小的那个添加进forecast表格中
	if len(willOrderUnLuckNum) > 0 {
		// 从willOrderLuckNum中找到数学期望最小的那个luckNum, 并将它写入到forecast_{jx/gd}表格中
		// 选出与5.0/11.0最接近的猜想数值
		//vWant := willOrderUnLuckNum[0]
		//for _, v := range willOrderUnLuckNum {
		//	if math.Abs(v.StopProbability-5.0/11.0) < math.Abs(vWant.StopProbability-5.0/11.0) {
		//		vWant = v
		//	}
		//}
		// 将所有符合要求的统统写入forecast2_{gd|jx}表格之中

		// 将预测的数学期望最低的unluckNum写入forecast_{jx/gd}中
		forecastOrderNum, err := nextOne(nearestOrderNum)
		if err != nil {
			log.Println(err)
			return err
		}
		// PushMsg(fmt.Sprintf("%s:forecast:%s 可选数字: %d, 遗漏数值: %d, 停止概率: %.6f, 数学期望: %.6f\n", prefix, forecastOrderNum, vWant.SpecificNum, vWant.LeaveValue, vWant.StopProbability, vWant.HopeIncome))
		// 停止推送消息到 Bot
		//err = pushMsgToBot(fmt.Sprintf("%s:forecast:%s 可选数字: %d, 遗漏数值: %d, 停止概率: %.6f, 数学期望: %.6f\n", prefix, forecastOrderNum, vWant.SpecificNum, vWant.LeaveValue, vWant.StopProbability, vWant.HopeIncome))
		//if err!=nil {
		//	log.Println(err)
		//	return err
		//}

		unluckNumListV := make([]int, 0)
		// 首先, 将预测的结果存放到forecast表格中
		for _, v := range willOrderUnLuckNum {
			err = mysql.StoreResultToForecast2Table(prefix, forecastOrderNum, v.SpecificNum)
			unluckNumListV = append(unluckNumListV, v.SpecificNum)
			if err != nil {
				log.Println(err)
				return err
			}
		}
		// 从原始数据表格中确认是否符合
		err = mysql.DetectForecast2(prefix)
		if err != nil {
			// 这里有错误不可以直接return err, 如果这样, 新的预测数据就发不出去了
			// 比如: 你预测了一个数字, 然后在统计数据中找不到对应的值, (这种情况很常见), 那你的预测数据发不出去了
			log.Println("DetectForecast2: ", err)
			//return err
		}
		msg, err := mysql.StatisticsForecast2(prefix)
		if err != nil {
			log.Println("StatisticsForecast: ", err)
			return err
		}
		message := msg + fmt.Sprintf("%s %s %+v", prefix, forecastOrderNum, unluckNumListV)
		// 尝试三次, 将消息发送出去
		go sendMsgThreeTimes("unluck: \n" + message)
		//return errors.New("send msg fail")
	} else {
		err := mysql.DetectForecast2(prefix)
		if err != nil {
			log.Println("没有预测数值时, DetectForecast2: ", err)
			return err
		}
	}
	return nil
}

// 发送msg三次, 最好使用goroutine, 防止发送端出现延迟影响 main goroutine
func sendMsgThreeTimes(msg string) {
	for i := 0; i < 3; i++ {
		err := pushMsgToBot(msg)
		if err == nil {
			return
		}else {
			log.Printf("send msg fail No.%d\n", i)
		}
	}
}

// 根据输入的orderNum生成下一期orderNum
func nextOne(thisOne string) (string, error) {
	if len(thisOne) != 10 {
		return "", errors.New("error length")
	}
	dateStr := thisOne[:8]
	number := thisOne[8:]
	numberInt, _ := strconv.Atoi(number)

	if numberInt < 42 {
		return dateStr + fmt.Sprintf("%02d", numberInt+1), nil
	} else if numberInt == 42 {
		timeThisOne, err := time.Parse("20060102", dateStr)
		if err != nil {
			return "", err
		}
		return timeThisOne.AddDate(0, 0, 1).Format("20060102") + "01", nil
	} else {
		return "", errors.New("unknown error")
	}
}

func PushMsg(content string) {
	url := "http://sthink.top:8080/pushMsg"
	reqBody, err := json.Marshal(map[string]string{
		"content": content,
	})
	if err != nil {
		log.Printf("make reqBody error \n")
		return
	}
	c := http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Authorization", "eb356e0f-2ad8-4352-9c1c-1af1db033f81")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		log.Printf("%+v\n", err)
		return
	}
	defer resp.Body.Close()

	Bytebody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("%+v\n", err)
		return
	}
	log.Println(string(Bytebody))
}
