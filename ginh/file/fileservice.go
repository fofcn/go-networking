package file

import (
	"go-networking/log"

	"gorm.io/gorm"
)

type FileService struct {
	db *gorm.DB
}

func NewFileService(db *gorm.DB) *FileService {
	return &FileService{
		db: db,
	}
}

func (f *FileService) GetFile(fileId uint) (*FileInfo, error) {
	var file FileInfo
	err := f.db.Where("id = ?", fileId).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (f *FileService) ListFile(cmd *ListFileCmd, userId uint) (*ListFilePageDto, error) {
	files, err := f.queryFileList(f.db, userId, cmd.FileId, cmd.Page, cmd.PageSize)
	if err != nil {
		return nil, err
	}

	var dtos []*ListFileDto
	for _, file := range files {
		dtos = append(dtos, &ListFileDto{
			FileId:      file.ID,
			FileName:    file.Fileame,
			OrgFilename: file.OrgFilename,
			ContentType: file.ContentType,
		})
	}

	return &ListFilePageDto{
		Pagenation: cmd.Pagenation,
		List:       dtos,
	}, nil
}

func (f *FileService) queryFileList(db *gorm.DB, userId uint, fileId uint, page int, pageSize int) ([]*FileInfo, error) {
	var files []*FileInfo

	query := db.Where("user_id = ?", userId).Offset((page - 1) * pageSize).Order("id DESC").Find(&files)

	// 如果fileId不为0，添加额外的查询条件
	if fileId != 0 {
		query = query.Where("file_id = ?", fileId)
	}

	if err := query.Error; err != nil {
		// 更细致的错误处理，可以考虑添加日志记录等
		log.Errorf("query file list error, %v", err)
		return nil, err
	}

	return files, nil
}

func (f *FileService) UploadFile(userId uint, parentId uint, fileName string, contentType uint8, orgFilename string) (*FileInfo, error) {
	// file := &FileInfo{
	// 	UserId:      userId,
	// 	ParentId:    parentId,
	// 	Filename:    fileName,
	// 	ContentType: contentType,
	// 	OrgFilename: orgFilename,
	// 	FileType:    0,
	// }

	// return file, f.db.Create(file).Error
	return nil, nil
}
