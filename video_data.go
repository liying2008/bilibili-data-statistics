package main

import (
	"flag"
	"fmt"
	"github.com/bilibili-data-statistics/data"
	. "github.com/bilibili-data-statistics/tool/db/config"
	"github.com/bilibili-data-statistics/tool/db/mysql"
	"github.com/bilibili-data-statistics/tool/db/sqlite3"
	"github.com/bilibili-data-statistics/tool/error"
	"github.com/bilibili-data-statistics/tool/file"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	URL_PREFIX     = "http://api.bilibili.com/archive_stat/stat?aid="
	CONNECTION     = "keep-alive"
	DB_CONFIG_FILE = "db_config.json"
	GROUP_COUNT    = 200
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

	var id int64
	oldLastCount := start
	lastCount := start
	times := getGroupCount()
	fmt.Println("times = " + strconv.FormatUint(times, 10))
	for i := uint64(0); i < times; i++ {
		oldLastCount = lastCount
		if lastCount+GROUP_COUNT < end {
			lastCount = lastCount + GROUP_COUNT
		} else {
			lastCount = end
		}
		groupData := getGroupData(oldLastCount, lastCount)
		// 插入数据库
		id = sqlite3.InsertGroupData(groupData)
		fmt.Println("last insert id: " + strconv.FormatInt(id, 10))
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

	var id int64
	oldLastCount := start
	lastCount := start
	times := getGroupCount()
	fmt.Println("times = " + strconv.FormatUint(times, 10))
	for i := uint64(0); i < times; i++ {
		oldLastCount = lastCount
		if lastCount+GROUP_COUNT < end {
			lastCount = lastCount + GROUP_COUNT
		} else {
			lastCount = end
		}
		groupData := getGroupData(oldLastCount, lastCount)
		// 插入数据库
		id = mysql.InsertGroupData(config, groupData)
		fmt.Println("last insert id: " + strconv.FormatInt(id, 10))
	}
}

// 得到分组个数，决定了集中写数据库的次数
func getGroupCount() (times uint64) {
	if (end-start+1)%GROUP_COUNT == 0 {
		times = (end - start + 1) / GROUP_COUNT
	} else {
		times = (end-start+1)/GROUP_COUNT + 1
	}
	return times
}

func getGroupData(s uint64, e uint64) []*data.Data {
	groupData := make([]*data.Data, e-s+1)
	for i := s; i < e; i++ {
		jsonStr := getVideoData(i)
		if jsonStr == "" {
			fmt.Errorf("%s", "av"+strconv.FormatUint(i, 10)+": result is nil")
			continue
		}
		video := data.ParseVideoData(jsonStr)
		if video.Code == 0 {
			// 获取信息成功
			groupData[i-s] = video.Data
			fmt.Println("av" + strconv.FormatUint(video.Aid, 10) + " Success!")
		} else {
			fmt.Println("failed to fetch av" + strconv.FormatUint(i, 10) + " data!")
			//fmt.Println(jsonStr)
		}
	}
	return groupData
}

func getVideoData(aid uint64) (data string) {
	if client == nil {
		client = &http.Client{}
	}
	url := URL_PREFIX + strconv.FormatUint(aid, 10)
	req, err := http.NewRequest("GET", url, nil)
	error.CheckErr(err)
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
