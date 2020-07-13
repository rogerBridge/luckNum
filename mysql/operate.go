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