package util

import "net/http"

type ApiResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

var ServerError = &ApiResponse{
	Code: http.StatusInternalServerError,
}
