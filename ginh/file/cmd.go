package file

import "go-networking/ginh/common"

type ListFileCmd struct {
	common.Pagenation
	FileId uint `json:"file_id"`
}

type UploadFileCmd struct {
	UserId      uint   `json:"user_id"`      // 上传用户
	ParentId    uint   `json:"parent_id"`    // 上传文件父级
	Filename    string `json:"filename"`     // 上传文件名
	ContentType uint8  `json:"content_type"` // 上传文件类型
	OrgFilename string `json:"org_filename"` // 原始文件名
}
