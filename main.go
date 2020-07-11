package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-echarts/go-echarts/charts"
	"go11x5/mysql"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	//"github.com/go-echarts/go-echarts/charts"
)

func main() {
	//dayList := []int{1, 2, 3, 5, 7, 14, 21, 30, 60, 90, 120, 180, 360}
	//for i := 0; i < len(dayList); i++ {
	//	fmt.Printf("近%d天\n", dayList[i])
	//	showThink("jx/", 42*dayList[i])
	//}

	//// 每天运行一遍定时任务
	saveData2MysqlByDate("jx/", time.Now().Format("2006-01-02"))

	////fmt.Println(orderNumConvertToDate("20200521022"))
	////saveData2MysqlByDate("jx/", "2019-10-01")

	mux := http.NewServeMux()
	mux.HandleFunc("/show", showDiffInRangeHTTP)
	log.Println("Listening on 127.0.0.1:4000")
	err := http.ListenAndServe("127.0.0.1:4000", mux)
	if err != nil {
		log.Println(err)
		return
	}
}

// 展示一个时间range内某一个特定数字的理论概率和实际概率之间的差值
func showDiffInRangeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("解析form时 ", err)
		return
	}
	key := r.Form.Get("key")
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
	log.Println("输入的参数是: ", specificNum, key)
	m := showDiffInRange(constructDate2(), data, specificNum)
	xList := make([]string, 0)
	yList := make([]float64, 0)
	for k, _ := range m {
		xList = append(xList, k)
	}
	sort.Strings(xList)
	for _, v := range xList {
		yList = append(yList, m[v][3])
	}
	prefixMap := make(map[string]string)
	prefixMap["jx"] = "江西"
	prefixMap["gd"] = "广东"
	// 整理eCharts需要的数据类型
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.TitleOpts{Title: "概率分布"})
	bar.AddXAxis(xList).AddYAxis(prefixMap[key]+"数字:"+strconv.Itoa(specificNum), yList)
	f, err := os.Create("bar.html")
	if err != nil {
		log.Println(err)
	}
	bar.Render(w, f)
}

// 统计一个时间range范围内的某一个特定数字的理论概率和实际概率之间的差值
// 已出现次数, 未出现次数, 与理论值的差(百分比), 与理论值的数字差
func showDiffInRange(timeRange []string, queryData []mysql.QueryData, specificNum int) map[string][]float64 {
	timeRange = constructDate2()
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
				m[t] = []float64{float64(count), float64(allNum) - float64(count), math.Trunc(rate*10000+0.5) / 10000, math.Trunc(diffNum*100 + 0.5) / 100}
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
				m[t] = []float64{float64(count), float64(allNum) - float64(count), math.Trunc(rate*10000+0.5) / 10000, math.Trunc(diffNum*100 + 0.5) / 100}
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
				m[t] = []float64{float64(count), float64(allNum) - float64(count), math.Trunc(rate*10000+0.5) / 10000, math.Trunc(diffNum*100 + 0.5) / 100}
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
				m[t] = []float64{float64(count), float64(allNum) - float64(count), math.Trunc(rate*10000+0.5) / 10000, math.Trunc(diffNum*100 + 0.5) / 100}
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
				m[t] = []float64{float64(count), float64(allNum) - float64(count), math.Trunc(rate*10000+0.5) / 10000, math.Trunc(diffNum*100 + 0.5) / 100}
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
				m[t] = []float64{float64(count), float64(allNum) - float64(count), math.Trunc(rate*10000+0.5) / 10000, math.Trunc(diffNum*100 + 0.5) / 100}
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
				m[t] = []float64{float64(count), float64(allNum) - float64(count), math.Trunc(rate*10000+0.5) / 10000, math.Trunc(diffNum*100 + 0.5) / 100}
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
				m[t] = []float64{float64(count), float64(allNum) - float64(count), math.Trunc(rate*10000+0.5) / 10000, math.Trunc(diffNum*100 + 0.5) / 100}
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
				m[t] = []float64{float64(count), float64(allNum) - float64(count), math.Trunc(rate*10000+0.5) / 10000, math.Trunc(diffNum*100 + 0.5) / 100}
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
				m[t] = []float64{float64(count), float64(allNum) - float64(count), math.Trunc(rate*10000+0.5) / 10000, math.Trunc(diffNum*100 + 0.5) / 100}
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
				m[t] = []float64{float64(count), float64(allNum) - float64(count), math.Trunc(rate*10000+0.5) / 10000, math.Trunc(diffNum*100 + 0.5) / 100}
			}
		}
		return m
	}
	log.Fatalln("输入的数字有误")
	return m
}

// 展示思考后的结果, flag为: "jx/" || "gd/"
// limit 最近多少期, 也可以自行改为最近多少天
func showThink(flag string, limit int) {
	if flag == "gd/" {
		fmt.Println("开始评估gd近", limit, "期...")
		data, err := mysql.QueryDataFromMysqlGd()
		if err != nil {
			log.Fatalln(err)
		}
		if limit != 0 {
			data = data[len(data)-limit:]
		}
		// 统计每个数字出现的概率
		for i := 1; i < 12; i++ {
			//r := calSpecificNumTimes(data, i)
			countTimesArise(data, i)
		}

		// 统计遗漏值和时间期数的关系
		for i := 1; i < 12; i++ {
			r := calSpecificNumTimes(data, i)
			calLeaveAndTimes(r, i)
		}
		// 统计单日遗漏值与理论之间的差距
		//for i := 1; i < 12; i++ {
		//	countOffsetCompareWithTheory(data, i, "20200710")
		//}
		return
	}
	if flag == "jx/" {
		fmt.Println("开始评估jx...")
		data, err := mysql.QueryDataFromMysqlJx()
		if err != nil {
			log.Fatalln(err)
		}
		if limit != 0 {
			data = data[len(data)-limit:]
		}
		// 统计每个数字出现的概率
		for i := 1; i < 12; i++ {
			//r := calSpecificNumTimes(data, i)
			countTimesArise(data, i)
		}

		// 统计遗漏值和时间期数的关系
		for i := 1; i < 12; i++ {
			r := calSpecificNumTimes(data, i)
			calLeaveAndTimes(r, i)
		}
		// 统计单日遗漏值与理论之间的差距
		//for i := 1; i < 12; i++ {
		//	fmt.Printf("开始统计 %s\n", "20200710")
		//	countOffsetCompareWithTheory(data, i, "20200710")
		//}
		return
	}
	log.Fatalln("flag暂不支持")
}

// 期数号码对应日期时间, 例如: 20200101001 20200101 08:50
func orderNumConvertToDate(orderNum string) string {
	if len(orderNum) != 11 {
		log.Fatalln(orderNum)
	}
	n, err := strconv.Atoi(orderNum[len(orderNum)-3:])
	if err != nil {
		log.Fatalln(err)
	}
	startMins := 8*60 + 30
	endMins := startMins + 20*n
	endHour := endMins / 60
	endHourStr := fmt.Sprintf("%02s", strconv.Itoa(endHour))
	endMin := endMins - endHour*60
	endMinStr := fmt.Sprintf("%2s", strconv.Itoa(endMin))
	return orderNum[:8] + fmt.Sprintf(" %s:%s", endHourStr, endMinStr)
}

// 统计某个数字当天出现的频率与标准值之间的偏差
func countOffsetCompareWithTheory(queryData []mysql.QueryData, specificNum int, dateStr string) {
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

// 统计所有次数中出现值, 未出现值
func countTimesArise(queryData []mysql.QueryData, specificNum int) {
	allTimes := len(queryData)
	ariseTimes := 0 // 统计出现次数
	switch specificNum {
	case 1:
		for i := 0; i < allTimes; i++ {
			if queryData[i].One == 1 {
				ariseTimes += 1
			}
		}
		fmt.Printf("数字%d", specificNum)
		fmt.Printf("出现次数: %d, 占比: %.4f, 未出现次数: %d, 占比: %.4f, 总次数: %d, 与理论值偏差: %.4f%%\n", ariseTimes, float64(ariseTimes)/float64(allTimes), allTimes-ariseTimes, float64(allTimes-ariseTimes)/float64(allTimes), allTimes, (float64(ariseTimes)/float64(allTimes)-5.0/11.0)*100/(5.0/11.0))
	case 2:
		for i := 0; i < allTimes; i++ {
			if queryData[i].Two == 1 {
				ariseTimes += 1
			}
		}
		fmt.Printf("数字%d", specificNum)
		fmt.Printf("出现次数: %d, 占比: %.4f, 未出现次数: %d, 占比: %.4f, 总次数: %d, 与理论值偏差: %.4f%%\n", ariseTimes, float64(ariseTimes)/float64(allTimes), allTimes-ariseTimes, float64(allTimes-ariseTimes)/float64(allTimes), allTimes, (float64(ariseTimes)/float64(allTimes)-5.0/11.0)*100/(5.0/11.0))
	case 3:
		for i := 0; i < allTimes; i++ {
			if queryData[i].Three == 1 {
				ariseTimes += 1
			}
		}
		fmt.Printf("数字%d", specificNum)
		fmt.Printf("出现次数: %d, 占比: %.4f, 未出现次数: %d, 占比: %.4f, 总次数: %d, 与理论值偏差: %.4f%%\n", ariseTimes, float64(ariseTimes)/float64(allTimes), allTimes-ariseTimes, float64(allTimes-ariseTimes)/float64(allTimes), allTimes, (float64(ariseTimes)/float64(allTimes)-5.0/11.0)*100/(5.0/11.0))

	case 4:
		for i := 0; i < allTimes; i++ {
			if queryData[i].Four == 1 {
				ariseTimes += 1
			}
		}
		fmt.Printf("数字%d", specificNum)
		fmt.Printf("出现次数: %d, 占比: %.4f, 未出现次数: %d, 占比: %.4f, 总次数: %d, 与理论值偏差: %.4f%%\n", ariseTimes, float64(ariseTimes)/float64(allTimes), allTimes-ariseTimes, float64(allTimes-ariseTimes)/float64(allTimes), allTimes, (float64(ariseTimes)/float64(allTimes)-5.0/11.0)*100/(5.0/11.0))

	case 5:
		for i := 0; i < allTimes; i++ {
			if queryData[i].Five == 1 {
				ariseTimes += 1
			}
		}
		fmt.Printf("数字%d", specificNum)
		fmt.Printf("出现次数: %d, 占比: %.4f, 未出现次数: %d, 占比: %.4f, 总次数: %d, 与理论值偏差: %.4f%%\n", ariseTimes, float64(ariseTimes)/float64(allTimes), allTimes-ariseTimes, float64(allTimes-ariseTimes)/float64(allTimes), allTimes, (float64(ariseTimes)/float64(allTimes)-5.0/11.0)*100/(5.0/11.0))

	case 6:
		for i := 0; i < allTimes; i++ {
			if queryData[i].Six == 1 {
				ariseTimes += 1
			}
		}
		fmt.Printf("数字%d", specificNum)
		fmt.Printf("出现次数: %d, 占比: %.4f, 未出现次数: %d, 占比: %.4f, 总次数: %d, 与理论值偏差: %.4f%%\n", ariseTimes, float64(ariseTimes)/float64(allTimes), allTimes-ariseTimes, float64(allTimes-ariseTimes)/float64(allTimes), allTimes, (float64(ariseTimes)/float64(allTimes)-5.0/11.0)*100/(5.0/11.0))

	case 7:
		for i := 0; i < allTimes; i++ {
			if queryData[i].Seven == 1 {
				ariseTimes += 1
			}
		}
		fmt.Printf("数字%d", specificNum)
		fmt.Printf("出现次数: %d, 占比: %.4f, 未出现次数: %d, 占比: %.4f, 总次数: %d, 与理论值偏差: %.4f%%\n", ariseTimes, float64(ariseTimes)/float64(allTimes), allTimes-ariseTimes, float64(allTimes-ariseTimes)/float64(allTimes), allTimes, (float64(ariseTimes)/float64(allTimes)-5.0/11.0)*100/(5.0/11.0))

	case 8:
		for i := 0; i < allTimes; i++ {
			if queryData[i].Eight == 1 {
				ariseTimes += 1
			}
		}
		fmt.Printf("数字%d", specificNum)
		fmt.Printf("出现次数: %d, 占比: %.4f, 未出现次数: %d, 占比: %.4f, 总次数: %d, 与理论值偏差: %.4f%%\n", ariseTimes, float64(ariseTimes)/float64(allTimes), allTimes-ariseTimes, float64(allTimes-ariseTimes)/float64(allTimes), allTimes, (float64(ariseTimes)/float64(allTimes)-5.0/11.0)*100/(5.0/11.0))

	case 9:
		for i := 0; i < allTimes; i++ {
			if queryData[i].Nine == 1 {
				ariseTimes += 1
			}
		}
		fmt.Printf("数字%d", specificNum)
		fmt.Printf("出现次数: %d, 占比: %.4f, 未出现次数: %d, 占比: %.4f, 总次数: %d, 与理论值偏差: %.4f%%\n", ariseTimes, float64(ariseTimes)/float64(allTimes), allTimes-ariseTimes, float64(allTimes-ariseTimes)/float64(allTimes), allTimes, (float64(ariseTimes)/float64(allTimes)-5.0/11.0)*100/(5.0/11.0))

	case 10:
		for i := 0; i < allTimes; i++ {
			if queryData[i].Ten == 1 {
				ariseTimes += 1
			}
		}
		fmt.Printf("数字%d", specificNum)
		fmt.Printf("出现次数: %d, 占比: %.4f, 未出现次数: %d, 占比: %.4f, 总次数: %d, 与理论值偏差: %.4f%%\n", ariseTimes, float64(ariseTimes)/float64(allTimes), allTimes-ariseTimes, float64(allTimes-ariseTimes)/float64(allTimes), allTimes, (float64(ariseTimes)/float64(allTimes)-5.0/11.0)*100/(5.0/11.0))

	case 11:
		for i := 0; i < allTimes; i++ {
			if queryData[i].Eleven == 1 {
				ariseTimes += 1
			}
		}
		fmt.Printf("数字%d", specificNum)
		fmt.Printf("出现次数: %d, 占比: %.4f, 未出现次数: %d, 占比: %.4f, 总次数: %d, 与理论值偏差: %.4f%%\n", ariseTimes, float64(ariseTimes)/float64(allTimes), allTimes-ariseTimes, float64(allTimes-ariseTimes)/float64(allTimes), allTimes, (float64(ariseTimes)/float64(allTimes)-5.0/11.0)*100/(5.0/11.0))

	default:
		log.Fatalln("未知错误!")
	}
}

// 计算遗漏值出现的频次
func calLeaveAndTimes(m map[int][]string, specificNum int) {
	// 先把key排序
	keyList := make([]int, 0)
	for k, _ := range m {
		keyList = append(keyList, k)
	}
	sort.Ints(keyList)
	//for _, v := range keyList {
	//	if v>=9 {
	//		fmt.Printf("数字%d连续遗漏%d次的数量为:%d\n", specificNum, v, len(m[v]))
	//	}
	//}
	result := make([]string, 0)
	// keyList 最大值, 即最大遗漏值
	for _, v := range m[keyList[len(keyList)-1]] {
		result = append(result, orderNumConvertToDate(v))
	}
	fmt.Println("数字", specificNum, "遗漏的最大数字是: ", keyList[len(keyList)-1], "遗漏的期数是: ", m[keyList[len(keyList)-1]], "对应的时间是: ", result)
}

// 计算0出现的最大次数以及对应的orderNum
func calSpecificNumTimes(queryList []mysql.QueryData, specificNum int) map[int][]string {
	count := 0
	countMap := make(map[int][]string, 0)
	// 计算数字1
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
			orderNum = datePrefix + orderNum[len(orderNum)-3:]
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
		if len(k) != 11 || len(m[k]) != 12 {
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

// 从http链接解析数据
func saveData2Mysql2(prefix string) {
	dateRange := constructDate()
	dateRangeLen := len(dateRange)
	for i := 0; i < dateRangeLen; i++ {
		log.Println("当前登记的日期是: ", dateRange[i])
		saveData2MysqlByDate("gd/", dateRange[i])
		rand.Seed(time.Now().UnixNano())
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))
	}
}

// prefix: "jx/" || "gd/"
func saveData2Mysql(prefix string) {
	dateList := constructDate()
	dateRangeLength := len(dateList) - 1
	m := make(map[string][]int, 0)
	for i := 0; i < dateRangeLength; i++ {
		fileReader, err := os.Open(prefix + strings.ReplaceAll(dateList[i], "-", "") + ".html")
		if err != nil {
			log.Fatalln(err)
		}
		doc, err := goquery.NewDocumentFromReader(fileReader)
		if err != nil {
			log.Fatalln(err)
		}
		datePrefix := strings.ReplaceAll(dateList[i], "-", "")
		content := doc.Find(".middle").Find(".bo_ocs,.bt_ocs_orange,.bg_white,.mt15").Find(".ov,.bx_bb").Find(".w_333,.ov,.bx_bb").Find(".bg_white,.bg_grayf2")
		content.Each(func(i int, selection *goquery.Selection) {
			// 在这里分解num和value
			orderNum := selection.Find(".w_15").Text()
			if len(orderNum) == 6 {
				orderNum = datePrefix + orderNum[len(orderNum)-3:]
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
		if len(k) != 11 || len(m[k]) != 12 {
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

// 返回两个部分 [01 08 04 02 07] [0 1 1 0 1 0 0 1 1 0 0 0]
func orderValue(selection *goquery.Selection) []int {
	valueList := make([]string, 0)
	selection.Find(".w_63").Each(func(i int, selection *goquery.Selection) {
		selection.Find(".ds_ib").Each(func(i int, selection *goquery.Selection) {
			valueList = append(valueList, selection.Text())
		})
	})
	if len(valueList) == 0 {
	}
	binList := make([]int, 12)
	for _, v := range valueList {
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalln(err)
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
