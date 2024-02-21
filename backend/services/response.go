package services

type Msg struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
}
