package file

import "go-networking/ginh/common"

type ListFileCmd struct {
	common.Pagenation
	FileId uint `json:"file_id"`
}
