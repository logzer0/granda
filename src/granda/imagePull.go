package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/docker/docker/api/types"
)

func CreateFunctionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	timeout := int(10)
	var err error
	funcName := r.PostFormValue("name")
	description := r.PostFormValue("description")
	imageName := r.PostFormValue("image")
	timeout = getTimeout(r.PostFormValue("time"))

	img := Image{funcName: funcName, Description: description, imageName: imageName, Timeout: timeout}

	response := CreateFunction(img)

	jsonResp, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	w.Write(jsonResp)
}

func CreateFunction(img Image) FunctionCreatedResponse {

	funcPath := string(RandStringBytesMaskImprSrc(8))
	routerPath := "/granda/" + img.funcName + funcPath

	canonicalImage := getCanonicalName(img.imageName)
	tagName := getTagName(img.imageName)

	response := FunctionCreatedResponse{Status: true, Path: prefix + routerPath}

	if !checkIfImageExists(tagName) {
		fmt.Println("Image unavailable, so pulling it from the repo ", tagName)
		response.Message = "Pulling the image from the docker repo. The function will be ready for use in 10 mins"
		_, err := cli.ImagePull(ctx, canonicalImage, types.ImagePullOptions{})
		if err != nil {
			panic(err)
			response.Status = false
		}
	}

	router.HandleFunc(routerPath, ContainerRunHandler)
	pathToImage[routerPath] = img

	//Save the info to the DB
	StoreFunctionToDB(img.funcName, img.imageName, routerPath, img.Timeout)

	return response
}
