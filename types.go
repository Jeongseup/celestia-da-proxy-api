package main

import "encoding/json"

type Data struct {
	// RetrieveHeight string          `json:"retrieve_height"`
	Data         []byte          `json:"data"`
	NamespaceKey string          `json:"namespace_key"`
	MetaData     json.RawMessage `json:"metadata"`
}

// Response 구조체 정의
type Response struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result,omitempty"`
	Error   string      `json:"error,omitempty"`
}
