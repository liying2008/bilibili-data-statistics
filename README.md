# bilibili 视频数据采集

## 编译执行

采集的数据可以存储到 **SQLite3** 数据库 或 **MySQL** 数据库，取决于 `db_config.json` 文件中的配置： 
 
`db_config.json`   
```json
{
  "use_mysql": true,
  "database": "video_data",
  "username": "root",
  "password": "root"
}
```
`use_mysql` 置为 `true` 则使用 **MySQL** 数据库，置为 `false` 则使用 **SQLite3** 数据库。其余 `database` 、 `username` 和 `password` 三项配置只针对 **MySQL** 数据库，如果使用 **MySQL** 数据库，请填写正确的配置。 

### 1. 下载依赖库  

使用 **SQLite3** 数据库 需要依赖 `github.com/mattn/go-sqlite3` ：
- 下载方法：
```shell
go get github.com/mattn/go-sqlite3
```

使用 **MySQL** 数据库 需要依赖 `github.com/go-sql-driver/mysql` ：
- 下载方法：
```shell
go get github.com/go-sql-driver/mysql
```

### 2. 编译（在当前目录下执行命令）：  

```shell
go build
```

### 3. 运行：

- Windows:
```shell
bilibili-data-statistics.exe --start={start aid} --end={end aid} --samedb={true/false}
```

- Linux:
```shell
./bilibili-data-statistics --start={start aid} --end={end aid} --samedb={true/false}
```

> **参数说明：**  
**start**: 起始AV号，从此AV号开始采集；  
**end**: 结束AV号，采集到该AV号就停止采集（不采集该AV号）；  
**samedb**: 本次采集的数据是否和上次的数据存在一起（默认值：`false`）  
（false: 如果使用 SQLite3 ，则会将上次采集生成的数据库文件重命名为以当前日期为结尾的文件名，并且新生成一个名为 `video_data.db` 的数据库文件，本次采集的数据都存在这个文件中；如果使用 MySQL ，则会将上次生成的表重命名,名称以当前日期结尾，然后生成一个新表，本次采集的数据都存在于新表中。  
true: 如果使用 SQLite3 ，不重命名上次采集生成的数据库文件，本次采集的数据依然放在同一文件中；如果使用 MySQL ，不重名上次采集生成的表，本次采集的数据依然放在同一个表中。）；  

> 命令示例：  
`bilibili-data-statistics.exe --start=30000 --end=100000 --samedb=true`

目前使用单线程请求接口采集数据。  
采集速率：约 **6.2条/s**

请求的接口：http://api.bilibili.com/archive_stat/stat?aid={aid}

## 字段说明

```go
type Data struct {
	Aid       uint64 `视频AV号`
	View      int   `播放数`
	Danmaku   int   `弹幕数`
	Reply     int   `评论数`
	Favorite  int   `收藏数`
	Coin      int   `硬币数`
	Share     int   `分享数`
	NowRank   int   `当前排名`
	HisRank   int   `历史最高排名`
	Like      int   `喜欢数`
	NoReprint int   `暂不清楚含义`
	Copyright int   `稿件授权方式`
}
```
