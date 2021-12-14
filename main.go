package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/kyicy/readimension/model"
	"github.com/kyicy/readimension/utility/config"
	"github.com/labstack/echo-contrib/session"
	"github.com/michaeljs1990/sqlitestore"
)

func parseFlag() (string, string) {
	var environment string
	flag.StringVar(&environment, "env", "production", "running environment")

	var path string
	flag.StringVar(&path, "path", ".", "working path")

	flag.Parse()
	return environment, path
}

func main() {
	env, workingPath := parseFlag()

	// book folder
	os.MkdirAll(path.Join(workingPath, "books"), 0777)

	// cover folder
	os.MkdirAll(path.Join(workingPath, "covers"), 0777)

	// upload folder
	os.MkdirAll(path.Join(workingPath, "uploads"), 0777)

	file, err := os.Open(path.Join(workingPath, "config.json"))
	checkError(err)

	bytes, err := ioutil.ReadAll(file)
	checkError(err)

	json.Unmarshal(bytes, &config.Configuration)
	envConfig := config.Configuration[env]
	config.SetENV(env)

	// Session Store
	dbPath := path.Join(workingPath, "readimension.db")
	sessionStore, err := sqlitestore.NewSqliteStore(dbPath, "sessions", "/", 3600*24*365, []byte(envConfig.SessionSecret))
	checkError(err)
	defer sessionStore.Close()

	// Database Model
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	checkError(err)
	model.LoadModel(db)

	// Create Echo Server Instance
	e := createInstance(env)
	e.Use(session.Middleware(sessionStore))

	// Start the Server
	addr := fmt.Sprintf("%s:%s", envConfig.Addr, envConfig.Port)
	s := &http.Server{
		Addr:         addr,
		ReadTimeout:  20 * time.Minute,
		WriteTimeout: 20 * time.Minute,
	}
	e.Logger.Fatal(e.StartServer(s))
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
