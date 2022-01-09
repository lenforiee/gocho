package globals

import (
	"encoding/json"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/lenforiee/gocho/logger"
	"github.com/oschwald/maxminddb-golang"
	"gopkg.in/redis.v5"
)

type Conf struct {
	Port int    `json:"port"`
	DSN  string `json:"mysql_dsn"`
}

var (
	DbConn    *sqlx.DB
	RedisConn *redis.Client
	Config    *Conf
	IPDb      *maxminddb.Reader
)

func ReadConfig() {
	dat, err := os.ReadFile("config.json")
	if err != nil {
		logger.Error("There was issue reading config file!")
		panic(err)
	}
	err = json.Unmarshal(dat, &Config)
	if err != nil {
		logger.Error("There was an error while parsing config!")
		panic(err)
	}
}

func InitialiseConnections() {
	db, err := sqlx.Open("mysql", Config.DSN+"?parseTime=true&allowNativePasswords=true")
	if err != nil {
		logger.Error("There was error connecting to mysql!")
		panic(err)
	} // pass the pointer.
	DbConn = db
	logger.Info("Connection with mysql established!")

	red := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	RedisConn = red
	logger.Info("Connection with redis established!")

	ip, err := maxminddb.Open("geoloc.mmdb")
	if err != nil {
		logger.Error("There was error connecting to IP database!")
		panic(err)
	}
	IPDb = ip
	logger.Info("Connection with IP database established!")
}
