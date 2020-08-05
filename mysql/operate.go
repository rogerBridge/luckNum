package mysql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"sort"
	"strconv"
)

func ReadConfig(file string) string {
	f, err := os.Open(file)
	if err!=nil {
		log.Fatalln(err)
	}
	type config struct {
		Username string `json:"username"`
		Password string `json:"password"`
		IpAddress string `json:"ipAddress"`
		Port string `json:"port"`
		Database string `json:"database"`
	}
	var c config
	err = json.NewDecoder(f).Decode(&c)
	if err!=nil {
		log.Fatalln(err)
	}
	//fmt.Println(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.Username, c.Password, c.IpAddress, c.Port, c.Database))
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.Username, c.Password, c.IpAddress, c.Port, c.Database)
}

func InitMysqlConn() *sql.DB {
	db, err := sql.Open("mysql", ReadConfig("mysqlConfig.json"))
	if err != nil {
		log.Fatalf("conn establish error\n")
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("conn establish error\n")
	}
	return db
}

// conn 设置为全局变量, 方便复用, 初始化失败时, 服务不能启动
//var conn = InitMysqlConn()
var conn = new(sql.DB)

func SaveResultToMysqlJx(orderNum string, orderTime string, valueList []int, valueString string) error {
	if len(valueList) != 12 {
		log.Fatalln("源头数据错误, 不能写入数据库!")
	}
	_, err := conn.Exec("INSERT INTO db_play.jx11x5 (order_number, order_time, one, two, three, four, five, six, seven, eight, nine, ten, eleven, result) VALUES (?, ?, ?, ?, ?,?, ?, ?, ?,?, ?, ?, ?, ?)", orderNum, orderTime, valueList[1], valueList[2], valueList[3], valueList[4], valueList[5], valueList[6], valueList[7], valueList[8], valueList[9], valueList[10], valueList[11], valueString)
	if err != nil {
		log.Printf("%+v insert error: %+v", conn, err)
		return err
	}
	return nil
}

func SaveResultToMysqlGd(orderNum string, orderTime string, valueList []int, valueString string) error {
	if len(valueList) != 12 {
		log.Fatalln("源头数据错误, 不能写入数据库!")
	}
	_, err := conn.Exec("INSERT INTO db_play.gd11x5 (order_number, order_time,  one, two, three, four, five, six, seven, eight, nine, ten, eleven, result) VALUES (?,?, ?, ?, ?,?, ?, ?, ?,?, ?, ?, ?, ?)", orderNum, orderTime, valueList[1], valueList[2], valueList[3], valueList[4], valueList[5], valueList[6], valueList[7], valueList[8], valueList[9], valueList[10], valueList[11], valueString)
	if err != nil {
		log.Printf("%+v insert error: %+v", conn, err)
		return err
	}
	return nil
}

type QueryData struct {
	OrderNum string
	One      int
	Two      int
	Three    int
	Four     int
	Five     int
	Six      int
	Seven    int
	Eight    int
	Nine     int
	Ten      int
	Eleven   int
}

func QueryDataFromMysqlJx() ([]QueryData, error) {
	queryDataList := make([]QueryData, 0)
	rows, err := conn.Query("SELECT order_number, one, two, three, four, five, six, seven, eight, nine, ten, eleven FROM jx11x5")
	if err != nil {
		return []QueryData{}, err
	}
	defer rows.Close()
	for rows.Next() {
		q := QueryData{}
		err := rows.Scan(&q.OrderNum, &q.One, &q.Two, &q.Three, &q.Four, &q.Five, &q.Six, &q.Seven, &q.Eight, &q.Nine, &q.Ten, &q.Eleven)
		if err != nil {
			return []QueryData{}, err
		}
		queryDataList = append(queryDataList, q)
	}
	return queryDataList, nil
}

func QueryDataFromMysqlGd() ([]QueryData, error) {
	queryDataList := make([]QueryData, 0)
	rows, err := conn.Query("SELECT order_number, one, two, three, four, five, six, seven, eight, nine, ten, eleven FROM gd11x5")
	if err != nil {
		return []QueryData{}, err
	}
	defer rows.Close()
	for rows.Next() {
		q := QueryData{}
		err := rows.Scan(&q.OrderNum, &q.One, &q.Two, &q.Three, &q.Four, &q.Five, &q.Six, &q.Seven, &q.Eight, &q.Nine, &q.Ten, &q.Eleven)
		if err != nil {
			return []QueryData{}, err
		}
		queryDataList = append(queryDataList, q)
	}
	return queryDataList, nil
}

// orderNum 是否存在于数据库
func IsExistInMysql(prefix string, orderNum string) (int, error) {
	if prefix == "gd/" {
		var ifExist int
		err := conn.QueryRow("SELECT IF (  EXISTS(SELECT order_number FROM db_play.gd11x5 WHERE order_number = ?), 1, 0)", orderNum).Scan(&ifExist)
		if err != nil {
			log.Println(err)
			return -1, err
		}
		if ifExist == 1 {
			log.Println(orderNum, "已存在")
			return 1, nil
		}
		if ifExist == 0 {
			log.Println(orderNum, "不存在")
			return 0, nil
		}
	}
	if prefix == "jx/" {
		var ifExist int
		err := conn.QueryRow("SELECT IF (  EXISTS(SELECT order_number FROM db_play.jx11x5 WHERE order_number = ?), 1, 0)", orderNum).Scan(&ifExist)
		if err != nil {
			log.Println(err)
			return -1, err
		}
		if ifExist == 1 {
			log.Println(orderNum, "已存在")
			return 1, nil
		}
		if ifExist == 0 {
			log.Println(orderNum, "不存在")
			return 0, nil
		}
	}
	return -1, errors.New("prefix not supported")
}

// 获取 "20200710" 之类的orderNum前缀
func QueryDataFromMysqlGdSomeDay(dateStr string) ([]QueryData, error) {
	queryDataList := make([]QueryData, 0)
	rows, err := conn.Query("SELECT order_number, one, two, three, four, five, six, seven, eight, nine, ten, eleven FROM gd11x5 WHERE order_number LIKE ?", dateStr+"__")
	if err != nil {
		return []QueryData{}, err
	}
	defer rows.Close()
	for rows.Next() {
		q := QueryData{}
		err := rows.Scan(&q.OrderNum, &q.One, &q.Two, &q.Three, &q.Four, &q.Five, &q.Six, &q.Seven, &q.Eight, &q.Nine, &q.Ten, &q.Eleven)
		if err != nil {
			return []QueryData{}, err
		}
		queryDataList = append(queryDataList, q)
	}
	return queryDataList, nil
}

// 获取 "20200710" 之类的orderNum前缀
func QueryDataFromMysqlJxSomeDay(dateStr string) ([]QueryData, error) {
	queryDataList := make([]QueryData, 0)
	rows, err := conn.Query("SELECT order_number, one, two, three, four, five, six, seven, eight, nine, ten, eleven FROM jx11x5 WHERE order_number LIKE ?", dateStr+"__")
	if err != nil {
		return []QueryData{}, err
	}
	defer rows.Close()
	for rows.Next() {
		q := QueryData{}
		err := rows.Scan(&q.OrderNum, &q.One, &q.Two, &q.Three, &q.Four, &q.Five, &q.Six, &q.Seven, &q.Eight, &q.Nine, &q.Ten, &q.Eleven)
		if err != nil {
			return []QueryData{}, err
		}
		queryDataList = append(queryDataList, q)
	}
	return queryDataList, nil
}

func Write2Luck(prefix string, specificNum int, leaveNum int, stopProbability float64, hopeIncome float64) error {
	if prefix == "jx/" {
		// Result.LastInsertId() 表示受影响的最近的columnsID
		_, err := conn.Exec("INSERT INTO db_play.jx_luck(specific_num, leave_value, stop_probability, hope_income) VALUES (?, ?, ?, ?)", specificNum, leaveNum, stopProbability, hopeIncome)
		if err != nil {
			log.Printf("%+v insert error: %+v", conn, err)
			return err
		}
		return nil
	}
	if prefix == "gd/" {
		_, err := conn.Exec("INSERT INTO db_play.gd_luck(specific_num, leave_value, stop_probability, hope_income) VALUES (?, ?, ?, ?)", specificNum, leaveNum, stopProbability, hopeIncome)
		if err != nil {
			log.Printf("%+v insert error: %+v", conn, err)
			return err
		}
		return nil
	}
	return errors.New("unsupported type")
}

func DeleteLuckTable(prefix string) error {
	if prefix == "jx/" {
		_, err := conn.Exec("DELETE FROM jx_luck")
		if err != nil {
			return err
		}
		return nil
	}
	if prefix == "gd/" {
		_, err := conn.Exec("DELETE FROM gd_luck")
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("unsupported table types")
}

type LuckNum struct {
	SpecificNum     int
	LeaveValue      int
	StopProbability float64
	HopeIncome      float64
}

func GetDataFromLuckTable(prefix string) ([]LuckNum, error) {
	if prefix == "jx/" {
		luckNumList := make([]LuckNum, 0)
		rows, err := conn.Query("SELECT specific_num, leave_value, stop_probability, hope_income FROM jx_luck")
		if err != nil {
			return []LuckNum{}, err
		}
		defer rows.Close()
		for rows.Next() {
			q := LuckNum{}
			err := rows.Scan(&q.SpecificNum, &q.LeaveValue, &q.StopProbability, &q.HopeIncome)
			if err != nil {
				return []LuckNum{}, err
			}
			luckNumList = append(luckNumList, q)
		}
		return luckNumList, nil
	}
	if prefix == "gd/" {
		luckNumList := make([]LuckNum, 0)
		rows, err := conn.Query("SELECT specific_num, leave_value, stop_probability, hope_income FROM gd_luck")
		if err != nil {
			return []LuckNum{}, err
		}
		defer rows.Close()
		for rows.Next() {
			q := LuckNum{}
			err := rows.Scan(&q.SpecificNum, &q.LeaveValue, &q.StopProbability, &q.HopeIncome)
			if err != nil {
				return []LuckNum{}, err
			}
			luckNumList = append(luckNumList, q)
		}
		return luckNumList, nil
	}
	return []LuckNum{}, errors.New("unsupport prefix")
}

// 将预测的数据存入forecast_{jx|gd} 表格
func StoreResultToForecastTable(prefix string, orderNum string, forecastNum int) error {
	if prefix == "gd/" {
		// 首先, 看一下预测orderNum是否存在于table中
		var ifExist int
		err := conn.QueryRow("SELECT IF (  EXISTS(SELECT order_num FROM db_play.forecast_gd WHERE order_num = ? AND forecast_num = ?), 1, 0)", orderNum, forecastNum).Scan(&ifExist)
		if err != nil {
			log.Println(err)
			return err
		}
		if ifExist == 1 {
			log.Printf("%s:%s 存在于forecast表格中, 不用再次添加\n", prefix, orderNum)
			//log.Println(orderNum, "已存在")
			return nil
		}
		if ifExist == 0 {
			log.Printf("%s:%s 不存在于forecast表格中, 可以添加\n", prefix, orderNum)
			//log.Println(orderNum, "不存在")
			_, err := conn.Exec("INSERT INTO db_play.forecast_gd (order_num, forecast_num) VALUES (?, ?)", orderNum, forecastNum)
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}

	if prefix == "jx/" {
		// 首先, 看一下预测orderNum是否存在于table中
		var ifExist int
		err := conn.QueryRow("SELECT IF (EXISTS(SELECT order_num FROM db_play.forecast_jx WHERE order_num = ? AND forecast_num = ?), 1, 0) ", orderNum, forecastNum).Scan(&ifExist)
		if err != nil {
			log.Println(err)
			return err
		}
		if ifExist == 1 {
			log.Printf("%s:%s 存在于forecast表格中, 不用再次添加\n", prefix, orderNum)
			//log.Println(orderNum, "已存在")
			return nil
		}
		if ifExist == 0 {
			//log.Println(orderNum, "不存在")
			log.Printf("%s:%s 不存在于forecast表格中, 可以添加\n", prefix, orderNum)
			_, err := conn.Exec("INSERT INTO db_play.forecast_jx (order_num, forecast_num) VALUES (?, ?)", orderNum, forecastNum)
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}
	return nil
}

// 验证之前的预测是否为真
func DetectForecast(prefix string) error {
	contrast := map[int]string{
		1:  "one",
		2:  "two",
		3:  "three",
		4:  "four",
		5:  "five",
		6:  "six",
		7:  "seven",
		8:  "eight",
		9:  "nine",
		10: "ten",
		11: "eleven",
	}
	type DetectForecast struct {
		OrderNum    string
		ForecastNum int
	}
	var detectList []DetectForecast

	if prefix == "gd/" {
		// 下一个优化的点: 只验证当天的order_num和forecast_num
		rows, err := conn.Query("SELECT order_num, forecast_num FROM db_play.forecast_gd")
		if err != nil {
			log.Println(err)
			return err
		}
		defer rows.Close()
		for rows.Next() {
			d := DetectForecast{}
			err := rows.Scan(&d.OrderNum, &d.ForecastNum)
			if err != nil {
				log.Println(err)
				return err
			}
			detectList = append(detectList, d)
		}
		// 拿到detectList之后, 把真实的值写入到forecast表里
		for _, v := range detectList {
			var isTrue int
			queryString := "SELECT " + contrast[v.ForecastNum] + " FROM gd11x5 WHERE order_number = ?"
			//fmt.Println(queryString)
			err := conn.QueryRow(queryString, v.OrderNum).Scan(&isTrue)
			if err != nil {
				log.Println("获取真实结果时: ", err)
				return err
			}
			_, err = conn.Exec("UPDATE forecast_gd SET forecast_result = ? WHERE order_num = ? AND forecast_num = ?", isTrue, v.OrderNum, v.ForecastNum)
			if err != nil {
				log.Println(err)
				return err
			}
		}
		return nil
	}

	if prefix == "jx/" {
		rows, err := conn.Query("SELECT order_num, forecast_num FROM db_play.forecast_jx")
		if err != nil {
			log.Println(err)
			return err
		}
		defer rows.Close()
		for rows.Next() {
			d := DetectForecast{}
			err := rows.Scan(&d.OrderNum, &d.ForecastNum)
			if err != nil {
				log.Println(err)
				return err
			}
			detectList = append(detectList, d)
		}
		// 拿到detectList之后, 把真实的值写入到forecast表里
		for _, v := range detectList {
			var isTrue int
			queryString := "SELECT " + contrast[v.ForecastNum] + " FROM jx11x5 WHERE order_number = ?"
			//fmt.Println(queryString)
			err := conn.QueryRow(queryString, v.OrderNum).Scan(&isTrue)
			if err != nil {
				log.Println("获取真实结果时: ", err)
				return err
			}
			_, err = conn.Exec("UPDATE forecast_jx SET forecast_result = ? WHERE order_num = ? AND forecast_num = ?", isTrue, v.OrderNum, v.ForecastNum)
			if err != nil {
				log.Println(err)
				return err
			}
		}
		return nil
	}
	return errors.New("unsupported prefix")
}

// 验证之前的预测是否为真
func DetectForecastImmediately(prefix string) error {
	conn := InitMysqlConn()
	defer conn.Close()

	type DetectForecast struct {
		OrderNum    string
		ForecastNum int
	}
	contrast := map[int]string{
		1:  "one",
		2:  "two",
		3:  "three",
		4:  "four",
		5:  "five",
		6:  "six",
		7:  "seven",
		8:  "eight",
		9:  "nine",
		10: "ten",
		11: "eleven",
	}
	var detectList []DetectForecast

	if prefix == "gd/" {
		// 下一个优化的点: 只验证当天的order_num和forecast_num
		rows, err := conn.Query("SELECT order_num, forecast_num FROM db_play.forecast_gd")
		if err != nil {
			log.Println(err)
			return err
		}
		defer rows.Close()
		for rows.Next() {
			d := DetectForecast{}
			err := rows.Scan(&d.OrderNum, &d.ForecastNum)
			if err != nil {
				log.Println(err)
				return err
			}
			detectList = append(detectList, d)
		}
		// 拿到detectList之后, 把真实的值写入到forecast表里
		for _, v := range detectList {
			var isTrue int
			queryString := "SELECT " + contrast[v.ForecastNum] + " FROM gd11x5 WHERE order_number = ?"
			//fmt.Println(queryString)
			err := conn.QueryRow(queryString, v.OrderNum).Scan(&isTrue)
			if err != nil {
				log.Println("获取真实结果时: ", err)
				return err
			}
			_, err = conn.Exec("UPDATE forecast_gd SET forecast_result = ? WHERE order_num = ? AND forecast_num = ?", isTrue, v.OrderNum, v.ForecastNum)
			if err != nil {
				log.Println(err)
				return err
			}
		}
		return nil
	}

	if prefix == "jx/" {
		rows, err := conn.Query("SELECT order_num, forecast_num FROM db_play.forecast_jx")
		if err != nil {
			log.Println(err)
			return err
		}
		defer rows.Close()
		for rows.Next() {
			d := DetectForecast{}
			err := rows.Scan(&d.OrderNum, &d.ForecastNum)
			if err != nil {
				log.Println(err)
				return err
			}
			detectList = append(detectList, d)
		}
		// 拿到detectList之后, 把真实的值写入到forecast表里
		for _, v := range detectList {
			var isTrue int
			queryString := "SELECT " + contrast[v.ForecastNum] + " FROM jx11x5 WHERE order_number = ?"
			//fmt.Println(queryString)
			err := conn.QueryRow(queryString, v.OrderNum).Scan(&isTrue)
			if err != nil {
				log.Println("获取真实结果时: ", err)
				return err
			}
			_, err = conn.Exec("UPDATE forecast_jx SET forecast_result = ? WHERE order_num = ? AND forecast_num = ?", isTrue, v.OrderNum, v.ForecastNum)
			if err != nil {
				log.Println(err)
				return err
			}
		}
		return nil
	}
	return errors.New("unsupported prefix")
}

type ForecastProve struct {
	OrderNum       string
	ForecastNum    int
	ForecastResult interface{}
}

// 统计整理statistics表格里面的各种数据
func StatisticsForecast(prefix string) (string, error) {
	var forecastList []ForecastProve
	if prefix == "gd/" {
		rows, err := conn.Query("SELECT order_num, forecast_num, forecast_result FROM forecast_gd")
		if err != nil {
			log.Printf("%v\n", err)
			return "", err
		}
		defer rows.Close()
		for rows.Next() {
			var f = ForecastProve{}
			err := rows.Scan(&f.OrderNum, &f.ForecastNum, &f.ForecastResult)
			if err != nil {
				log.Printf("%v\n", err)
				return "", err
			}
			// 只将forecastResult结果为int的数据存入forecastList
			fmt.Println()
			if _, ok := f.ForecastResult.([]uint8); ok {
				num, _ := strconv.Atoi(string(f.ForecastResult.([]uint8)[0]))
				forecastList = append(forecastList, ForecastProve{f.OrderNum, f.ForecastNum, num})
			}
		}
		//fmt.Println(prefix, forecastList)
		// 开始统计, 总猜测次数, 猜错次数, 占比, 猜对次数, 占比, 最大猜错次数
		return StatisticsForecastList(forecastList, prefix), nil
	}
	if prefix == "jx/" {
		//allLines := 0
		//err := conn.QueryRow("SELECT count(*) FROM forecast_gd").Scan(&allLines)
		//if err != nil {
		//	log.Println(err)
		//	return "", err
		//}
		rows, err := conn.Query("SELECT order_num, forecast_num, forecast_result FROM forecast_jx")
		if err != nil {
			log.Printf("%v\n", err)
			return "", err
		}
		defer rows.Close()
		for rows.Next() {
			var f = ForecastProve{}
			err := rows.Scan(&f.OrderNum, &f.ForecastNum, &f.ForecastResult)
			if err != nil {
				log.Printf("%v\n", err)
				return "", err
			}
			// 只将forecastResult结果为int的数据存入forecastList
			fmt.Println()
			if _, ok := f.ForecastResult.([]uint8); ok {
				num, _ := strconv.Atoi(string(f.ForecastResult.([]uint8)[0]))
				forecastList = append(forecastList, ForecastProve{f.OrderNum, f.ForecastNum, num})
			}
		}
		//fmt.Println(prefix, forecastList)
		// 开始统计, 总猜测次数, 猜错次数, 占比, 猜对次数, 占比, 最大猜错次数
		return StatisticsForecastList(forecastList, prefix), nil
	}
	return "", errors.New("不受支持的prefix类型")
}

// 对从 forecast 表格里面拿到的数据进行加工处理, 输出字符串
// 开始统计, 总猜测次数, 猜错次数, 占比, 猜对次数, 占比, 最大猜错次数
func StatisticsForecastList(data []ForecastProve, prefix string) string {
	allGuessNumber := len(data)
	allGuessWrong := 0
	allGuessTrue := 0
	allGuessContinousWrongList := make([]int, 0)
	for _, v := range data {
		if v.ForecastResult.(int) == 1 {
			allGuessTrue += 1
		}
		if v.ForecastResult.(int) == 0 {
			allGuessWrong += 1
		}
	}
	singleResult := 0
	for _, v := range data {
		if v.ForecastResult.(int) == 0 {
			singleResult += 1
		} else {
			allGuessContinousWrongList = append(allGuessContinousWrongList, singleResult)
			singleResult = 0
		}
	}
	sort.Ints(allGuessContinousWrongList)
	maxSingle := allGuessContinousWrongList[len(allGuessContinousWrongList)-1]
	return fmt.Sprintf("%s: 总猜测数量: %d, 猜对: %d次, 猜错: %d次, 最大连续猜错数量: %d次\n", prefix, allGuessNumber, allGuessTrue, allGuessWrong, maxSingle)
}
