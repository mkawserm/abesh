package model

type HTTPResponseObject struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Lang    string      `json:"lang"`
	Data    interface{} `json:"data"`
}
