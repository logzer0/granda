package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
)

var router *mux.Router
var ctx context.Context
var cli *client.Client
var pathToImage map[string]Image
var db *sql.DB

const (
	dockerHub        = "docker.io/library/"
	prefix           = "http://localhost:8881"
	latest           = ":latest"
	colon            = ":"
	dbPath           = "/repos/granda/src/sqlite/granda.db"
	userSpacePath    = "userspace/"
	goUserSpacePath  = userSpacePath + "golang/"
	py3UserSpacePath = userSpacePath + "py3/"
	go17             = "Go1.7"
	go16             = "Go1.6"
	go15             = "Go1.5"
	py3              = "Python3"
	py27             = "Python2.7"
)

func init() {
	var err error
	cli, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	ctx = context.Background()
	pathToImage = make(map[string]Image)
	db = initDB(dbPath)
	loadActiveState()
	createUserSpace()
}

//createUserSpace creates the user space for the runtimes
func createUserSpace() {
	if err := os.MkdirAll(userSpacePath, 0755); err != nil {
		panic(err)
	}

	if _, err := os.Stat(goUserSpacePath); os.IsNotExist(err) {
		if err := os.Mkdir(goUserSpacePath, 0755); err != nil {
			panic(err)
		}
	}

	if _, err := os.Stat(py3UserSpacePath); os.IsNotExist(err) {
		if err := os.Mkdir(py3UserSpacePath, 0755); err != nil {
			panic(err)
		}
	}
}

//loadActiveState loads the previous configuration from the DB
func loadActiveState() {
	getFunctionsFromDB()
	router = mux.NewRouter()
	for k, v := range pathToImage {
		router.HandleFunc(k, ContainerRunHandler)
		fmt.Println(k, "  ", v)
	}

}

type buildResp struct {
	Stream string `json:"stream"`
}
