package daemon

type BaseResponse struct {
	Result string      `json:"result"`
	Error  interface{} `json:"error"`
	Id     string      `json:"id"`
}
