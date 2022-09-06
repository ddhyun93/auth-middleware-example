package dto

type Msg struct {
	Result string      `json:"result"`
	Error  string      `json:"error"`
	Data   interface{} `json:"data"`
}
