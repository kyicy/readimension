package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kyicy/readimension/model"
	"github.com/kyicy/readimension/utility"
	"github.com/labstack/echo-contrib/session"
	redistore "gopkg.in/boj/redistore.v1"
)

func parseFlag() (string, string) {
	var configFile string
	flag.StringVar(&configFile, "conf", "config.json", "json config file")
	var environment string
	flag.StringVar(&environment, "env", "production", "running environment")
	flag.Parse()
	return configFile, environment
}

type configStruct map[string]struct {
	Addr string `json:"addr"`
	Port string `json:"port"`

	SessionSecret string `json:"session_secret"`
	MySQL         struct {
		Addr     string `json:"addr"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
		Size     int    `json:"size"`
	} `json:"mysql"`

	Redis struct {
		Addr     string `json:"addr"`
		Port     string `json:"port"`
		Password string `json:"password"`
		Size     int    `json:"size"`
		DB       int    `json:"db"`
	} `json:"redis"`

	SMTP struct {
		Sender   string `json:"sender"`
		Password string `json:"password"`
		SMTP     string `json:"smtp"`
	}
}

func main() {
	configFile, env := parseFlag()
	file, err := os.Open(configFile)
	checkError(err)

	bytes, err := ioutil.ReadAll(file)
	checkError(err)

	var configObj configStruct
	json.Unmarshal(bytes, &configObj)
	envConfig := configObj[env]

	// Redis Session Store
	redisConfig := envConfig.Redis
	sessionStore, err := redistore.NewRediStore(redisConfig.Size, "tcp", ":"+redisConfig.Port, redisConfig.Password, []byte(envConfig.SessionSecret))
	checkError(err)
	sessionStore.SetMaxAge(100 * 24 * 3600)
	defer sessionStore.Close()

	// Redis Cache
	redisClient := redis.NewClient(&redis.Options{
		DB:       redisConfig.DB,
		PoolSize: redisConfig.Size,
		Addr:     fmt.Sprintf("%s:%s", redisConfig.Addr, redisConfig.Port),
	})
	defer redisClient.Close()
	utility.SetUpRedis(redisClient)

	// Mysql and Model
	mysqlConfig := envConfig.MySQL
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", mysqlConfig.User, mysqlConfig.Password, mysqlConfig.Addr, mysqlConfig.Database))
	checkError(err)
	if env != "production" {
		db.LogMode(true)
	}
	defer db.Close()
	model.LoadModel(db)

	// Setup Postman
	smtpConf := envConfig.SMTP
	utility.SetUpMailer(env, smtpConf.Sender, smtpConf.Password, smtpConf.SMTP)

	// Create Echo Server Instance
	e := createInstance(env)
	e.Use(session.Middleware(sessionStore))

	// Start the Server
	addr := fmt.Sprintf("%s:%s", envConfig.Addr, envConfig.Port)
	e.Logger.Fatal(e.Start(addr))
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
