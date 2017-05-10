package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func GetFuncInvocationCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var err error
	rResp := InvocationCountResponse{Status: true}
	duration := 10
	functionName := r.URL.Query().Get("funcname")
	if len(functionName) == 0 {
		rResp.Status = false
	}

	durationParam := r.URL.Query().Get("duration")
	if len(durationParam) != 0 {
		duration, err = strconv.Atoi(durationParam)
		if err != nil {
			rResp.Status = false
		}
	}

	if rResp.Status {
		rResp.Count = getInvocationCount(duration)
	}

	response, err := json.Marshal(rResp)
	if err != nil {
		panic(err)
	}

	w.Write(response)
}

func InvocationDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var resp invocationDetails
	functionName := r.URL.Query().Get("function")
	if len(functionName) != 0 {
		resp.Data = getInvocationDetailsFromDB(functionName)
	}

	response, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	w.Write(response)
}
