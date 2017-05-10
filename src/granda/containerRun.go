package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func ContainerRunHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	path := r.RequestURI
	indx := strings.Index(r.RequestURI, "?")
	if indx > 0 {
		path = r.RequestURI[0:indx]
	}

	runResponse := make(chan string)
	go RunContainer(string(body), path, pathToImage[path], runResponse)

	func() {
		select {
		case status := <-runResponse:
			response := RunResponse{Status: status}
			reqResponse, _ := json.Marshal(response)
			w.Write(reqResponse)
			return
		}
	}()

}

func RunContainer(userMessage string, urlPath string, img Image, runResponse chan string) {
	begin := time.Now()

	done := make(chan bool)
	var commandString []string
	if len(userMessage) > 0 {
		commandString = append(commandString, userMessage)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: img.imageName,
		Cmd:   commandString,
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	timeoutContext, cancel := context.WithTimeout(ctx, time.Duration(img.Timeout)*time.Second)
	defer cancel()
	defer func() {
		if r := recover(); r != nil {
			if err := cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
				panic(err)
			}
			//save data to db
			StoreRuntimeToDB(img.funcName, urlPath, begin.String(), time.Since(begin).Seconds(), "Failure")
			runResponse <- "Operation Timed Out"
		}
	}()

	go func() {
		select {
		case <-done:
			//save data to db
			StoreRuntimeToDB(img.funcName, urlPath, begin.String(), time.Since(begin).Seconds(), "Success")
			runResponse <- "Success"
		case <-timeoutContext.Done():
			//Well, the timeout has been exceeded
		}
	}()

	if _, err := cli.ContainerWait(timeoutContext, resp.ID); err != nil {
		panic(err)
	} else {
		done <- true
	}

	if err := cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
		panic(err)
	}

}
