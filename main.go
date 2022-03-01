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

	"github.com/gorilla/sessions"
	"github.com/kyicy/readimension/model"
	"github.com/kyicy/readimension/utility/config"
	"github.com/labstack/echo-contrib/session"
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
	err := os.MkdirAll(path.Join(workingPath, "books"), 0777)
	checkError(err)

	// cover folder
	err = os.MkdirAll(path.Join(workingPath, "covers"), 0777)
	checkError(err)

	// upload folder
	err = os.MkdirAll(path.Join(workingPath, "uploads"), 0777)
	checkError(err)

	file, err := os.Open(path.Join(workingPath, "config.json"))
	checkError(err)

	bytes, err := ioutil.ReadAll(file)
	checkError(err)

	err = json.Unmarshal(bytes, &config.Configuration)
	checkError(err)
	envConfig := config.Configuration[env]
	envConfig.WorkDir = workingPath
	config.SetENV(env)

	// Database Model
	dbPath := path.Join(workingPath, "readimension.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	checkError(err)
	model.LoadModel(db)

	// Create Echo Server Instance
	e := createInstance(env)
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(envConfig.SessionSecret))))

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
