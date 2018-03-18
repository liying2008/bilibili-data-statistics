package main

import (
	"net/http"
	"io/ioutil"
	"strconv"
	"fmt"
	"./tool/error"
	"./tool/db/sqlite3"
	"./tool/db/mysql"
	"./tool/file"
	"./data"
	"flag"
	"time"
	"os"
	"strings"
	. "./tool/db/config"
)

const (
	URL_PREFIX     = "http://api.bilibili.com/archive_stat/stat?aid="
	CONTENT_TYPE   = "application/json; charset=utf-8"
	REFERER        = "https://www.bilibili.com/video/av11809669/?"
	USER_AGENT     = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36"
	CONNECTION     = "keep-alive"
	DB_CONFIG_FILE = "db_config.json"
)

// 开始 aid
var start uint64
// 结束 aid
var end uint64
// 写入同一数据库
var samedb bool

var client *http.Client

func main() {
	// 解析命令行参数
	flag.Uint64Var(&start, "start", 1, "the start aid (Include)")
	flag.Uint64Var(&end, "end", 100, "the end aid (Exclude)")
	flag.BoolVar(&samedb, "samedb", false, "write new data to an old database.")
	flag.Parse()

	if start > end {
		fmt.Println("start aid is greater than end aid, so start and end exchange.")
		start, end = end, start
	}

	if !file.Exists(DB_CONFIG_FILE) {
		fmt.Println(DB_CONFIG_FILE + " file does not exist. use sqlite3 database.")
		// 使用 SQLite3 数据库
		dataToSqlite3()
		return
	}
	configData, err := ioutil.ReadFile(DB_CONFIG_FILE)
	error.CheckErr(err)
	dbConfig := ParseDBConfig(configData)
	fmt.Println("use_mysql=" + strconv.FormatBool(dbConfig.UseMysql))

	if dbConfig.UseMysql {
		// 使用 MySQL 数据库
		dataToMysql(dbConfig)
	} else {
		// 使用 SQLite3 数据库
		dataToSqlite3()
	}
}

func dataToSqlite3() {
	// 数据库
	if file.Exists(sqlite3.DB_NAME) {
		// 旧数据库文件存在
		if !samedb {
			currTime := time.Now().Format("20060102150405")
			err := os.Rename(sqlite3.DB_NAME, strings.TrimRight(sqlite3.DB_NAME, ".db")+"-"+currTime+".db")
			error.CheckErr(err)
			sqlite3.InitDB()
		}
	} else {
		sqlite3.InitDB()
	}

	for i := start; i < end; i++ {
		jsonStr := getVideoData(i)
		if i%1000 == 0 {
			time.Sleep(10 * time.Second)
		} else if i%200 == 0 {
			time.Sleep(2 * time.Second)
		}
		if jsonStr == "" {
			fmt.Errorf("%s", "av"+strconv.FormatUint(i, 10)+": result is nil")
			continue
		}
		video := data.ParseVideoData(jsonStr)
		if video.Code == 0 {
			// 获取信息成功
			sqlite3.InsertData(video.Data)
			fmt.Println("av" + strconv.FormatUint(video.Aid, 10) + " Success!")
		} else {
			fmt.Println("failed to fetch av" + strconv.FormatUint(i, 10) + " data!")
			//fmt.Println(jsonStr)
		}
	}
}

func dataToMysql(config *Config) {
	mysql.InitDB(config)
	if !samedb {
		currTime := time.Now().Format("20060102150405")
		newTableName := mysql.TB_VIDEO_DATA + "-" + currTime
		mysql.RenameTable(config, newTableName)
		mysql.InitDB(config)
	}

	for i := start; i < end; i++ {
		jsonStr := getVideoData(i)
		if i%1000 == 0 {
			time.Sleep(10 * time.Second)
		} else if i%200 == 0 {
			time.Sleep(2 * time.Second)
		}
		if jsonStr == "" {
			fmt.Errorf("%s", "av"+strconv.FormatUint(i, 10)+": result is nil")
			continue
		}
		video := data.ParseVideoData(jsonStr)
		if video.Code == 0 {
			// 获取信息成功
			mysql.InsertData(config, video.Data)
			fmt.Println("av" + strconv.FormatUint(video.Aid, 10) + " Success!")
		} else {
			fmt.Println("failed to fetch av" + strconv.FormatUint(i, 10) + " data!")
			//fmt.Println(jsonStr)
		}
	}
}

func getVideoData(aid uint64) (data string) {
	if client == nil {
		client = &http.Client{}
	}
	url := URL_PREFIX + strconv.FormatUint(aid, 10)
	req, err := http.NewRequest("GET", url, nil)
	error.CheckErr(err)

	//req.Header.Set("Content-Type", CONTENT_TYPE)
	//req.Header.Set("Referer", REFERER)
	//req.Header.Set("User-Agent", USER_AGENT)
	//req.Header.Set("Connection", CONNECTION)

	resp, err := client.Do(req)

	if err != nil || resp == nil {
		fmt.Errorf("%s", err)
		return ""
	}
	body := resp.Body
	defer body.Close()
	if body != nil {
		body, err := ioutil.ReadAll(resp.Body)
		error.CheckErr(err)
		data = string(body)
		return
	} else {
		fmt.Errorf("%s", "body is nil")
		return ""
	}
}
