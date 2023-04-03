package aliyundrive_open

import (
	"os"
	"strings"
)

// OrderSortedField 排序字段
type OrderSortedField string

const (
	OrderFieldCreated OrderSortedField = "created_at"
	OrderFieldUpdate  OrderSortedField = "updated_at"
	OrderFieldSize    OrderSortedField = "size"
	OrderFieldName    OrderSortedField = "name"
)

// OrderSortedDirection 排序方式
type OrderSortedDirection string

const (
	OrderSortedDirectionAsc  OrderSortedDirection = "ASC"
	OrderSortedDirectionDesc OrderSortedDirection = "DESC"
)

// FileCategory 返回文件类型分类
type FileCategory string

func (fc FileCategory) String() string {
	return string(fc)
}

const (
	FileCategoryVideo  FileCategory = "video"  // 视频
	FileCategoryAudio  FileCategory = "audio"  // 音频
	FileCategoryImage  FileCategory = "image"  // 图片
	FileCategoryDoc    FileCategory = "doc"    // 文档
	FileCategoryZip    FileCategory = "zip"    // 压缩包
	FileCategoryOthers FileCategory = "others" // 其他
)

// FileType 返回文件类型
type FileType string

const (
	FileTypeAll    FileType = "all"    // 所有
	FileTypeFile   FileType = "file"   // 文件
	FileTypeFolder FileType = "folder" // 目录
)

// ResponseFieldName 指定返回的字段类型
type ResponseFieldName string

func (rf ResponseFieldName) String() string {
	return string(rf)
}

const (
	ResponseFieldURL           ResponseFieldName = "url"
	ResponseFieldThumbnail     ResponseFieldName = "thumbnail"
	ResponseFieldVideoMetadata ResponseFieldName = "video_metadata"
)

// CheckNameMode 重命名时检查文件名模式
type CheckNameMode string

const (
	CheckNameModeRefuse     CheckNameMode = "refuse"      // 重名时拒绝创建文件
	CheckNameModeAutoRename CheckNameMode = "auto_rename" //自动重命名
	CheckNameModeIgnore     CheckNameMode = "ignore"      //允许重命名
)

// FileOption 文件列表参数
type FileOption struct {
	DriveID             string               `json:"drive_id"`              // 云盘ID(必填)
	ParentFileID        string               `json:"parent_file_id"`        // 目录ID(目录/必填)
	FileID              string               `json:"file_id,omitempty"`     // 文件ID(文件必填)
	Name                string               `json:"name"`                  // 文件名(重命名必填)
	Path                string               `json:"path"`                  // 文件完整路径(不包括 /root)
	ExpireSec           int64                `json:"expire_sec"`            // 下载链接有效期(链接必填)
	URLExpireSec        int64                `json:"url_expire_sec"`        // 视频播放地址有效期(播放必填)
	ToParentFileID      string               `json:"to_parent_file_id"`     // 移动到的目录ID(移动必填)
	CheckNameMode       CheckNameMode        `json:"check_name_mode"`       // 重命名时检查文件名模式(重命名)
	NewName             string               `json:"new_name"`              // 移动时重名时的新文件名(移动)
	Marker              string               `json:"marker"`                // 分页标记(目录)
	Limit               int64                `json:"limit"`                 // 分页大小(目录)
	OrderBy             OrderSortedField     `json:"order_by"`              // 排序字段(目录)
	OrderDirection      OrderSortedDirection `json:"order_direction"`       // 排序方式(目录)
	Category            string               `json:"category"`              // 指定返回的文件类型(目录/文件)
	Type                FileType             `json:"type"`                  // 指定返回文件还是目录, type不为空时, category参数无效(目录)
	VideoThumbnailTime  int64                `json:"video_thumbnail_time"`  // 视频预览时间 (单位:秒) (目录/文件)
	VideoThumbnailWidth int64                `json:"video_thumbnail_width"` // 视频预览宽度 (目录/文件)
	ImageThumbnailWidth int64                `json:"image_thumbnail_width"` // 视频预览图片宽度 (目录/文件)
	Fields              string               `json:"fields"`                // 只返回指定字段 (目录)
	//ParallelUpload      bool                 `json:"parallel_upload"`       // Deprecated: 并发上传已经停止支持
	PartInfoList []FileUpdatePartInfo `json:"part_info_list"` // 分片上传信息(上传)
	OpenFile     *os.File             `json:"-"`              // 文件流(上传)
	UploadID     string               `json:"upload_id"`      // 上传ID(上传)
}

// FileUpdatePartInfo 分片上传选项
type FileUpdatePartInfo struct {
	ParallelSha1Ctx ParallelSha1Ctx `json:"parallel_sha1_ctx"` // 分片sha1
	PartNumber      int64           `json:"part_number"`       // 分片序号

}

type ParallelSha1Ctx struct {
	PartOffset int64    `json:"part_offset"` // 分片偏移量
	PartSize   int64    `json:"part_size"`   // 分片大小
	H          []uint32 `json:"h"`           // 分片sha1
}

// NewFileCreateOption 创建文件参数
func NewFileCreateOption(parentFileID, name string) *FileOption {
	if parentFileID == "" {
		parentFileID = "root"
	}

	return &FileOption{
		ParentFileID:  parentFileID,
		Name:          name,
		CheckNameMode: "auto_rename",
	}
}

// NewFileUploadOption 创建文件上传参数
func NewFileUploadOption(parentFileID, name string, of *os.File) *FileOption {
	option := NewFileCreateOption(parentFileID, name)
	option.Type = FileTypeFile
	option.OpenFile = of
	return option
}

// NewFileVideoPlayInfoOption 创建获取视频播放信息参数
func NewFileVideoPlayInfoOption(fileID string) *FileOption {
	return &FileOption{
		FileID:       fileID,
		Category:     "live_transcoding",
		URLExpireSec: 14400,
	}
}

// NewFileTrashAndDeleteOption 创建文件删除参数
func NewFileTrashAndDeleteOption(fileID string) *FileOption {
	return &FileOption{
		FileID: fileID,
	}
}

// NewFileMoveAndCopyOption 创建文件复制/移动参数
func NewFileMoveAndCopyOption(fileID, toParentFileID string) *FileOption {
	return &FileOption{
		FileID:         fileID,
		ToParentFileID: toParentFileID,
		CheckNameMode:  "auto_rename",
	}
}

// NewFileDownloadURLOption 创建获取单个文件下载链接默认参数
func NewFileDownloadURLOption(fileID string) *FileOption {
	return &FileOption{
		FileID:    fileID,
		ExpireSec: 115200,
	}
}

// NewFileRenameOption 创建重命名参数
func NewFileRenameOption(fileID, newName string) *FileOption {
	return &FileOption{
		FileID: fileID,
		Name:   newName,
	}
}

// NewFileOption 创建获取单个文件默认参数
func NewFileOption(fileID string) *FileOption {
	return &FileOption{
		FileID: fileID,
	}
}

// NewFileOptionByPath 根据文件路径获取文件信息, 该接口暂为灰度测试接口
func NewFileOptionByPath(path string) *FileOption {
	return &FileOption{
		Path: path,
	}
}

// NewFilesOption 创建获取多个文件默认参数
func NewFilesOption(fileIDs []string) (options []*FileOption) {
	for _, id := range fileIDs {
		options = append(options, NewFileOption(id))
	}
	return
}

// NewFileListOption 创建默认文件列表参数
func NewFileListOption(parentFileID, marker string) *FileOption {
	return &FileOption{
		ParentFileID:   parentFileID,
		Marker:         marker,
		OrderBy:        OrderFieldName,
		OrderDirection: OrderSortedDirectionAsc,
		Limit:          100,
		URLExpireSec:   86400,
		Fields:         "*",
	}
}

// SetDriveID 设置目录ID
func (option *FileOption) SetDriveID(driveID string) *FileOption {
	option.DriveID = driveID
	return option
}

// SetParentFileID 设置目录ID
func (option *FileOption) SetParentFileID(parentFileID string) *FileOption {
	option.ParentFileID = parentFileID
	return option
}

// SetFileID 设置文件ID
func (option *FileOption) SetFileID(fileID string) *FileOption {
	option.FileID = fileID
	return option
}

// SetFilePath 设置文件ID
func (option *FileOption) SetFilePath(path string) *FileOption {
	option.Path = path
	return option
}

// SetName 设置文件名
func (option *FileOption) SetName(name string) *FileOption {
	option.Name = name
	return option
}

// SetExpireSec 设置下载链接有效期
func (option *FileOption) SetExpireSec(expireSec int64) *FileOption {
	option.ExpireSec = expireSec
	return option
}

// SetURLExpireSec 设置视频播放地址有效期
func (option *FileOption) SetURLExpireSec(urlExpireSec int64) *FileOption {
	option.URLExpireSec = urlExpireSec
	return option
}

// SetNewName 设置移动时重名时的新文件名
func (option *FileOption) SetNewName(newName string) *FileOption {
	option.NewName = newName
	return option
}

// SetCheckNameMode 设置重命名时检查文件名模式
func (option *FileOption) SetCheckNameMode(checkNameMode CheckNameMode) *FileOption {
	option.CheckNameMode = checkNameMode
	return option
}

// SetMarker 设置分页标记
func (option *FileOption) SetMarker(marker string) *FileOption {
	option.Marker = marker
	return option
}

// SetLimit 设置分页大小
func (option *FileOption) SetLimit(limit int64) *FileOption {
	option.Limit = limit
	return option
}

// SetOrderBy 设置排序字段
func (option *FileOption) SetOrderBy(orderBy OrderSortedField) *FileOption {
	option.OrderBy = orderBy
	return option
}

// SetOrder 设置排序方式
func (option *FileOption) SetOrder(direction OrderSortedDirection) *FileOption {
	option.OrderDirection = direction
	return option
}

// SetCategory 设置返回文件类型分类
func (option *FileOption) SetCategory(category []FileCategory) *FileOption {
	option.Category = joinCustomString(category, ",")
	return option
}

// SetType 设置返回文件类型
func (option *FileOption) SetType(fileType FileType) *FileOption {
	option.Type = fileType
	return option
}

// SetVideoThumbnailTime 设置视频预览时间
func (option *FileOption) SetVideoThumbnailTime(time int64) *FileOption {
	option.VideoThumbnailTime = time
	return option
}

// SetFields 设置返回字段
func (option *FileOption) SetFields(fields []string) *FileOption {
	option.Fields = strings.Join(fields, ",")
	return option
}

// SetThumbnailWidth 设置视频预览宽度
func (option *FileOption) SetThumbnailWidth(width int64) *FileOption {
	option.VideoThumbnailWidth = width
	option.ImageThumbnailWidth = width
	return option
}

// SetResponseFields 设置返回字段
func (option *FileOption) SetResponseFields(fields []ResponseFieldName) *FileOption {
	option.Fields = joinCustomString(fields, ",")
	return option
}

// SetUploadOpenFile 设置上传数据流
func (option *FileOption) SetUploadOpenFile(f *os.File) *FileOption {
	option.OpenFile = f
	return option
}
