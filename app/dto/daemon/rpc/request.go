package daemon

type BaseRequest map[string]interface{}

func NewBaseRequest(method string, params []interface{}) BaseRequest {
	req := make(BaseRequest, 4)
	req["jsonrpc"] = "1.0"
	req["id"] = "bitcoin"
	req["method"] = method
	req["params"] = params

	return req
}
