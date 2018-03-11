# bilibili 视频数据采集

## 编译执行

采集的数据存储到 **SQLite3** 数据库，需要依赖 `github.com/mattn/go-sqlite3` 。  
1. 下载 **go-sqlite3** 依赖库：
```shell
go get github.com/mattn/go-sqlite3
```
2. 编译（在当前目录下执行命令）：  
```shell
go build
```
3. 运行：
- Windows:
```shell
bilibili-data-statistics.exe --start={start aid} --end={end aid} --samedb={true/false}
```

- Linux:
```shell
./bilibili-data-statistics --start={start aid} --end={end aid} --samedb={true/false}
```
> 参数说明：  
start: 起始AV号，从此AV号开始采集；  
end: 结束AV号，采集到该AV号就停止采集（不采集该AV号）；  
samedb: 本次采集的数据是否和上次的数据存在一起（默认值：`false`）  
（false: 将上次采集生成的数据库文件重命名为以当前日期为结尾的文件名，并且新生成一个名为 `video_data.db` 的数据库文件，本次采集的数据都存在这个文件中；  
true: 不重命名上次采集生成的数据库文件，本次采集的数据依然放在同一文件中，便于多次分段采集。）；  

> 命令示例：  
`bilibili-data-statistics.exe --start=30000 --end=100000 --samedb=true`

目前使用单线程请求接口采集数据，每请求 200 次，线程休息 2s，每请求 1000 次，休息 10s（主要是为了不给bilibili的服务器造成压力）。  
采集速率：约 **5.8条/s**

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
