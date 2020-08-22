package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/beevik/etree"
	"github.com/go-echarts/go-echarts/charts"
	"go11x5/mysql"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var gdLuck = 1.01 / gdWinOnce
var jxLuck = 1.01 / jxWinOnce
var gdWinOnce = 2.134
var gdUnluck = (5.0 / 11.0) * 0.99
var jxUnluck = (5.0 / 11.0) * 0.99
var jxWinOnce = 2.156

func main() {
	// 测试区域
	////mysql.StatisticsForecast("gd/")
	//msgGd, err := mysql.StatisticsForecast("jx/")
	//if err != nil {
	//	log.Println(err)
	//}
	//fmt.Println(msgGd)

	////立即更新不含今天数据的所有数据到mysql中
	//dateList := constructDate()
	//dateRangeLength := len(dateList) - 1 // 不包含今天
	//for _, v := range dateList[:dateRangeLength] {
	//	saveData2MysqlByDate("gd/", v)
	//	saveData2MysqlByDate("jx/", v)
	//	time.Sleep(time.Second)
	//}

	// 每隔1min 更新最新数据
	//timingGetData()

	////every 10 mins, 筛选符合要求的数学期望 lucky
	//for {
	//	t0 := time.Now()
	//	t0minute := t0.Minute()
	//	t0second := t0.Second()
	//	if t0minute%10 == 3 && t0second == 20 {
	//		log.Println("设定开启的时分到!")
	//		break
	//	}
	//	timer := time.NewTimer(time.Second)
	//	<-timer.C
	//}
	//
	//f, err := os.OpenFile("luckyGd.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//if err != nil {
	//	log.Fatalf("error opening file: %v", err)
	//}
	//defer f.Close()
	//log.SetOutput(f)
	//
	//log.Println("start luckyGd ...")
	//for {
	//	t0 := time.Now()
	//	todayZero := time.Date(t0.Year(), t0.Month(), t0.Day(), 0, 0, 0, 0, time.Local)
	//	// 每天特定时间范围内执行
	//	if t0.After(todayZero.Add(time.Minute*(9*60))) && t0.Before(todayZero.Add(time.Minute*(23*60+30))) && (t0.Minute()%10==3 && t0.Second()==20) {
	//		err := getLuckNum("gd/")
	//		if err != nil {
	//			log.Printf("%+v\n", err)
	//		}
	//	}
	//	// 只阻断 main goroutine, 不阻断其他goroutine的运行
	//	timer := time.NewTimer(time.Second)
	//	<-timer.C
	//	//time.Sleep(10 * time.Minute)
	//}

	//every 10 mins, 筛选符合要求的数学期望 lucky
	for {
		t0 := time.Now()
		t0minute := t0.Minute()
		t0second := t0.Second()
		if t0minute%10 == 3 && t0second == 21 {
			log.Println("设定开启的时分到!")
			break
		}
		timer := time.NewTimer(time.Second)
		<-timer.C
	}
	f, err := os.OpenFile("luckyJx.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	log.Println("start luckyJx ...")
	for {
		t0 := time.Now()
		todayZero := time.Date(t0.Year(), t0.Month(), t0.Day(), 0, 0, 0, 0, time.Local)
		// 每天特定时间范围内执行
		if t0.After(todayZero.Add(time.Minute*(9*60))) && t0.Before(todayZero.Add(time.Minute*(23*60+30))) && (t0.Minute()%10 == 3 && t0.Second() == 21) {
			err := getLuckNum("jx/")
			if err != nil {
				log.Printf("%+v\n", err)
			}
		}
		// 只阻断 main goroutine, 不阻断其他goroutine的运行
		timer := time.NewTimer(time.Second)
		<-timer.C
		//time.Sleep(10 * time.Minute)
	}

	//// 立刻统计forecast_uncd sourluck表格中的数值
	//msgGd, err := mysql.StatisticsForecast2("gd/")
	//if err!=nil {
	//	log.Println(err)
	//}
	//log.Println(msgGd)
	//msgJx, err := mysql.StatisticsForecast2("jx/")
	//if err!=nil {
	//	log.Println(err)
	//}
	//log.Println(msgJx)

	//// 立刻发送探测后的信息给msgbot
	//err := mysql.DetectForecast("jx/")
	//if err != nil {
	//	log.Println(err)
	//}
	//err = mysql.DetectForecast("gd/")
	//if err != nil {
	//	log.Println(err)
	//}
	//// 发送统计信息给bot
	//msgGd, err := mysql.StatisticsForecast("gd/")
	//if err != nil {
	//	log.Println(err)
	//}
	//msgJx, err := mysql.StatisticsForecast("jx/")
	//if err != nil {
	//	log.Println(err)
	//}
	//// 尝试三次
	//for i := 0; i < 3; i++ {
	//	err := pushMsgToBot(msgGd + msgJx)
	//	if err == nil {
	//		time.Sleep(time.Minute) // 如果发送成功, 那就sleep一分钟, 避免重复发送数据
	//		break
	//	}
	//}

	//// 每天23:15 || 11:15拿到匹配结果
	//for {
	//	t0 := time.Now()
	//	if (t0.Hour() == 23 || t0.Hour() == 11) && t0.Minute() == 15 {
	//		err := mysql.DetectForecast("jx/")
	//		if err != nil {
	//			log.Println(err)
	//		}
	//		err = mysql.DetectForecast("gd/")
	//		if err != nil {
	//			log.Println(err)
	//		}
	//		// 发送统计信息给bot
	//		msgGd, err := mysql.StatisticsForecast("gd/")
	//		if err != nil {
	//			log.Println(err)
	//		}
	//		msgJx, err := mysql.StatisticsForecast("jx/")
	//		if err != nil {
	//			log.Println(err)
	//		}
	//		// 尝试三次
	//		for i := 0; i < 3; i++ {
	//			err := pushMsgToBot(msgGd + msgJx)
	//			if err == nil {
	//				time.Sleep(time.Minute) // 如果发送成功, 那就sleep一分钟, 避免重复发送数据
	//				break
	//			}
	//		}
	//	}
	//}

	// 立刻验证当前的猜测
	//err := mysql.DetectForecast("jx/")
	//if err != nil {
	//	log.Println(err)
	//}
	//err = mysql.DetectForecast("gd/")
	//if err != nil {
	//	log.Println(err)
	//}

	//// 发送统计信息给bot
	//msgGd, err := mysql.StatisticsForecast("gd/")
	//if err!=nil {
	//	log.Println(err)
	//}
	//msgJx, err := mysql.StatisticsForecast("jx/")
	//if err!=nil {
	//	log.Println(err)
	//}
	//err = pushMsgToBot(msgGd+msgJx)
	//if err!=nil {
	//	log.Printf("Send msg to bot error\n")
	//}
}

// 定时获取最新数据
func timingGetData() {
	log.Println("start timingGetData ...")
	for {
		t0 := time.Now()
		todayZero := time.Date(t0.Year(), t0.Month(), t0.Day(), 0, 0, 0, 0, time.Local)
		// 每天特定时间范围内执行
		if t0.After(todayZero.Add(time.Minute*(9*60))) && t0.Before(todayZero.Add(time.Minute*(23*60+30))) {
			// 更新最新数据
			go getNewestData("gd/")
			go getNewestData("jx/")
		}
		timer := time.NewTimer(time.Minute)
		<-timer.C
	}
}

// 获取当天单个数字的乐观概率和悲观概率
func getSingleDayProbability(dateStr string, prefix string) error {
	if prefix == "gd/" {
		allData, err := mysql.QueryDataFromMysqlGd()
		if err != nil {
			log.Println(err)
			return err
		}
		lenAllData := len(allData)
		dayData, err := mysql.QueryDataFromMysqlGdSomeDay(dateStr)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println(dayData)
		// 获取今天specificNum出现的次数
		numComeList := make([]int, 12)
		for _, v := range dayData {
			numComeList[1] += v.One
			numComeList[2] += v.Two
			numComeList[3] += v.Three
			numComeList[4] += v.Four
			numComeList[5] += v.Five
			numComeList[6] += v.Six
			numComeList[7] += v.Seven
			numComeList[8] += v.Eight
			numComeList[9] += v.Nine
			numComeList[10] += v.Ten
			numComeList[11] += v.Eleven
		}
		todayHadCome := len(dayData) // 今天已经出现的期数
		allData = allData[:lenAllData-todayHadCome]
		// 获取这个数字最近35天的某一天的最小出现次数
		timeList := constructDate2()
		day := 35
		for i := 1; i < 12; i++ {
			m := showDiffInRange(timeList[len(timeList)-day:], allData, i)
			minNumList := make([]float64, 0)
			for _, v := range m {
				minNumList = append(minNumList, v[3])
			}
			sort.Float64s(minNumList)
			minNum := minNumList[0] // 找出minNumList中的最小值
			//log.Println(minNumList)
			// 乐观和悲观的估计计算数值
			theoryProbability := 42.0 * 5.0 / 11.0
			happyValue := (theoryProbability - float64(numComeList[i])) / (42.0 - float64(todayHadCome))
			unhappyValue := (theoryProbability + minNum - float64(numComeList[i])) / (42.0 - float64(todayHadCome))
			fmt.Printf("%s 数字%d 的乐观估计值是: %.4f, 悲观估计值是: %.4f\n", prefix, i, happyValue, unhappyValue)
		}
		return nil
	}
	if prefix == "jx/" {
		allData, err := mysql.QueryDataFromMysqlJx()
		if err != nil {
			log.Println(err)
			return err
		}
		lenAllData := len(allData)
		dayData, err := mysql.QueryDataFromMysqlJxSomeDay(dateStr)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println(dayData)
		// 获取今天specificNum出现的次数
		numComeList := make([]int, 12)
		for _, v := range dayData {
			numComeList[1] += v.One
			numComeList[2] += v.Two
			numComeList[3] += v.Three
			numComeList[4] += v.Four
			numComeList[5] += v.Five
			numComeList[6] += v.Six
			numComeList[7] += v.Seven
			numComeList[8] += v.Eight
			numComeList[9] += v.Nine
			numComeList[10] += v.Ten
			numComeList[11] += v.Eleven
		}
		log.Println(numComeList)
		todayHadCome := len(dayData) // 今天已经出现的期数
		allData = allData[:lenAllData-todayHadCome]
		// 获取这个数字最近35天的某一天的最小出现次数
		timeList := constructDate2()
		day := 35
		for i := 1; i < 12; i++ {
			m := showDiffInRange(timeList[len(timeList)-day:], allData, i)
			minNumList := make([]float64, 0)
			for _, v := range m {
				minNumList = append(minNumList, v[3])
			}
			sort.Float64s(minNumList)
			minNum := minNumList[0] // 找出minNumList中的最小值
			//log.Println(minNumList)
			// 乐观和悲观的估计计算数值
			theoryProbability := 42.0 * 5.0 / 11.0
			happyValue := (theoryProbability - float64(numComeList[i])) / (42.0 - float64(todayHadCome))
			unhappyValue := (theoryProbability + minNum - float64(numComeList[i])) / (42.0 - float64(todayHadCome))
			fmt.Printf("%s 数字%d 的乐观估计值是: %.4f, 悲观估计值是: %.4f\n", prefix, i, happyValue, unhappyValue)
		}
		return nil
	}
	return errors.New("unsupport prefix")
}

// 实时的获取数据然后写入数据库
func getNewestData(prefix string) error {
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
	// 获取最近的42期数值:结果, 并且将结果打包到 map[string][]string, 例如: "20200710" []string{"01", "02", "03", "10", "11"}
	m := make(map[string][]string, 0)
	iwant := trList[len(trList)-42:]
	for _, tds := range iwant {
		tdsValue := tds.SelectElements("td")
		if len(tdsValue) < 6 {
			log.Fatalln("tds 长度不够")
		}
		m[tdsValue[0].Text()] = []string{tdsValue[1].Text(), tdsValue[2].Text(), tdsValue[3].Text(), tdsValue[4].Text(), tdsValue[5].Text()}
	}
	// 找出最近的42期
	//fmt.Println(m)
	keyList := make([]string, 0)
	for k, _ := range m {
		keyList = append(keyList, k)
	}
	sort.Strings(keyList) // 将keyList排序一下
	for _, k := range keyList {
		//fmt.Printf("%v: %v\n", k, m[k])
		// 开始把这些代码写入数据库
		orderNum := k
		orderValue := stringList2Binlist(m[k])
		orderString := convertValueList2ValueString(orderValue)
		orderTime := orderNumConvertToDate(k)
		// 对数据需要做进一步的校验
		if len(k) != 10 || len(orderValue) != 12 {
			log.Println("数据长度不合法!")
			return errors.New("数据长度不合法!")
		}
		// 遍历查找数据库中是否存在重复数据
		isExist, err := mysql.IsExistInMysql(prefix, orderNum)
		if err != nil {
			log.Fatalln(err)
		}
		if isExist == 0 {
			// 数据不存在于mysql, 可以写入
			if prefix == "jx/" {
				err = mysql.SaveResultToMysqlJx(orderNum, orderTime, orderValue, orderString)
				if err != nil {
					log.Println(err)
					return err
				}
			}
			if prefix == "gd/" {
				err = mysql.SaveResultToMysqlGd(orderNum, orderTime, orderValue, orderString)
				if err != nil {
					log.Println(err)
					return err
				}
			}
		}
	}
	return nil
}

// []string{"01", "02", "03", "04", "05"} -> []int{0, 1, 1, 1,1,1,0,0,0,0,0,0}
func stringList2Binlist(strList []string) []int {
	binList := make([]int, 12)
	for _, v := range strList {
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalln("strconv.Atoi error", err)
		}
		binList[i] += 1
	}
	return binList
}

func httpServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/show", showDiffInRangeHTTP)
	log.Println("Listening on 127.0.0.1:4000")
	err := http.ListenAndServe("127.0.0.1:4000", mux)
	if err != nil {
		log.Println(err)
		return
	}
}

//
// 展示一个时间range内某一个特定数字的理论概率和实际概率之间的差值
func showDiffInRangeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("解析form时 ", err)
		return
	}
	key := r.Form.Get("key")
	day, _ := strconv.Atoi(r.Form.Get("day"))
	if day == 0 {
		day = len(constructDate2())
	}
	if day > len(constructDate2()) || day < 0 {
		log.Println("day的数值不符合我们的要求")
		http.Error(w, "day的数值不符合我们的要求", 400)
		return
	}
	num := r.Form.Get("num")
	if num == "" {
		log.Println("num参数不可以为空")
		return
	}
	specificNum, err := strconv.Atoi(num)
	if err != nil {
		log.Println("str to int error", err)
		return
	}
	log.Println("输入的参数是: ", specificNum, key, day)
	var data []mysql.QueryData
	if key == "jx" {
		data, err = mysql.QueryDataFromMysqlJx()
		if err != nil {
			log.Fatalln(err)
		}
	} else if key == "gd" {
		data, err = mysql.QueryDataFromMysqlGd()
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Println("key错误!")
	}
	timeList := constructDate2()
	m := showDiffInRange(timeList[len(timeList)-day:], data, specificNum)
	xList := make([]string, 0)
	yList := make([]float64, 0)
	for k, _ := range m {
		xList = append(xList, k)
	}
	sort.Strings(xList)   // 对日期字符串进行排序
	var expectNum float64 // 下一日为了达到平衡, expectNum的数值
	for _, v := range xList {
		expectNum += m[v][3]
		yList = append(yList, m[v][3])
	}
	prefixMap := make(map[string]string)
	prefixMap["jx"] = "江西"
	prefixMap["gd"] = "广东"
	// 整理eCharts需要的数据类型
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.TitleOpts{Title: "概率分布"})
	bar.AddXAxis(xList).AddYAxis(prefixMap[key]+"数字:"+strconv.Itoa(specificNum)+" 总体走势: "+fmt.Sprintf("%.4f", expectNum), yList)
	f, err := os.Create("bar.html")
	if err != nil {
		log.Println(err)
	}
	bar.Render(w, f)
}

// 统计一个时间range范围内的某一个特定数字的理论概率和实际概率之间的差值
// 已出现次数, 未出现次数, 与理论值的差(百分比), 与理论值的数字差
func showDiffInRange(timeRange []string, queryData []mysql.QueryData, specificNum int) map[string][]float64 {
	m := make(map[string][]float64) // 记录每一天中单个数字出现的次数
	switch specificNum {
	case 1:
		for _, t := range timeRange {
			count := 0
			allNum := 0
			for _, v := range queryData {
				if v.OrderNum[:8] == t {
					allNum += 1
					if v.One == 1 {
						count += 1
					}
				}
			}
			rate := (float64(count) - 42*5.0/11) * 100 / 42
			diffNum := float64(count) - 42.0*5.0/11.0
			if count != 0 && allNum != 0 {
				m[t] = []float64{float64(count), float64(allNum) - float64(count), math.Trunc(rate*10000+0.5) / 10000, math.Trunc(diffNum*100+0.5) / 100}
			}
		}
		return m
	case 2:
		for _, t := range timeRange {
			count := 0
			allNum := 0
			for _, v := range queryData {
				if v.OrderNum[:8] == t {
					allNum += 1
					if v.Two == 1 {
						count += 1
					}
				}
			}
			rate := (float64(count) - 42*5.0/11) * 100 / 42
			diffNum := float64(count) - 42.0*5.0/11.0
			if count != 0 && allNum != 0 {
				m[t] = []float64{float64(count), float64(allNum) - float64(count), math.Trunc(rate*10000+0.5) / 10000, math.Trunc(diffNum*100+0.5) / 100}
			}
		}
		return m
	case 3:
		for _, t := range timeRange {
			count := 0
			allNum := 0
			for _, v := range queryData {
				if v.OrderNum[:8] == t {
					allNum += 1
					if v.Three == 1 {
						count += 1
					}
				}
			}
			rate := (float64(count) - 42*5.0/11) * 100 / 42
			diffNum := float64(count) - 42.0*5.0/11.0
			if count != 0 && allNum != 0 {
				m[t] = []float64{float64(count), float64(allNum) - float64(count), math.Trunc(rate*10000+0.5) / 10000, math.Trunc(diffNum*100+0.5) / 100}
			}
		}
		return m
	case 4:
		for _, t := range timeRange {
			count := 0
			allNum := 0
			for _, v := range queryData {
				if v.OrderNum[:8] == t {
					allNum += 1
					if v.Four == 1 {
						count += 1
					}
				}
			}
			rate := (float64(count) - 42*5.0/11) * 100 / 42
			diffNum := float64(count) - 42.0*5.0/11.0
			if count != 0 && allNum != 0 {
				m[t] = []float64{float64(count), float64(allNum) - float64(count), math.Trunc(rate*10000+0.5) / 10000, math.Trunc(diffNum*100+0.5) / 100}
			}
		}
		return m
	case 5:
		for _, t := range timeRange {
			count := 0
			allNum := 0
			for _, v := range queryData {
				if v.OrderNum[:8] == t {
					allNum += 1
					if v.Five == 1 {
						count += 1
					}
				}
			}
			rate := (float64(count) - 42*5.0/11) * 100 / 42
			diffNum := float64(count) - 42.0*5.0/11.0
			if count != 0 && allNum != 0 {
				m[t] = []float64{float64(count), float64(allNum) - float64(count), math.Trunc(rate*10000+0.5) / 10000, math.Trunc(diffNum*100+0.5) / 100}
			}
		}
		return m
	case 6:
		for _, t := range timeRange {
			count := 0
			allNum := 0
			for _, v := range queryData {
				if v.OrderNum[:8] == t {
					allNum += 1
					if v.Six == 1 {
						count += 1
					}
				}
			}
			rate := (float64(count) - 42*5.0/11) * 100 / 42
			diffNum := float64(count) - 42.0*5.0/11.0
			if count != 0 && allNum != 0 {
				m[t] = []float64{float64(count), float64(allNum) - float64(count), math.Trunc(rate*10000+0.5) / 10000, math.Trunc(diffNum*100+0.5) / 100}
			}
		}
		return m
	case 7:
		for _, t := range timeRange {
			count := 0
			allNum := 0
			for _, v := range queryData {
				if v.OrderNum[:8] == t {
					allNum += 1
					if v.Seven == 1 {
						count += 1
					}
				}
			}
			rate := (float64(count) - 42*5.0/11) * 100 / 42
			diffNum := float64(count) - 42.0*5.0/11.0
			if count != 0 && allNum != 0 {
				m[t] = []float64{float64(count), float64(allNum) - float64(count), math.Trunc(rate*10000+0.5) / 10000, math.Trunc(diffNum*100+0.5) / 100}
			}
		}
		return m
	case 8:
		for _, t := range timeRange {
			count := 0
			allNum := 0
			for _, v := range queryData {
				if v.OrderNum[:8] == t {
					allNum += 1
					if v.Eight == 1 {
						count += 1
					}
				}
			}
			rate := (float64(count) - 42*5.0/11) * 100 / 42
			diffNum := float64(count) - 42.0*5.0/11.0
			if count != 0 && allNum != 0 {
				m[t] = []float64{float64(count), float64(allNum) - float64(count), math.Trunc(rate*10000+0.5) / 10000, math.Trunc(diffNum*100+0.5) / 100}
			}
		}
		return m
	case 9:
		for _, t := range timeRange {
			count := 0
			allNum := 0
			for _, v := range queryData {
				if v.OrderNum[:8] == t {
					allNum += 1
					if v.Nine == 1 {
						count += 1
					}
				}
			}
			rate := (float64(count) - 42*5.0/11) * 100 / 42
			diffNum := float64(count) - 42.0*5.0/11.0
			if count != 0 && allNum != 0 {
				m[t] = []float64{float64(count), float64(allNum) - float64(count), math.Trunc(rate*10000+0.5) / 10000, math.Trunc(diffNum*100+0.5) / 100}
			}
		}
		return m
	case 10:
		for _, t := range timeRange {
			count := 0
			allNum := 0
			for _, v := range queryData {
				if v.OrderNum[:8] == t {
					allNum += 1
					if v.Ten == 1 {
						count += 1
					}
				}
			}
			rate := (float64(count) - 42*5.0/11) * 100 / 42
			diffNum := float64(count) - 42.0*5.0/11.0
			if count != 0 && allNum != 0 {
				m[t] = []float64{float64(count), float64(allNum) - float64(count), math.Trunc(rate*10000+0.5) / 10000, math.Trunc(diffNum*100+0.5) / 100}
			}
		}
		return m
	case 11:
		for _, t := range timeRange {
			count := 0
			allNum := 0
			for _, v := range queryData {
				if v.OrderNum[:8] == t {
					allNum += 1
					if v.Eleven == 1 {
						count += 1
					}
				}
			}
			rate := (float64(count) - 42*5.0/11) * 100 / 42
			diffNum := float64(count) - 42.0*5.0/11.0
			if count != 0 && allNum != 0 {
				m[t] = []float64{float64(count), float64(allNum) - float64(count), math.Trunc(rate*10000+0.5) / 10000, math.Trunc(diffNum*100+0.5) / 100}
			}
		}
		return m
	}
	log.Fatalln("输入的数字有误")
	return m
}

// 展示思考后的结果, flag为: "jx/" || "gd/", 并且将思考后的结果写入: ["gd", "jx"]_luck表中
// limit 最近多少期
// 展示思考后的结果, flag为: "jx/" || "gd/", 并且将思考后的结果写入: ["gd", "jx"]_unluck表中
func showThink(flag string, limit int) {
	if flag == "gd/" {
		//fmt.Println("开始评估gd近", limit, "期...")
		data, err := mysql.QueryDataFromMysqlGd()
		if err != nil {
			log.Fatalln(err)
		}
		if limit != 0 {
			data = data[len(data)-limit:]
		}

		// 统计每个数字出现的概率
		ariseTimesMap := make(map[int]float64)
		for i := 1; i < 12; i++ {
			// ariseTimesMap, 天才的想法呀, 因为map是指针型的, 所以函数内部对它的修改也会反映到外部
			countTimesArise(data, i, ariseTimesMap)
		}
		//log.Println(ariseTimesMap)

		// 统计遗漏值和时间期数的关系
		//for i := 1; i < 12; i++ {
		//	r := calSpecificNumTimes(data, i)
		//	//fmt.Println("specificNum: ",i, r)
		//	calLeaveAndTimes(r, i)
		//}
		// 统计数字q出现遗漏值大于等于x的遗漏数值和对应期数
		// 每次向luck表中遍历添加数值时, 先清除luck表中的全部内容
		err = mysql.DeleteLuckTable(flag)
		if err != nil {
			log.Printf("清除luck %s 失败\n", flag)
			return
		}
		// 每次向unluck表中遍历添加数值时, 先清除unluck表中的全部内容
		err = mysql.DeleteUnLuckTable(flag)
		if err != nil {
			log.Printf("清除unluck %s 失败\n", flag)
			return
		}
		// 从数字 1 to 11, 将每个数字的遗漏比例统计一下
		for q := 1; q < 12; q++ {
			// 将 某一个特定数字的遗漏数字和遗漏期数拿出来, 并且将遗漏数字排序
			r := calSpecificNumTimes(data, q)
			keyList := make([]int, 0)
			for k, _ := range r {
				keyList = append(keyList, k)
			}
			sort.Ints(keyList)
			//fmt.Printf("%s: specific Num %d\n", flag, q)
			for i := 0; i < len(keyList); i++ {
				if keyList[i] >= 0 {
					count := 0
					for k, v := range r {
						// 这个遗漏值到了之后没有停止
						if k >= keyList[i] {
							count += len(v)
						}
					}
					// 这个数到了之后就停止/这个数到了之后没有停止 的 比值
					if (float64(len(r[keyList[i]]))/float64(count))*(1-ariseTimesMap[q]) > gdLuck {
						//fmt.Printf("遗漏值: %d, 遗漏期数: %d, 就此终止几率: %.4f, 期望收益: %.4f\n", keyList[i], len(r[keyList[i]]), float64(len(r[keyList[i]]))/float64(count), 2.156*float64(len(r[keyList[i]]))/float64(count))
						// 将数学期望大于hopeWin的数值写入数据库
						err := mysql.Write2Luck(flag, q, keyList[i], (float64(len(r[keyList[i]]))/float64(count))/(1+ariseTimesMap[q]), gdWinOnce*((float64(len(r[keyList[i]]))/float64(count))/(1+ariseTimesMap[q])))
						if err != nil {
							log.Printf("write to mysql_luck error\n")
							return
						}
					}

					// 将低于预期的数值写入unluck表格中
					if (float64(len(r[keyList[i]]))/float64(count))*(1-ariseTimesMap[q]) < gdUnluck {
						//fmt.Printf("遗漏值: %d, 遗漏期数: %d, 就此终止几率: %.4f, 期望收益: %.4f\n", keyList[i], len(r[keyList[i]]), float64(len(r[keyList[i]]))/float64(count), 2.156*float64(len(r[keyList[i]]))/float64(count))
						// 将数学期望大于hopeWin的数值写入数据库
						err := mysql.Write2UnLuck(flag, q, keyList[i], (float64(len(r[keyList[i]]))/float64(count))/(1+ariseTimesMap[q]), gdWinOnce*((float64(len(r[keyList[i]]))/float64(count))/(1+ariseTimesMap[q])))
						if err != nil {
							log.Printf("write to mysql_unluck error\n")
							return
						}
					}
				}
			}
		}
		return
	}
	if flag == "jx/" {
		//fmt.Println("开始评估jx近", limit, "期...")
		data, err := mysql.QueryDataFromMysqlJx()
		if err != nil {
			log.Fatalln(err)
		}
		if limit != 0 {
			data = data[len(data)-limit:]
		}

		// 统计每个数字出现的概率
		ariseTimesMap := make(map[int]float64)
		for i := 1; i < 12; i++ {
			countTimesArise(data, i, ariseTimesMap)
		}
		//log.Println(ariseTimesMap)

		// 统计遗漏值和时间期数的关系
		//for i := 1; i < 12; i++ {
		//	r := calSpecificNumTimes(data, i)
		//	//fmt.Println("specificNum: ",i, r)
		//	calLeaveAndTimes(r, i)
		//}
		// 统计数字1出现遗漏值大于等于1的遗漏数值和对应期数
		// 每次向luck表中遍历添加数值时, 先清除luck表中的全部内容
		err = mysql.DeleteLuckTable(flag)
		if err != nil {
			log.Printf("清除luck %s 失败\n", flag)
			return
		}
		err = mysql.DeleteUnLuckTable(flag)
		if err != nil {
			log.Printf("清除unluck %s 失败\n", flag)
			return
		}

		for q := 1; q < 12; q++ {
			r := calSpecificNumTimes(data, q)
			keyList := make([]int, 0)
			for k, _ := range r {
				keyList = append(keyList, k)
			}
			sort.Ints(keyList)
			//fmt.Printf("%s: specific Num %d\n", flag, q)
			for i := 0; i < len(keyList); i++ {
				if keyList[i] >= 0 {
					count := 0
					for k, v := range r {
						if k >= keyList[i] {
							count += len(v)
						}
					}
					if (float64(len(r[keyList[i]]))/float64(count))*(1-ariseTimesMap[q]) > jxLuck {
						//fmt.Printf("遗漏值: %d, 遗漏期数: %d, 就此终止几率: %.4f, 期望收益: %.4f\n", keyList[i], len(r[keyList[i]]), float64(len(r[keyList[i]]))/float64(count), 2.156*float64(len(r[keyList[i]]))/float64(count))
						// 将数学期望大于hopeWin的数值写入数据库
						err := mysql.Write2Luck(flag, q, keyList[i], (float64(len(r[keyList[i]]))/float64(count))/(1+ariseTimesMap[q]), jxWinOnce*((float64(len(r[keyList[i]]))/float64(count))/(1+ariseTimesMap[q])))
						if err != nil {
							log.Printf("write to mysql_luck error\n")
							return
						}
					}
					if (float64(len(r[keyList[i]]))/float64(count))*(1-ariseTimesMap[q]) < jxUnluck {
						//fmt.Printf("遗漏值: %d, 遗漏期数: %d, 就此终止几率: %.4f, 期望收益: %.4f\n", keyList[i], len(r[keyList[i]]), float64(len(r[keyList[i]]))/float64(count), 2.156*float64(len(r[keyList[i]]))/float64(count))
						// 将数学期望大于hopeWin的数值写入数据库
						err := mysql.Write2UnLuck(flag, q, keyList[i], (float64(len(r[keyList[i]]))/float64(count))/(1+ariseTimesMap[q]), jxWinOnce*((float64(len(r[keyList[i]]))/float64(count))/(1+ariseTimesMap[q])))
						if err != nil {
							log.Printf("write to mysql_unluck error\n")
							return
						}
					}
				}
			}
		}
		return
	}
	log.Fatalln("flag暂不支持")
}

// 期数号码对应日期时间, 例如: 2020010101 20200101 09:30
func orderNumConvertToDate(orderNum string) string {
	if len(orderNum) != 10 {
		log.Fatalln("convert 2020010101 20200101 09:30 error", orderNum)
	}
	n, err := strconv.Atoi(orderNum[len(orderNum)-2:])
	if err != nil {
		log.Fatalln(err)
	}
	startMins := 9*60 + 10 // 09:10 一般开始第一期
	endMins := startMins + 20*n
	endHour := endMins / 60
	endHourStr := fmt.Sprintf("%02s", strconv.Itoa(endHour))
	endMin := endMins - endHour*60
	endMinStr := fmt.Sprintf("%2s", strconv.Itoa(endMin))
	return orderNum[:8] + fmt.Sprintf(" %s:%s", endHourStr, endMinStr)
}

// 统计某个数字当天出现的频率与标准值之间的偏差
func countInaccurateOneDay(queryData []mysql.QueryData, specificNum int, dateStr string) {
	allTimes := len(queryData)
	switch specificNum {
	case 1:
		countTimes := 0
		for i := 0; i < allTimes; i++ {
			if queryData[i].OrderNum[:8] == dateStr {
				if queryData[i].One == 1 {
					countTimes += 1
				}
			}
		}
		fmt.Printf("数字%d 出现的次数是: %d, 概率是: %.4f, 与标准值偏差: %.4f%%\n", 1, countTimes, float64(countTimes)/42.0, (float64(countTimes)-42.0*5/11)*100/42.0)
	case 2:
		countTimes := 0
		for i := 0; i < allTimes; i++ {
			if queryData[i].OrderNum[:8] == dateStr {
				if queryData[i].Two == 1 {
					countTimes += 1
				}
			}
		}
		fmt.Printf("数字%d 出现的次数是: %d, 概率是: %.4f, 与标准值偏差: %.4f%%\n", 2, countTimes, float64(countTimes)/42.0, (float64(countTimes)-42.0*5/11)*100/42.0)

	case 3:
		countTimes := 0
		for i := 0; i < allTimes; i++ {
			if queryData[i].OrderNum[:8] == dateStr {
				if queryData[i].Three == 1 {
					countTimes += 1
				}
			}
		}
		fmt.Printf("数字%d 出现的次数是: %d, 概率是: %.4f, 与标准值偏差: %.4f%%\n", 3, countTimes, float64(countTimes)/42.0, (float64(countTimes)-42.0*5/11)*100/42.0)

	case 4:
		countTimes := 0
		for i := 0; i < allTimes; i++ {
			if queryData[i].OrderNum[:8] == dateStr {
				if queryData[i].Four == 1 {
					countTimes += 1
				}
			}
		}
		fmt.Printf("数字%d 出现的次数是: %d, 概率是: %.4f, 与标准值偏差: %.4f%%\n", 4, countTimes, float64(countTimes)/42.0, (float64(countTimes)-42.0*5/11)*100/42.0)

	case 5:
		countTimes := 0
		for i := 0; i < allTimes; i++ {
			if queryData[i].OrderNum[:8] == dateStr {
				if queryData[i].Five == 1 {
					countTimes += 1
				}
			}
		}
		fmt.Printf("数字%d 出现的次数是: %d, 概率是: %.4f, 与标准值偏差: %.4f%%\n", 5, countTimes, float64(countTimes)/42.0, (float64(countTimes)-42.0*5/11)*100/42.0)

	case 6:
		countTimes := 0
		for i := 0; i < allTimes; i++ {
			if queryData[i].OrderNum[:8] == dateStr {
				if queryData[i].Six == 1 {
					countTimes += 1
				}
			}
		}
		fmt.Printf("数字%d 出现的次数是: %d, 概率是: %.4f, 与标准值偏差: %.4f%%\n", 6, countTimes, float64(countTimes)/42.0, (float64(countTimes)-42.0*5/11)*100/42.0)

	case 7:
		countTimes := 0
		for i := 0; i < allTimes; i++ {
			if queryData[i].OrderNum[:8] == dateStr {
				if queryData[i].Seven == 1 {
					countTimes += 1
				}
			}
		}
		fmt.Printf("数字%d 出现的次数是: %d, 概率是: %.4f, 与标准值偏差: %.4f%%\n", 7, countTimes, float64(countTimes)/42.0, (float64(countTimes)-42.0*5/11)*100/42.0)

	case 8:
		countTimes := 0
		for i := 0; i < allTimes; i++ {
			if queryData[i].OrderNum[:8] == dateStr {
				if queryData[i].Eight == 1 {
					countTimes += 1
				}
			}
		}
		fmt.Printf("数字%d 出现的次数是: %d, 概率是: %.4f, 与标准值偏差: %.4f%%\n", 8, countTimes, float64(countTimes)/42.0, (float64(countTimes)-42.0*5/11)*100/42.0)

	case 9:
		countTimes := 0
		for i := 0; i < allTimes; i++ {
			if queryData[i].OrderNum[:8] == dateStr {
				if queryData[i].Nine == 1 {
					countTimes += 1
				}
			}
		}
		fmt.Printf("数字%d 出现的次数是: %d, 概率是: %.4f, 与标准值偏差: %.4f%%\n", 9, countTimes, float64(countTimes)/42.0, (float64(countTimes)-42.0*5/11)*100/42.0)

	case 10:
		countTimes := 0
		for i := 0; i < allTimes; i++ {
			if queryData[i].OrderNum[:8] == dateStr {
				if queryData[i].Ten == 1 {
					countTimes += 1
				}
			}
		}
		fmt.Printf("数字%d 出现的次数是: %d, 概率是: %.4f, 与标准值偏差: %.4f%%\n", 10, countTimes, float64(countTimes)/42.0, (float64(countTimes)-42.0*5/11)*100/42.0)

	case 11:
		countTimes := 0
		for i := 0; i < allTimes; i++ {
			if queryData[i].OrderNum[:8] == dateStr {
				if queryData[i].Eleven == 1 {
					countTimes += 1
				}
			}
		}
		fmt.Printf("数字%d 出现的次数是: %d, 概率是: %.4f, 与标准值偏差: %.4f%%\n", 11, countTimes, float64(countTimes)/42.0, (float64(countTimes)-42.0*5/11)*100/42.0)

	default:
		log.Fatalln("未知错误")
	}

}

// 统计所有数据中某个值出现的次数, 未出现的次数, 与理论值的偏差
func countTimesArise(queryData []mysql.QueryData, specificNum int, ariseTimesMap map[int]float64) {
	allTimes := len(queryData)
	ariseTimes := 0 // 统计出现次数
	switch specificNum {
	case 1:
		for i := 0; i < allTimes; i++ {
			if queryData[i].One == 1 {
				ariseTimes += 1
			}
		}
		log.Printf("数字%d", specificNum)
		log.Printf("出现次数: %d, 占比: %.4f, 未出现次数: %d, 占比: %.4f, 总次数: %d, 与理论值偏差: %.4f%%\n", ariseTimes, float64(ariseTimes)/float64(allTimes), allTimes-ariseTimes, float64(allTimes-ariseTimes)/float64(allTimes), allTimes, (float64(ariseTimes)/float64(allTimes)-5.0/11.0)*100/(5.0/11.0))
		ariseTimesMap[1] = (float64(ariseTimes)/float64(allTimes) - 5.0/11.0) / (5.0 / 11.0)
	case 2:
		for i := 0; i < allTimes; i++ {
			if queryData[i].Two == 1 {
				ariseTimes += 1
			}
		}
		log.Printf("数字%d", specificNum)
		log.Printf("出现次数: %d, 占比: %.4f, 未出现次数: %d, 占比: %.4f, 总次数: %d, 与理论值偏差: %.4f%%\n", ariseTimes, float64(ariseTimes)/float64(allTimes), allTimes-ariseTimes, float64(allTimes-ariseTimes)/float64(allTimes), allTimes, (float64(ariseTimes)/float64(allTimes)-5.0/11.0)*100/(5.0/11.0))
		ariseTimesMap[2] = (float64(ariseTimes)/float64(allTimes) - 5.0/11.0) / (5.0 / 11.0)
	case 3:
		for i := 0; i < allTimes; i++ {
			if queryData[i].Three == 1 {
				ariseTimes += 1
			}
		}
		log.Printf("数字%d", specificNum)
		log.Printf("出现次数: %d, 占比: %.4f, 未出现次数: %d, 占比: %.4f, 总次数: %d, 与理论值偏差: %.4f%%\n", ariseTimes, float64(ariseTimes)/float64(allTimes), allTimes-ariseTimes, float64(allTimes-ariseTimes)/float64(allTimes), allTimes, (float64(ariseTimes)/float64(allTimes)-5.0/11.0)*100/(5.0/11.0))
		ariseTimesMap[3] = (float64(ariseTimes)/float64(allTimes) - 5.0/11.0) / (5.0 / 11.0)
	case 4:
		for i := 0; i < allTimes; i++ {
			if queryData[i].Four == 1 {
				ariseTimes += 1
			}
		}
		log.Printf("数字%d", specificNum)
		log.Printf("出现次数: %d, 占比: %.4f, 未出现次数: %d, 占比: %.4f, 总次数: %d, 与理论值偏差: %.4f%%\n", ariseTimes, float64(ariseTimes)/float64(allTimes), allTimes-ariseTimes, float64(allTimes-ariseTimes)/float64(allTimes), allTimes, (float64(ariseTimes)/float64(allTimes)-5.0/11.0)*100/(5.0/11.0))
		ariseTimesMap[4] = (float64(ariseTimes)/float64(allTimes) - 5.0/11.0) / (5.0 / 11.0)

	case 5:
		for i := 0; i < allTimes; i++ {
			if queryData[i].Five == 1 {
				ariseTimes += 1
			}
		}
		log.Printf("数字%d", specificNum)
		log.Printf("出现次数: %d, 占比: %.4f, 未出现次数: %d, 占比: %.4f, 总次数: %d, 与理论值偏差: %.4f%%\n", ariseTimes, float64(ariseTimes)/float64(allTimes), allTimes-ariseTimes, float64(allTimes-ariseTimes)/float64(allTimes), allTimes, (float64(ariseTimes)/float64(allTimes)-5.0/11.0)*100/(5.0/11.0))
		ariseTimesMap[5] = (float64(ariseTimes)/float64(allTimes) - 5.0/11.0) / (5.0 / 11.0)

	case 6:
		for i := 0; i < allTimes; i++ {
			if queryData[i].Six == 1 {
				ariseTimes += 1
			}
		}
		log.Printf("数字%d", specificNum)
		log.Printf("出现次数: %d, 占比: %.4f, 未出现次数: %d, 占比: %.4f, 总次数: %d, 与理论值偏差: %.4f%%\n", ariseTimes, float64(ariseTimes)/float64(allTimes), allTimes-ariseTimes, float64(allTimes-ariseTimes)/float64(allTimes), allTimes, (float64(ariseTimes)/float64(allTimes)-5.0/11.0)*100/(5.0/11.0))
		ariseTimesMap[6] = (float64(ariseTimes)/float64(allTimes) - 5.0/11.0) / (5.0 / 11.0)

	case 7:
		for i := 0; i < allTimes; i++ {
			if queryData[i].Seven == 1 {
				ariseTimes += 1
			}
		}
		log.Printf("数字%d", specificNum)
		log.Printf("出现次数: %d, 占比: %.4f, 未出现次数: %d, 占比: %.4f, 总次数: %d, 与理论值偏差: %.4f%%\n", ariseTimes, float64(ariseTimes)/float64(allTimes), allTimes-ariseTimes, float64(allTimes-ariseTimes)/float64(allTimes), allTimes, (float64(ariseTimes)/float64(allTimes)-5.0/11.0)*100/(5.0/11.0))
		ariseTimesMap[7] = (float64(ariseTimes)/float64(allTimes) - 5.0/11.0) / (5.0 / 11.0)

	case 8:
		for i := 0; i < allTimes; i++ {
			if queryData[i].Eight == 1 {
				ariseTimes += 1
			}
		}
		log.Printf("数字%d", specificNum)
		log.Printf("出现次数: %d, 占比: %.4f, 未出现次数: %d, 占比: %.4f, 总次数: %d, 与理论值偏差: %.4f%%\n", ariseTimes, float64(ariseTimes)/float64(allTimes), allTimes-ariseTimes, float64(allTimes-ariseTimes)/float64(allTimes), allTimes, (float64(ariseTimes)/float64(allTimes)-5.0/11.0)*100/(5.0/11.0))
		ariseTimesMap[8] = (float64(ariseTimes)/float64(allTimes) - 5.0/11.0) / (5.0 / 11.0)

	case 9:
		for i := 0; i < allTimes; i++ {
			if queryData[i].Nine == 1 {
				ariseTimes += 1
			}
		}
		log.Printf("数字%d", specificNum)
		log.Printf("出现次数: %d, 占比: %.4f, 未出现次数: %d, 占比: %.4f, 总次数: %d, 与理论值偏差: %.4f%%\n", ariseTimes, float64(ariseTimes)/float64(allTimes), allTimes-ariseTimes, float64(allTimes-ariseTimes)/float64(allTimes), allTimes, (float64(ariseTimes)/float64(allTimes)-5.0/11.0)*100/(5.0/11.0))
		ariseTimesMap[9] = (float64(ariseTimes)/float64(allTimes) - 5.0/11.0) / (5.0 / 11.0)

	case 10:
		for i := 0; i < allTimes; i++ {
			if queryData[i].Ten == 1 {
				ariseTimes += 1
			}
		}
		log.Printf("数字%d", specificNum)
		log.Printf("出现次数: %d, 占比: %.4f, 未出现次数: %d, 占比: %.4f, 总次数: %d, 与理论值偏差: %.4f%%\n", ariseTimes, float64(ariseTimes)/float64(allTimes), allTimes-ariseTimes, float64(allTimes-ariseTimes)/float64(allTimes), allTimes, (float64(ariseTimes)/float64(allTimes)-5.0/11.0)*100/(5.0/11.0))
		ariseTimesMap[10] = (float64(ariseTimes)/float64(allTimes) - 5.0/11.0) / (5.0 / 11.0)

	case 11:
		for i := 0; i < allTimes; i++ {
			if queryData[i].Eleven == 1 {
				ariseTimes += 1
			}
		}
		log.Printf("数字%d", specificNum)
		log.Printf("出现次数: %d, 占比: %.4f, 未出现次数: %d, 占比: %.4f, 总次数: %d, 与理论值偏差: %.4f%%\n", ariseTimes, float64(ariseTimes)/float64(allTimes), allTimes-ariseTimes, float64(allTimes-ariseTimes)/float64(allTimes), allTimes, (float64(ariseTimes)/float64(allTimes)-5.0/11.0)*100/(5.0/11.0))
		ariseTimesMap[11] = (float64(ariseTimes)/float64(allTimes) - 5.0/11.0) / (5.0 / 11.0)

	default:
		log.Fatalln("未知错误!")
	}
}

// 按照key排序后的结果打印某个数字遗漏的期数和对应的时间
func calLeaveAndTimes(m map[int][]string, specificNum int) {
	// 先把key排序
	keyList := make([]int, 0)
	for k, _ := range m {
		keyList = append(keyList, k)
	}
	sort.Ints(keyList)
	result := make([]string, 0)
	// keyList 最大值, 即最大遗漏值和最大遗漏值对应的期数
	for _, v := range m[keyList[len(keyList)-1]] {
		result = append(result, orderNumConvertToDate(v))
	}
	fmt.Println("数字", specificNum, "遗漏的最大数字是: ", keyList[len(keyList)-1], "遗漏的期数是: ", m[keyList[len(keyList)-1]], "对应的时间是: ", result)
}

// 计算某个数字在queryList中遗漏的每个次数以及对应的orderNum
func calSpecificNumTimes(queryList []mysql.QueryData, specificNum int) map[int][]string {
	count := 0
	countMap := make(map[int][]string, 0)
	// 计算数字 specificNum
	switch specificNum {
	case 1:
		for i := 0; i < len(queryList); i++ {
			if queryList[i].One == 0 {
				count += 1
			} else {
				countMap[count] = append(countMap[count], queryList[i].OrderNum)
				count = 0
			}
		}
	case 2:
		for i := 0; i < len(queryList); i++ {
			if queryList[i].Two == 0 {
				count += 1
			} else {
				countMap[count] = append(countMap[count], queryList[i].OrderNum)
				count = 0
			}
		}
	case 3:
		for i := 0; i < len(queryList); i++ {
			if queryList[i].Three == 0 {
				count += 1
			} else {
				countMap[count] = append(countMap[count], queryList[i].OrderNum)
				count = 0
			}
		}
	case 4:
		for i := 0; i < len(queryList); i++ {
			if queryList[i].Four == 0 {
				count += 1
			} else {
				countMap[count] = append(countMap[count], queryList[i].OrderNum)
				count = 0
			}
		}
	case 5:
		for i := 0; i < len(queryList); i++ {
			if queryList[i].Five == 0 {
				count += 1
			} else {
				countMap[count] = append(countMap[count], queryList[i].OrderNum)
				count = 0
			}
		}
	case 6:
		for i := 0; i < len(queryList); i++ {
			if queryList[i].Six == 0 {
				count += 1
			} else {
				countMap[count] = append(countMap[count], queryList[i].OrderNum)
				count = 0
			}
		}
	case 7:
		for i := 0; i < len(queryList); i++ {
			if queryList[i].Seven == 0 {
				count += 1
			} else {
				countMap[count] = append(countMap[count], queryList[i].OrderNum)
				count = 0
			}
		}
	case 8:
		for i := 0; i < len(queryList); i++ {
			if queryList[i].Eight == 0 {
				count += 1
			} else {
				countMap[count] = append(countMap[count], queryList[i].OrderNum)
				count = 0
			}
		}
	case 9:
		for i := 0; i < len(queryList); i++ {
			if queryList[i].Nine == 0 {
				count += 1
			} else {
				countMap[count] = append(countMap[count], queryList[i].OrderNum)
				count = 0
			}
		}
	case 10:
		for i := 0; i < len(queryList); i++ {
			if queryList[i].Ten == 0 {
				count += 1
			} else {
				countMap[count] = append(countMap[count], queryList[i].OrderNum)
				count = 0
			}
		}
	case 11:
		for i := 0; i < len(queryList); i++ {
			if queryList[i].Eleven == 0 {
				count += 1
			} else {
				countMap[count] = append(countMap[count], queryList[i].OrderNum)
				count = 0
			}
		}
	default:
		log.Fatalln("未知错误")
	}
	return countMap
}

// 使用hash表, 比较两个字符串数组中不同的数
func compare2StringSlice(str1 []string, str2 []string) []string {
	m := make(map[string]int)
	for _, v := range str1 {
		m[v] += 1
	}
	for _, v := range str2 {
		m[v] += 1
	}
	result := make([]string, 0)
	for k, v := range m {
		if v == 1 {
			result = append(result, k)
		}
	}
	sort.Strings(result)
	log.Println(result, len(result))
	return result
}

// prefix 要保存文件的目录, 如 gd/
// dateStr 日期字符串 例如: 2020-07-10
// 每天晚上11:30可以把结果保存在mysql中, 定时任务
func saveData2MysqlByDate(prefix string, dateStr string) {
	baseReqURL := ""
	if prefix == "jx/" {
		baseReqURL = "https://www.polocai.com/hfkj_jx11x5/"
	} else if prefix == "gd/" {
		baseReqURL = "https://www.polocai.com/hfkj_gd11x5/"
	} else {
		log.Fatalln("prefix 不在支持范围内")
	}
	c := http.Client{
		Timeout: 60 * time.Second,
	}
	m := make(map[string][]int, 0)
	req, err := http.NewRequest(http.MethodGet, baseReqURL+dateStr+".html", nil)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	resp, err := c.Do(req)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	r, _ := ioutil.ReadAll(resp.Body)
	// 保存文件到本地硬盘上
	err = ioutil.WriteFile(prefix+strings.ReplaceAll(dateStr, "-", "")+".html", r, 0644)
	if err != nil {
		log.Printf("%s, 保存文件到本地硬盘出错!", err.Error())
	}
	time.Sleep(time.Duration(rand.Intn(1000000001)) * time.Nanosecond)
	//将之间的io.Reader文件读取一下
	//或者直接使用os.Open()
	fileReader, err := os.Open(prefix + strings.ReplaceAll(dateStr, "-", "") + ".html")
	if err != nil {
		log.Fatalln(err)
	}
	doc, err := goquery.NewDocumentFromReader(fileReader)
	if err != nil {
		log.Fatalln(err)
	}
	// datePrefix 在存储orderNumber的时候会作为前缀而用到
	datePrefix := strings.ReplaceAll(dateStr, "-", "")
	content := doc.Find(".middle").Find(".bo_ocs,.bt_ocs_orange,.bg_white,.mt15").Find(".ov,.bx_bb").Find(".w_333,.ov,.bx_bb").Find(".bg_white,.bg_grayf2")
	content.Each(func(i int, selection *goquery.Selection) {
		// 在这里分解num和value
		orderNum := selection.Find(".w_15").Text()
		if len(orderNum) == 6 {
			orderNum = datePrefix + orderNum[len(orderNum)-2:]
			m[orderNum] = orderValue(selection)
		}
	})
	//fmt.Println(m)

	// 准备将数据写入数据库
	keyList := make([]string, 0)
	for k, _ := range m {
		keyList = append(keyList, k)
	}
	sort.Strings(keyList) // 将keyList排序一下
	for _, k := range keyList {
		//fmt.Printf("%v: %v\n", k, m[k])
		// 开始把这些代码写入数据库
		orderNum := k
		orderValue := m[k]
		orderString := convertValueList2ValueString(orderValue)
		orderTime := orderNumConvertToDate(orderNum)
		// 对数据需要做进一步的校验
		if len(k) != 10 || len(m[k]) != 12 {
			continue
		}
		// 根据prefix来决定写入哪个数据库
		if prefix == "jx/" {
			err := mysql.SaveResultToMysqlJx(orderNum, orderTime, orderValue, orderString)
			if err != nil {
				log.Fatalln(orderNum, orderValue, orderString)
			}
		}
		if prefix == "gd/" {
			err := mysql.SaveResultToMysqlGd(orderNum, orderTime, orderValue, orderString)
			if err != nil {
				log.Fatalln(orderNum, orderValue, orderString)
			}
		}
	}
}

// prefix: "jx/" || "gd/"
func saveData2Mysql(prefix string) {
	dateList := constructDate()
	dateRangeLength := len(dateList) - 1 // 不包含今天
	m := make(map[string][]int, 0)
	for i := 0; i < dateRangeLength; i++ {
		fileReader, err := os.Open(prefix + strings.ReplaceAll(dateList[i], "-", "") + ".html")
		if err != nil {
			log.Printf("open file error, %+v \n", err)
			return
		}
		doc, err := goquery.NewDocumentFromReader(fileReader)
		if err != nil {
			log.Printf("goQuery fileReader error, %+v \n", err)
		}
		datePrefix := strings.ReplaceAll(dateList[i], "-", "")
		content := doc.Find(".middle").Find(".bo_ocs,.bt_ocs_orange,.bg_white,.mt15").Find(".ov,.bx_bb").Find(".w_333,.ov,.bx_bb").Find(".bg_white,.bg_grayf2")
		content.Each(func(i int, selection *goquery.Selection) {
			// 在这里分解num和value
			orderNum := selection.Find(".w_15").Text()
			if len(orderNum) == 6 {
				orderNum = datePrefix + orderNum[len(orderNum)-2:]
				m[orderNum] = orderValue(selection)
			}
		})
	}
	keyList := make([]string, 0)
	for k, _ := range m {
		keyList = append(keyList, k)
	}
	sort.Strings(keyList) // 将keyList排序一下
	for _, k := range keyList {
		//fmt.Printf("%v: %v\n", k, m[k])
		// 开始把这些代码写入数据库
		orderNum := k
		orderValue := m[k]
		orderString := convertValueList2ValueString(orderValue)
		orderTime := orderNumConvertToDate(k)
		// 对数据需要做进一步的校验
		if len(k) != 10 || len(m[k]) != 12 {
			log.Println("数据长度不合法!")
			continue
		}
		if prefix == "jx/" {
			err := mysql.SaveResultToMysqlJx(orderNum, orderTime, orderValue, orderString)
			if err != nil {
				log.Fatalln(orderNum, orderValue, orderString)
			}
		}
		if prefix == "gd/" {
			err := mysql.SaveResultToMysqlGd(orderNum, orderTime, orderValue, orderString)
			if err != nil {
				log.Fatalln(orderNum, orderValue, orderString)
			}
		}
	}
}

// 转换: [0 0 0 0 0 1 1 1 1 1 0] to "0607080910"
func convertValueList2ValueString(valueList []int) string {
	r := ""
	for i := 1; i < 12; i++ {
		if valueList[i] == 1 {
			if i < 10 {
				r += "0" + strconv.Itoa(i)
			} else {
				r += strconv.Itoa(i)
			}
		}
	}
	return r
}

// 返回 [0 1 1 0 1 0 0 1 1 0 0 0]
func orderValue(selection *goquery.Selection) []int {
	valueList := make([]string, 0)
	selection.Find(".w_63").Each(func(i int, selection *goquery.Selection) {
		selection.Find(".ds_ib").Each(func(i int, selection *goquery.Selection) {
			valueList = append(valueList, selection.Text())
		})
	})
	binList := make([]int, 12)
	for _, v := range valueList {
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalln("strconv.Atoi error", err)
		}
		binList[i] += 1
	}
	return binList
}

// 设定初始年月日, 每次提高1天, 2019-02-11
func constructDate() []string {
	startDate := time.Date(2019, time.Month(2), 11, 0, 0, 0, 0, time.Local)
	now := time.Now()
	diff := int(now.Sub(startDate).Hours() / 24)
	resultList := make([]string, 0)
	for i := 0; i <= diff; i++ {
		nextDay := startDate.AddDate(0, 0, i)
		nextDayYear := strconv.Itoa(nextDay.Year())
		nextDayMonth := "0" + strconv.Itoa(int(nextDay.Month()))
		nextDayMonth = nextDayMonth[len(nextDayMonth)-2:]
		nextDayDay := "0" + strconv.Itoa(nextDay.Day())
		nextDayDay = nextDayDay[len(nextDayDay)-2:]
		result := nextDayYear + "-" + nextDayMonth + "-" + nextDayDay
		resultList = append(resultList, result)
	}
	//fmt.Println(resultList)
	return resultList
}

// 设定初始年月日, 每次提高1天, 20190211
func constructDate2() []string {
	startDate := time.Date(2019, time.Month(2), 11, 0, 0, 0, 0, time.Local)
	now := time.Now()
	diff := int(now.Sub(startDate).Hours() / 24)
	resultList := make([]string, 0)
	for i := 0; i <= diff; i++ {
		nextDay := startDate.AddDate(0, 0, i)
		nextDayYear := strconv.Itoa(nextDay.Year())
		nextDayMonth := "0" + strconv.Itoa(int(nextDay.Month()))
		nextDayMonth = nextDayMonth[len(nextDayMonth)-2:]
		nextDayDay := "0" + strconv.Itoa(nextDay.Day())
		nextDayDay = nextDayDay[len(nextDayDay)-2:]
		result := nextDayYear + nextDayMonth + nextDayDay
		resultList = append(resultList, result)
	}
	//fmt.Println(resultList)
	return resultList
}
