package common

type CommonResp struct {
	Data    interface{} `json:"data"`
	Message interface{} `json:"message"`
}

var NoDataSuccessResposne = CommonResp{
	Data:    "OK",
	Message: "OK",
}
var NoDataFailureResposne = CommonResp{
	Data:    "Failed",
	Message: "Failed",
}

type Pagenation struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}
