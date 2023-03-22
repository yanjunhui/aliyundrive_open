package aliyundrive_open

import "fmt"

// DriveInfo 云盘信息
type DriveInfo struct {
	Avatar         string        `json:"avatar"`
	Email          string        `json:"email"`
	Phone          string        `json:"phone"`
	Role           string        `json:"role"`
	Status         string        `json:"status"`
	Description    string        `json:"description"`
	Punishments    []interface{} `json:"punishments"`
	PunishFlagEnum int           `json:"punishFlagEnum"`
	UserId         string        `json:"user_id"`
	DomainId       string        `json:"domain_id"`
	UserName       string        `json:"user_name"`
	NickName       string        `json:"nick_name"`
	DefaultDriveId string        `json:"default_drive_id"`
	CreatedAt      int64         `json:"created_at"`
	UpdatedAt      int64         `json:"updated_at"`
	UserData       struct {
		BackUpConfig struct {
			手机备份 struct {
				FolderId      string `json:"folder_id"`
				PhotoFolderId string `json:"photo_folder_id"`
				SubFolder     struct {
				} `json:"sub_folder"`
				VideoFolderId string `json:"video_folder_id"`
			} `json:"手机备份"`
		} `json:"back_up_config"`
	} `json:"user_data"`
	PunishFlag int `json:"punish_flag"`
	ErrorInfo
}

// DriveInfo 获取云盘信息
func (a *Authorize) DriveInfo() (result DriveInfo, err error) {

	err = a.HttpPost(APIDriveInfo, map[string]string{}, &result)
	if err != nil {
		return result, err
	}

	if result.Code != "" {
		err = fmt.Errorf("获取云盘信息失败: %s", result.Message)
	}

	return result, err
}

type SpaceInfo struct {
	PersonalSpaceInfo struct {
		UsedSize  int64 `json:"used_size"`
		TotalSize int64 `json:"total_size"`
	} `json:"personal_space_info"`
	ErrorInfo
}

// DriveSpace 获取云盘空间信息
func (a *Authorize) DriveSpace() (result SpaceInfo, err error) {
	err = a.HttpPost(APISpaceInfo, nil, &result)
	if err != nil {
		return result, err
	}

	if result.Code != "" {
		err = fmt.Errorf("获取云盘空间信息失败: %s", result.Message)
	}

	return result, err
}
