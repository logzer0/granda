package main

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMaskImprSrc(n int) []byte {
	var src = rand.NewSource(time.Now().UnixNano())

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return b
}

func getCanonicalName(userSubmittedName string) (canonicalImage string) {
	canonicalImage = userSubmittedName
	if !strings.Contains(userSubmittedName, dockerHub) {
		canonicalImage = dockerHub + userSubmittedName
	}
	if !strings.Contains(userSubmittedName, colon) {
		canonicalImage = canonicalImage + latest
	}
	return
}

func getTagName(userSubmittedName string) (tagName string) {
	tagName = userSubmittedName
	if !strings.Contains(userSubmittedName, colon) {
		tagName = userSubmittedName + latest
	}
	return
}

func checkIfImageExists(name string) bool {
	imgs, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, img := range imgs {
		for _, eachName := range img.RepoTags {
			if name == eachName {
				return true
			}
		}
	}

	return false
}

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, nice meeting you. \n What are you doing here?"))
}

func getTimeout(time string) (timeout int) {
	var err error
	timeout = 10
	if time != "" {
		timeout, err = strconv.Atoi(time)
		if err != nil {
			panic(err)
		}
	}
	return
}

func getExternalPackages(code string) (out []string) {
	a := strings.Index(code, "(")
	b := strings.Index(code, "\"")
	c := strings.Index(code, ")")
	if a < b {
		z := code[a+1 : c]
		za := strings.Split(z, "\n")
		for _, v := range za {
			p := strings.TrimSpace(v)
			if strings.Contains(p, ".") {
				p = p[strings.Index(p, "\"")+1 : strings.LastIndex(p, "\"")]
				out = append(out, p)
			}
		}
	}
	return
}
