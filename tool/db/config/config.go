package config

type Config struct {
	UseMysql bool   `json:"use_mysql"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}
