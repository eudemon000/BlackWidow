package Database

import (
	_"github.com/go-sql-driver/mysql"
	"database/sql"
	"BlackWidow/src/Config"
	logUtils "BlackWidow/src/logPackage"
	//"fmt"
	//"fmt"
)



var Db *sql.DB

//mysql配置
type mysqlConfig struct{

	host		string
	port		string
	database	string
	user		string
	password	string
	charset		string
}

type Urls struct {
	Url		string
	Md5		string
	Is_crawl	string
	Layer		int
	Content		string
}

//已爬取URL结构
type FinishedUrls struct {
	Id	int
	Url	string
	Md5	string
}

//待爬取URL结构
type WaitUrls struct {
	Id	int
	Url	string
}

var mysqlConf *mysqlConfig

func init() {
	mysqlConf = new(mysqlConfig)
	config := new(Config.Config)
	config.InitConfig()
	mysqlConf.host = config.Read("MySQL", "host")
	mysqlConf.port = config.Read("MySQL", "port")
	mysqlConf.database = config.Read("MySQL", "database")
	mysqlConf.user = config.Read("MySQL", "user")
	mysqlConf.password = config.Read("MySQL", "password")
	mysqlConf.charset = config.Read("MySQL", "charset")
	var err error
	connStr := mysqlConf.user + ":" + mysqlConf.password + "@tcp(" + mysqlConf.host + ":" + mysqlConf.port + ")/" + mysqlConf.database + "?charset=" + mysqlConf.charset
	Db, err = sql.Open("mysql", connStr)
	if err != nil {
		logUtils.Msg(logUtils.Error, err)
		panic(err)
	}
	Db.SetMaxOpenConns(20)
}

//检查是否存在Url
func GetUrl(url string) bool {
	rows, err := Db.Query("SELECT * FROM tbl_urls as u where u.url like ?", url)
	defer rows.Close()
	if err != nil {
		logUtils.Msg(logUtils.Error, err)
	}
	if rows.Next() {
		return true
	} else {
		return false
	}

	/*columns, err := rows.Columns()
	if err != nil {
		logUtils.Msg(logUtils.Error, err)
	}*/

}

//检查待爬取列表是否存在该url
func GetWaitUrl(url string) bool {
	rows, err := Db.Query("SELECT * FROM tbl_continue_url as u where u.url like ?", url)
	defer rows.Close()
	if err != nil {
		logUtils.Msg(logUtils.Error, err)
	}
	if rows.Next() {
		return true
	} else {
		return false
	}

	/*columns, err := rows.Columns()
	if err != nil {
		logUtils.Msg(logUtils.Error, err)
	}*/

}

//从带爬取表删除
func RemoveWaitUrl(url string) int64 {
	result, err := Db.Exec("delete from tbl_continue_url where url like ?", url)
	if err != nil {
		logUtils.Msg(logUtils.Error, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logUtils.Msg(logUtils.Error, err)
	}
	return rowsAffected
}

//向已爬取列表插入数据
func InsertUrls(url, md5, is_crawl, layer string) int64 {
	result, err := Db.Exec("insert into tbl_urls(url, md5, is_crawl, layer) values(?, ?, ?, ?)", url, md5, is_crawl, layer)
	if err != nil {
		logUtils.Msg(logUtils.Error, err)
		return -1
	}

	id, err := result.LastInsertId()
	if err != nil {
		logUtils.Msg(logUtils.Error, err)
		return -1
	}

	return id
}

//向待爬取列表插入url
func InsertWaitUrl(url string) int64 {
	result, err := Db.Exec("insert into tbl_continue_url(url) values(?)", url)
	if err != nil {
		logUtils.Msg(logUtils.Error, err)
		return -1
	}

	id, err := result.LastInsertId()
	if err != nil {
		logUtils.Msg(logUtils.Error, err)
		return -1
	}
	return id

}

//获取待爬的url列表
func GetFinishedUrls() []FinishedUrls {
	var w []FinishedUrls = make([]FinishedUrls, 0)
	rows, err := Db.Query("select * from tbl_urls")
	if err != nil {
		logUtils.Msg(logUtils.Error, err)
		return w
	}

	for rows.Next() {
		var id	int
		var url	string
		var md5	string
		err := rows.Scan(&id, &url, &md5)
		if err != nil {
			logUtils.Msg(logUtils.Error, err)
			return w
		}
		record := FinishedUrls{}

		record.Id = id
		record.Url = url
		record.Md5 = md5
		w = append(w, record)
	}
	return w
}

//获取待爬的url列表
func GetWaitUrls() []WaitUrls {
	var w []WaitUrls = make([]WaitUrls, 0)
	rows, err := Db.Query("select * from tbl_continue_url")
	if err != nil {
		logUtils.Msg(logUtils.Error, err)
		return w
	}

	for rows.Next() {
		var id		int
		var url		string
		var create_time string
		var change_time	string
		err := rows.Scan(&id, &url, &create_time, &change_time)
		if err != nil {
			logUtils.Msg(logUtils.Error, err)
			return w
		}
		record := WaitUrls{}

		record.Id = id
		record.Url = url
		w = append(w, record)
	}
	return w
}







