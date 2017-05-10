package main

type FunctionCreatedResponse struct {
	Status  bool   `json:"status"`
	Path    string `json:"path"`
	Message string `json:"message"`
}

type RunResponse struct {
	Status string `json:"status"`
}

type InvocationCountResponse struct {
	Status bool `json:"status"`
	Count  int  `json:"count"`
}

type Image struct {
	funcName    string
	Description string
	imageName   string
	Timeout     int
}

type invocationDetails struct {
	Data [][]string `json:"data"`
}
