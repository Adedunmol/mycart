package util

type APIResponse struct {
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
	Status  string      `json:"status"`
}
