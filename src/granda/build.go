package main

import (
	"archive/tar"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
)

func CodeFuncHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var response FunctionCreatedResponse
	funcName := r.PostFormValue("name")
	runTime := r.PostFormValue("runtime")
	timeout := getTimeout(r.PostFormValue("time"))
	description := r.PostFormValue("description")
	code := r.PostFormValue("code")

	buildState, imageName, message := buildAnImage(runTime, funcName, code)
	response.Message = message
	if buildState {
		img := Image{funcName: funcName, imageName: imageName, Timeout: timeout, Description: description}
		response = CreateFunction(img)
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	w.Write(jsonResp)
}

func buildAnImage(runTime, funcName, code string) (status bool, imageName string, message string) {
	message = "Failed to create the image"
	fileMap := make(map[string]string)
	switch runTime {
	case go17, go16, go15:
		codeFile, codeContent := saveCodeFile(goUserSpacePath, funcName, code)
		fileMap[codeFile] = codeContent
		dockerFile, dockerFileContent := createGoDockerFile(goUserSpacePath, funcName, code, runTime)
		fileMap[dockerFile] = dockerFileContent
		tarPath := createTarFile(goUserSpacePath, funcName, fileMap)
		status, imageName = buildNewImage(funcName, tarPath)
		if status {
			message = "Success"
		}

	default:
		fmt.Println("Runtime unavailable")
	}
	return
}

func saveCodeFile(path, funcName, code string) (fileName, content string) {
	content = code
	fileName = "main.go"
	return
}

func createGoDockerFile(path, funcName, code, runTime string) (dockerFileName, content string) {
	//Set the base image
	switch runTime {
	case go17:
		content = "From golang:1.7\n"
	case go16:
		content = "From golang:1.6\n"
	case go15:
		content = "From golang:1.5\n"
	}

	//Create the dir
	content += "\n#Create the dir\nADD . /go/src/app\n"
	content += "\n#External packages\n"
	externalPackages := getExternalPackages(code)
	for _, eachPackage := range externalPackages {
		content += "RUN go get " + eachPackage + "\n"
	}

	content += "\n#Build and Install app \nRUN go install app\n\nENTRYPOINT [\"/go/bin/app\"]"
	dockerFileName = "Dockerfile"
	return
}

func createTarFile(path, funcName string, fileMap map[string]string) string {
	var err error
	var f *os.File

	tarDirPath := path + funcName + "/"
	tarPath := tarDirPath + funcName + ".tar"

	if _, err := os.Stat(tarDirPath); os.IsNotExist(err) {
		if err := os.Mkdir(tarDirPath, 0755); err != nil {
			panic(err)
		}
	}

	f, err = os.Create(tarPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	tw := tar.NewWriter(f)
	for fileName, fileContent := range fileMap {
		hdr := &tar.Header{
			Name: fileName,
			Mode: 0600,
			Size: int64(len(fileContent)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			log.Fatalln(err)
		}
		if _, err := tw.Write([]byte(fileContent)); err != nil {
			log.Fatalln(err)
		}
	}

	if err := tw.Close(); err != nil {
		panic(err)
	}

	return tarPath
}

func buildNewImage(imageName, tarPath string) (bool, string) {
	dockerBuildContext, err := os.Open(tarPath)
	defer dockerBuildContext.Close()

	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile", // optional, is the default
		Tags:       []string{strings.ToLower(imageName)},
	}

	buildResponse, err := cli.ImageBuild(context.Background(), dockerBuildContext, buildOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer buildResponse.Body.Close()

	var resp buildResp
	for {
		if err := json.NewDecoder(buildResponse.Body).Decode(&resp); err == io.EOF {
			break
		}
	}

	return strings.Contains(resp.Stream, "Success"), imageName
}
