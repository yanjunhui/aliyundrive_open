package aliyundrive_open

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type FileList struct {
	Items      []FileInfo `json:"items"`
	NextMarker string     `json:"next_marker"`
	ErrorInfo
}

type FileInfo struct {
	Trashed            bool      `json:"trashed"`
	DriveId            string    `json:"drive_id"`
	FileId             string    `json:"file_id"`
	Category           string    `json:"category,omitempty"`
	ContentHash        string    `json:"content_hash,omitempty"`
	ContentHashName    string    `json:"content_hash_name,omitempty"`
	ContentType        string    `json:"content_type,omitempty"`
	Crc64Hash          string    `json:"crc64_hash,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	DomainId           string    `json:"domain_id"`
	DownloadUrl        string    `json:"download_url,omitempty"` // Deprecated: download_url 即将废弃
	EncryptMode        string    `json:"encrypt_mode"`
	FileExtension      string    `json:"file_extension,omitempty"`
	Hidden             bool      `json:"hidden"`
	MimeType           string    `json:"mime_type,omitempty"`
	Name               string    `json:"name"`
	ParentFileId       string    `json:"parent_file_id"`
	PunishFlag         int       `json:"punish_flag,omitempty"`
	Size               int64     `json:"size,omitempty"`
	Starred            bool      `json:"starred"`
	Status             string    `json:"status"`
	Thumbnail          string    `json:"thumbnail,omitempty"`
	Type               FileType  `json:"type"`
	UpdatedAt          time.Time `json:"updated_at"`
	Url                string    `json:"url,omitempty"`
	UserMeta           string    `json:"user_meta,omitempty"`
	SyncFlag           bool      `json:"sync_flag,omitempty"`
	VideoMediaMetadata struct {
		Duration              string `json:"duration"`
		Height                int    `json:"height"`
		VideoMediaAudioStream []struct {
			BitRate       string `json:"bit_rate"`
			ChannelLayout string `json:"channel_layout"`
			Channels      int    `json:"channels"`
			CodeName      string `json:"code_name"`
			Duration      string `json:"duration"`
			SampleRate    string `json:"sample_rate"`
		} `json:"video_media_audio_stream"`
		VideoMediaVideoStream []struct {
			Bitrate  string `json:"bitrate"`
			Clarity  string `json:"clarity"`
			CodeName string `json:"code_name"`
			Duration string `json:"duration"`
			Fps      string `json:"fps"`
		} `json:"video_media_video_stream"`
		Width int `json:"width"`
	} `json:"video_media_metadata,omitempty"`
	ExFieldsInfo struct {
	} `json:"ex_fields_info,omitempty"`
	ErrorInfo
}

// FileList  获取文件列表
func (a *Authorize) FileList(option *FileOption) (result FileList, err error) {
	if option == nil {
		option = NewFileListOption(a.DriveID, "root", "")
	}

	err = a.HttpPost(APIList, option, &result)
	if err != nil {
		return result, err
	}

	if result.Code != "" {
		err = fmt.Errorf("获取文件列表失败: %s", result.Message)
	}

	return result, err
}

// File 获取文件信息
func (a *Authorize) File(option *FileOption) (result FileInfo, err error) {
	if option == nil {
		return result, fmt.Errorf("option is nil")
	}

	err = a.HttpPost(APIFile, option, &result)
	if err != nil {
		return result, err
	}

	if result.Code != "" {
		err = fmt.Errorf("获取文件信息失败: %s", result.Message)
	}

	return result, err
}

// Files 批量获取文件信息
func (a *Authorize) Files(option []*FileOption) (result FileList, err error) {
	if len(option) == 0 {
		return result, fmt.Errorf("option is nil")
	}

	err = a.HttpPost(APIFiles, map[string][]*FileOption{
		"file_list": option,
	}, &result)
	if err != nil {
		return result, err
	}

	if result.Code != "" {
		err = fmt.Errorf("获取文件信息失败: %s", result.Message)
	}

	return result, err
}

type FileDownloadURL struct {
	URL        string    `json:"url"`
	Expiration string    `json:"expiration"`
	ExpireTime time.Time `json:"expire_time"`
	ErrorInfo
}

// FileDownloadURL 获取文件下载链接
func (a *Authorize) FileDownloadURL(option *FileOption) (result FileDownloadURL, err error) {
	if option == nil {
		return result, fmt.Errorf("option is nil")
	}

	err = a.HttpPost(APIFileDownload, option, &result)
	if err != nil {
		return result, err
	}

	if result.Code != "" {
		err = fmt.Errorf("获取文件下载信息失败: %s", result.Message)
	}

	return result, err
}

// FileRename 重命名文件
func (a *Authorize) FileRename(option *FileOption) (result FileInfo, err error) {
	if option == nil {
		return result, fmt.Errorf("option is nil")
	}

	err = a.HttpPost(APIFileUpdate, option, &result)
	if err != nil {
		return result, err
	}

	if result.Code != "" {
		err = fmt.Errorf("重命名失败: %s", result.Message)
	}

	return result, err
}

type FileVideoPlayInfo struct {
	DriveId              string `json:"drive_id"`
	FileId               string `json:"file_id"`
	VideoPreviewPlayInfo struct {
		Category string `json:"category"`
		Meta     struct {
			Duration float64 `json:"duration"`
			Width    int     `json:"width"`
			Height   int     `json:"height"`
		} `json:"meta"`
		LiveTranscodingTaskList []struct {
			TemplateId     string `json:"template_id"`
			TemplateName   string `json:"template_name"`
			TemplateWidth  int    `json:"template_width"`
			TemplateHeight int    `json:"template_height"`
			Status         string `json:"status"`
			Stage          string `json:"stage"`
			Url            string `json:"url"`
		} `json:"live_transcoding_task_list"`
	} `json:"video_preview_play_info"`
	ErrorInfo
}

// FileVideoPlayInfo 获取视频转码播放信息
func (a *Authorize) FileVideoPlayInfo(option *FileOption) (result FileVideoPlayInfo, err error) {
	err = a.HttpPost(APIFileVideoPlayInfo, option, &result)
	if err != nil {
		return result, err
	}

	if result.Code != "" {
		err = fmt.Errorf("获取视频文件转码播放信息失败: %s", result.Message)
	}

	return result, err
}

type FileMoveCopyDelTask struct {
	DriveID     string `json:"drive_id"`
	FileID      string `json:"file_id"`
	AsyncTaskID string `json:"async_task_id"`
	Exist       bool   `json:"exist"`
	ErrorInfo
}

// FileMove 移动文件
func (a *Authorize) FileMove(option *FileOption) (result FileMoveCopyDelTask, err error) {
	return a.FileMoveAndCopy(option, true)
}

// FileCopy 复制文件
func (a *Authorize) FileCopy(option *FileOption) (result FileMoveCopyDelTask, err error) {
	return a.FileMoveAndCopy(option, false)
}

// FileMoveAndCopy  移动/复制文件
func (a *Authorize) FileMoveAndCopy(option *FileOption, isMove bool) (result FileMoveCopyDelTask, err error) {
	file, err := a.File(option)
	if err != nil {
		return result, err
	}

	apiURL := APIFileCopy
	if isMove {
		apiURL = APIFileMove
	}

	newName := strings.Join([]string{file.Name, file.FileId[len(file.FileId)-8:]}, "_")
	option.SetNewName(newName)

	err = a.HttpPost(apiURL, option, &result)
	if err != nil {
		return result, err
	}

	if result.Code != "" {
		err = fmt.Errorf("移动文件失败: %s", result.Message)
	}

	return result, err
}

type FileCreate struct {
	DriveId      string `json:"drive_id"`
	FileId       string `json:"file_id"`
	ParentFileId string `json:"parent_file_id"`
	FileName     string `json:"file_name"`

	Trashed         interface{} `json:"trashed"`
	Name            interface{} `json:"name"`
	Thumbnail       interface{} `json:"thumbnail"`
	Type            string      `json:"type"`
	Category        interface{} `json:"category"`
	Hidden          interface{} `json:"hidden"`
	Status          interface{} `json:"status"`
	Description     interface{} `json:"description"`
	Meta            interface{} `json:"meta"`
	Url             interface{} `json:"url"`
	Size            interface{} `json:"size"`
	Starred         interface{} `json:"starred"`
	Available       interface{} `json:"available"`
	Exist           interface{} `json:"exist"`
	UserTags        interface{} `json:"user_tags"`
	MimeType        interface{} `json:"mime_type"`
	FileExtension   interface{} `json:"file_extension"`
	RevisionId      string      `json:"revision_id"`
	ContentHash     interface{} `json:"content_hash"`
	ContentHashName interface{} `json:"content_hash_name"`
	EncryptMode     string      `json:"encrypt_mode"`
	DomainId        string      `json:"domain_id"`
	DownloadUrl     interface{} `json:"download_url"`
	UserMeta        interface{} `json:"user_meta"`
	ContentType     interface{} `json:"content_type"`
	CreatedAt       interface{} `json:"created_at"`
	UpdatedAt       interface{} `json:"updated_at"`
	LocalCreatedAt  interface{} `json:"local_created_at"`
	LocalModifiedAt interface{} `json:"local_modified_at"`
	TrashedAt       interface{} `json:"trashed_at"`
	PunishFlag      interface{} `json:"punish_flag"`
	UploadId        string      `json:"upload_id"`
	Location        string      `json:"location"`
	RapidUpload     bool        `json:"rapid_upload"`
	PartInfoList    []struct {
		Etag        interface{} `json:"etag"`
		PartNumber  int         `json:"part_number"`
		PartSize    interface{} `json:"part_size"`
		UploadUrl   string      `json:"upload_url"`
		ContentType string      `json:"content_type"`
	} `json:"part_info_list"`

	ErrorInfo
}

// FileCreate 创建文件
func (a *Authorize) FileCreate(option *FileOption) (result FileCreate, err error) {
	option.SetType(FileTypeFile)
	return a.fileAndFolderCreate(option)
}

// FolderCreate 创建目录
func (a *Authorize) FolderCreate(option *FileOption) (result FileCreate, err error) {
	option.SetType(FileTypeFolder)
	return a.fileAndFolderCreate(option)
}

// fileAndFolderCreate 创建文件和目录
func (a *Authorize) fileAndFolderCreate(option *FileOption) (result FileCreate, err error) {

	log.Printf("创建文件选项信息: %+v\n", option)

	err = a.HttpPost(APIFileCreate, option, &result)
	if err != nil {
		return result, err
	}

	if result.Code != "" {
		err = fmt.Errorf("创建 %s 失败: %s", option.Type, result.Message)
	}

	return result, err
}

const DefaultPartSize int64 = 1024 * 1024 * 64

// FileUpload 上传文件
func (a *Authorize) FileUpload(option *FileOption) (result FileInfo, err error) {
	if option.OpenFile == nil {
		return result, fmt.Errorf("OpenFile is nil")
	}
	defer option.OpenFile.Close()

	//获取文件分片信息
	option.PartInfoList, err = SplitFile(option)
	if err != nil {
		return result, err
	}

	//创建文件
	creatResp, err := a.FileCreate(option)
	if err != nil {
		return result, err
	}

	//上传文件(串行)
	httpClient := new(http.Client)
	for index, part := range creatResp.PartInfoList {
		size := option.PartInfoList[index].ParallelSha1Ctx.PartSize
		req, err := http.NewRequest("PUT", part.UploadUrl, io.LimitReader(option.OpenFile, size))
		if err != nil {
			return result, err
		}

		res, err := httpClient.Do(req)
		if err != nil {
			return result, err
		}
		res.Body.Close()
		log.Printf("上传文件分片: %s, 返回状态: %s\n", part.UploadUrl, res.Status)
	}

	//完成
	err = a.HttpPost(APIFileComplete, map[string]string{
		"file_id":   creatResp.FileId,
		"drive_id":  creatResp.DriveId,
		"upload_id": creatResp.UploadId,
	}, &result)

	if err != nil {
		log.Printf("完成文件返回信息: %+v\n", result)
		return result, err
	}

	if result.Code != "" {
		err = fmt.Errorf("上传文件失败: %s", result.Message)
	}

	return result, err
}

// FileTrash 放入回收站
func (a *Authorize) FileTrash(option *FileOption) (result FileMoveCopyDelTask, err error) {
	if option == nil {
		return result, fmt.Errorf("option is nil")
	}

	err = a.HttpPost(APIFileTrash, option, &result)
	if err != nil {
		return result, err
	}

	if result.Code != "" {
		err = fmt.Errorf("文件放入回收站失败: %s", result.Message)
	}

	return result, err
}

// FileDelete 删除文件
func (a *Authorize) FileDelete(option *FileOption) (result FileMoveCopyDelTask, err error) {
	if option == nil {
		return result, fmt.Errorf("option is nil")
	}

	err = a.HttpPost(APIFileDelete, option, &result)
	if err != nil {
		return result, err
	}

	if result.Code != "" {
		err = fmt.Errorf("删除文件失败: %s", result.Message)
	}

	return result, err
}

// FileReplaceName 批量替换文件名内指定字符(官方接口二次封装), 支持单文件和目录内所有子文件
func (a *Authorize) FileReplaceName(fileID, old, new string) error {

	//查询文件信息
	fileOption := NewFileOption(a.DriveID, fileID)
	file, err := a.File(fileOption)
	if err != nil {
		return err
	}

	errFileIDs := make([]string, 0)
	if file.IsDir() {
		marker := "first"
		listOption := NewFileListOption(a.DriveID, fileID, "")
		for marker != "" {
			if marker == "first" {
				marker = ""
			}
			listOption.SetMarker(marker)
			list, err := a.FileList(listOption)
			if err != nil {
				return err
			}
			marker = list.NextMarker

			for _, f := range list.Items {
				if strings.Contains(f.Name, old) {
					newName := strings.Replace(f.Name, old, new, -1)
					option := NewFileRenameOption(a.DriveID, f.FileId, newName)
					_, err := a.FileRename(option)
					if err != nil {
						errFileIDs = append(errFileIDs, strings.Join([]string{f.FileId, err.Error()}, ":"))
					}
				}
			}
		}

		if len(errFileIDs) > 0 {
			err = fmt.Errorf("失败信息: %s", strings.Join(errFileIDs, ","))
		}

		return err
	}

	if strings.Contains(file.Name, old) {
		newName := strings.Replace(file.Name, old, new, -1)
		option := NewFileRenameOption(a.DriveID, fileID, newName)
		_, err := a.FileRename(option)
		return err
	}

	return fmt.Errorf("文件(%s)名字(\"%s\")中不包含\"%s\"", file.FileId, file.Name, old)
}

func (f *FileInfo) IsDir() bool {
	return f.Type == FileTypeFolder
}
