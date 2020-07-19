package mysql

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func InitMysqlConn() *sql.DB {
	db, err := sql.Open("mysql", "leo:123456@tcp(127.0.0.1:3306)/db_play")
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
var conn = InitMysqlConn()

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

func Write2Luck(prefix string, specificNum int, leaveNum int, stopProbability float64, hopeIncome float64) error{
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
		if err!=nil {
			return err
		}
		return nil
	}
	if prefix == "gd/" {
		_, err := conn.Exec("DELETE FROM gd_luck")
		if err!=nil {
			return err
		}
		return nil
	}
	return errors.New("unsupported table types")
}

type LuckNum struct {
	SpecificNum int
	LeaveValue int
	StopProbability float64
	HopeIncome float64
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
			if err!=nil {
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
			if err!=nil {
				log.Println(err)
				return err
			}
		}
	}
	return nil
}

// 验证之前的预测是否为真
func DetectForecast(prefix string) error {
	type DetectForecast struct {
		OrderNum string
		ForecastNum int
	}
	contrast := map[int]string{
		1: "one",
		2: "two",
		3: "three",
		4: "four",
		5: "five",
		6: "six",
		7: "seven",
		8: "eight",
		9: "nine",
		10: "ten",
		11: "eleven",
	}
	var detectList []DetectForecast

	if prefix == "gd/" {
		// 下一个优化的点: 只验证当天的order_num和forecast_num
		rows, err := conn.Query("SELECT order_num, forecast_num FROM db_play.forecast_gd")
		if err!=nil {
			log.Println(err)
			return err
		}
		defer rows.Close()
		for rows.Next() {
			d := DetectForecast{}
			err := rows.Scan(&d.OrderNum, &d.ForecastNum)
			if err!=nil {
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
			if err!=nil {
				log.Println("获取真实结果时: ", err)
				return err
			}
			_, err = conn.Exec("UPDATE forecast_gd SET forecast_result = ? WHERE order_num = ? AND forecast_num = ?", isTrue, v.OrderNum, v.ForecastNum)
			if err!=nil {
				log.Println(err)
				return err
			}
		}
		return nil
	}

	if prefix == "jx/" {
		rows, err := conn.Query("SELECT order_num, forecast_num FROM db_play.forecast_jx")
		if err!=nil {
			log.Println(err)
			return err
		}
		defer rows.Close()
		for rows.Next() {
			d := DetectForecast{}
			err := rows.Scan(&d.OrderNum, &d.ForecastNum)
			if err!=nil {
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
			if err!=nil {
				log.Println("获取真实结果时: ", err)
				return err
			}
			_, err = conn.Exec("UPDATE forecast_jx SET forecast_result = ? WHERE order_num = ? AND forecast_num = ?", isTrue, v.OrderNum, v.ForecastNum)
			if err!=nil {
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
		OrderNum string
		ForecastNum int
	}
	contrast := map[int]string{
		1: "one",
		2: "two",
		3: "three",
		4: "four",
		5: "five",
		6: "six",
		7: "seven",
		8: "eight",
		9: "nine",
		10: "ten",
		11: "eleven",
	}
	var detectList []DetectForecast

	if prefix == "gd/" {
		// 下一个优化的点: 只验证当天的order_num和forecast_num
		rows, err := conn.Query("SELECT order_num, forecast_num FROM db_play.forecast_gd")
		if err!=nil {
			log.Println(err)
			return err
		}
		defer rows.Close()
		for rows.Next() {
			d := DetectForecast{}
			err := rows.Scan(&d.OrderNum, &d.ForecastNum)
			if err!=nil {
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
			if err!=nil {
				log.Println("获取真实结果时: ", err)
				return err
			}
			_, err = conn.Exec("UPDATE forecast_gd SET forecast_result = ? WHERE order_num = ? AND forecast_num = ?", isTrue, v.OrderNum, v.ForecastNum)
			if err!=nil {
				log.Println(err)
				return err
			}
		}
		return nil
	}

	if prefix == "jx/" {
		rows, err := conn.Query("SELECT order_num, forecast_num FROM db_play.forecast_jx")
		if err!=nil {
			log.Println(err)
			return err
		}
		defer rows.Close()
		for rows.Next() {
			d := DetectForecast{}
			err := rows.Scan(&d.OrderNum, &d.ForecastNum)
			if err!=nil {
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
			if err!=nil {
				log.Println("获取真实结果时: ", err)
				return err
			}
			_, err = conn.Exec("UPDATE forecast_jx SET forecast_result = ? WHERE order_num = ? AND forecast_num = ?", isTrue, v.OrderNum, v.ForecastNum)
			if err!=nil {
				log.Println(err)
				return err
			}
		}
		return nil
	}
	return errors.New("unsupported prefix")
}