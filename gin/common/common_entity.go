package common

type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

var NoDataSuccessResposne = SuccessResponse{
	Data:    "OK",
	Message: "OK",
}
var NoDataFailureResposne = SuccessResponse{
	Data:    "Failed",
	Message: "Failed",
}
