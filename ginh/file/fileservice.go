package file

import "gorm.io/gorm"

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
	var files []*FileInfo
	err := f.db.Where("user_id = ? and file_id = ?", userId, cmd.FileId).Offset((cmd.Page - 1) * cmd.PageSize).Order("id DESC").Find(&files).Error
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
