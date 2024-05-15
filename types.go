package main

type Data struct {
	RetrieveHeight string `json:"retrieve_height"`
	Data           []byte `json:"data"`
}

// Response 구조체 정의
type Response struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result,omitempty"`
	Error   string      `json:"error,omitempty"`
}
