package main

import (
	"log"
	"net/http"
)

func main() {
	router.HandleFunc("/", HelloWorldHandler)

	router.HandleFunc("/createFunc", CreateFunctionHandler)
	router.HandleFunc("/createCodeFunc", CodeFuncHandler)

	router.HandleFunc("/funcInvocationCount", GetFuncInvocationCount)
	router.HandleFunc("/invocationDetails", InvocationDetails)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8881", router))
}
