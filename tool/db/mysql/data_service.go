package mysql

import (
	"database/sql"
	. "../config"
	. "../../error"
	. "../../../data"
	_ "github.com/go-sql-driver/mysql"
)

// Initialize database
func InitDB(config *Config) {
	db, err := sql.Open(DRIVER_NAME, config.Username+":"+config.Password+"@/"+config.Database)
	CheckErr(err)
	// Create table.
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS " + TB_VIDEO_DATA + " (" +
		ID + " INTEGER PRIMARY KEY AUTO_INCREMENT," +
		AID + " INTEGER NOT NULL," +
		VIEW + " INTEGER DEFAULT -1," +
		DANMAKU + " INTEGER DEFAULT -1," +
		REPLY + " INTEGER DEFAULT -1," +
		FAVORITE + " INTEGER DEFAULT -1," +
		COIN + " INTEGER DEFAULT -1," +
		SHARE + " INTEGER DEFAULT -1," +
		NOW_RANK + " INTEGER DEFAULT -1," +
		HIS_RANK + " INTEGER DEFAULT -1," +
		LIKE + " INTEGER DEFAULT -1," +
		NO_REPRINT + " INTEGER DEFAULT -1," +
		COPYRIGHT + " INTEGER DEFAULT -1" +
		");")
	defer stmt.Close()
	CheckErr(err)
	_, err = stmt.Exec()
	CheckErr(err)
}

func RenameTable(config *Config, newTableName string) {
	db, err := sql.Open(DRIVER_NAME, config.Username+":"+config.Password+"@/"+config.Database)
	CheckErr(err)
	// rename table.
	stmt, err := db.Prepare("ALTER TABLE `" + TB_VIDEO_DATA + "` RENAME TO `" + newTableName + "`;")
	defer stmt.Close()
	stmt.Exec()
}

// Insert a data to mysql database
func InsertData(config *Config, data *Data) int64 {
	db, err := sql.Open(DRIVER_NAME, config.Username+":"+config.Password+"@/"+config.Database)
	CheckErr(err)
	stmt, err := db.Prepare("INSERT INTO " + TB_VIDEO_DATA + " (" +
		AID + ", " +
		VIEW + ", " +
		DANMAKU + ", " +
		REPLY + ", " +
		FAVORITE + ", " +
		COIN + ", " +
		SHARE + ", " +
		NOW_RANK + ", " +
		HIS_RANK + ", " +
		LIKE + ", " +
		NO_REPRINT + ", " +
		COPYRIGHT +
		") VALUES(?,?,?,?,?,?,?,?,?,?,?,?)")
	defer stmt.Close()
	CheckErr(err)

	res, err := stmt.Exec(data.Aid, data.View, data.Danmaku, data.Reply, data.Favorite,
		data.Coin, data.Share, data.NowRank, data.HisRank, data.Like, data.NoReprint, data.Copyright)
	CheckErr(err)

	id, err := res.LastInsertId()
	CheckErr(err)
	return id
}

// Insert data to mysql database
func InsertGroupData(config *Config, datas []*Data) int64 {
	db, err := sql.Open(DRIVER_NAME, config.Username+":"+config.Password+"@/"+config.Database)
	CheckErr(err)
	stmt, err := db.Prepare("INSERT INTO " + TB_VIDEO_DATA + " (" +
		AID + ", " +
		VIEW + ", " +
		DANMAKU + ", " +
		REPLY + ", " +
		FAVORITE + ", " +
		COIN + ", " +
		SHARE + ", " +
		NOW_RANK + ", " +
		HIS_RANK + ", " +
		LIKE + ", " +
		NO_REPRINT + ", " +
		COPYRIGHT +
		") VALUES(?,?,?,?,?,?,?,?,?,?,?,?)")
	defer stmt.Close()
	CheckErr(err)
	var res sql.Result
	for _, data := range datas {
		if data != nil {
			res, err = stmt.Exec(data.Aid, data.View, data.Danmaku, data.Reply, data.Favorite,
				data.Coin, data.Share, data.NowRank, data.HisRank, data.Like, data.NoReprint, data.Copyright)
			CheckErr(err)
		}
	}
	id, err := res.LastInsertId()
	CheckErr(err)
	return id
}

// Get all data
func GetAllData(config *Config, ) []Data {
	db, err := sql.Open(DRIVER_NAME, config.Username+":"+config.Password+"@/"+config.Database)
	// query data
	rows, err := db.Query("SELECT * FROM " + TB_VIDEO_DATA)
	defer rows.Close()
	CheckErr(err)

	var allData = make([]Data, 0)
	for rows.Next() {
		var data = Data{}
		err = rows.Scan(&data.Aid, &data.View, &data.Danmaku, &data.Reply, &data.Favorite,
			&data.Coin, &data.Share, &data.NowRank, &data.HisRank, &data.Like, &data.NoReprint, &data.Copyright)
		CheckErr(err)
		allData = append(allData, data)
	}
	return allData
}
