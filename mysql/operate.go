package mysql

import (
	"database/sql"
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
	if len(valueList) != 12 || len(orderNum) != 11 {
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
	One int
	Two int
	Three int
	Four int
	Five int
	Six int
	Seven int
	Eight int
	Nine int
	Ten int
	Eleven int
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
		if err!=nil {
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
		if err!=nil {
			return []QueryData{}, err
		}
		queryDataList = append(queryDataList, q)
	}
	return queryDataList, nil
}