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

func NewCommonResp(data interface{}, message interface{}) CommonResp {
	return CommonResp{
		Data:    data,
		Message: message,
	}
}

type Pagenation struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}
