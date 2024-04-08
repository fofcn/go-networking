package file

import "go-networking/ginh/common"

type ListFileDto struct {
	FileId      uint   `json:"file_id"`
	FileName    string `json:"file_name"`
	OrgFilename string `json:"org_filename"`
	ContentType uint8  `json:"content_type"`
}

type ListFilePageDto struct {
	common.Pagenation
	List []*ListFileDto `json:"list"`
}
