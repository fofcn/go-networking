package common

type CommonResp struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

var NoDataSuccessResposne = CommonResp{
	Data:    "OK",
	Message: "OK",
}
var NoDataFailureResposne = CommonResp{
	Data:    "Failed",
	Message: "Failed",
}
