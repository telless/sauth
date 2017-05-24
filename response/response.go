package response

import (
	"encoding/json"
	"sauth/models"
)

type BaseResponse interface {
	Serialize() string
}

type errorResponse struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

func NewErrorResponse(status string, data string) errorResponse {
	return errorResponse{Status: status, Data: data}
}

func (err *errorResponse) Serialize() string {
	jsonData, _ := json.Marshal(&err)
	return string(jsonData)
}

type successResponse struct {
	Status string `json:"status"`
	Data   models.User `json:"data"`
}

func NewSuccessResponse(status string, data models.User) successResponse {
	return successResponse{Status: status, Data: data}
}

func (success *successResponse) Serialize() string {
	jsonData, _ := json.Marshal(&success)
	return string(jsonData)
}
