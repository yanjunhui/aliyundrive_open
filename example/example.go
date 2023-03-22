//登录授权相关方法示例

package main

import (
	"fmt"
	"github.com/yanjunhui/aliyundrive_open"
	"log"
	"os"
	"path/filepath"
	"time"
)

var client = aliyundrive_open.NewClient("ClientID", "ClientSecret")

// GetQRCode 获取登录二维码. 直接打开返回的 qrCodeUrl 就可以看到二维码.
// sid 参数用于后续 QrCodeStatus 方法获取扫码状态
func GetQRCode() (result aliyundrive_open.AuthorizeQRCode, err error) {
	// 生成请求选项. 默认配置选项.自定义参数可以再通过其他 Set 方法设置值
	option := aliyundrive_open.NewDefaultSingleAuthorizeOption()
	return client.QRCode(option)
}

// CheckQrCodeStatus 检查二维码状态. 通过 QRCode方法返回的 sid 参数获取二维码状态.
// 扫码成功后,返回 authCode 用于最后的登录授权获取 access_token 和 refresh_token
func CheckQrCodeStatus(sid string) (result aliyundrive_open.AuthorizeQRCodeStatus, err error) {
	return client.QrCodeStatus(sid)
}

// Auth 登录授权.
// 通过 QrCodeStatus 方法返回的 authCode 参数获取 access_token 和 refresh_token
func Auth(authCode string) (result aliyundrive_open.Authorize, err error) {
	return client.Authorize(authCode)
}

// RefreshToken 刷新 access_token
// 通过 Auth 方法返回的 refresh_token 参数刷新 access_token
func RefreshToken(refreshToken string) (result aliyundrive_open.Authorize, err error) {
	return client.RefreshToken(refreshToken)
}

// 完整的登录授权流程
func LoginQRCode() (result aliyundrive_open.Authorize, err error) {
	//1. 获取登录二维码
	qrCode, err := GetQRCode()
	if err != nil {
		log.Printf("获取二维码失败: %s\n", err)
		return result, err
	}

	//2. 二维码可以通过任意方式加载图片展示给用户. 这里我们就直接通过浏览器打开以下链接
	log.Printf("点击或者复制以下链接通过浏览器打开\n%s\n", qrCode.QrCodeUrl)

	//3. 循环检查二维码状态. 这里可以由前端主动发起请求查询, 也可以由后端主动轮询
	//这里我们就直接使用后端轮询的方式
	//确认登录成功后, 会获得authCode
	authCode := ""
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			status, err := CheckQrCodeStatus(qrCode.Sid)
			if err != nil {
				log.Printf("获取二维码状态失败: %s\n", err)
				return result, err
			}

			//有三个状态, 我们可以根据状态做不同的处理
			switch status.Status {
			case "WaitLogin": // 二维码未扫描
				continue
			case "ScanSuccess": // 二维码已扫描
				log.Println("二维码已扫描,等待授权确认")
			case "LoginSuccess": // 二维码已确认
				log.Println("二维码已确认")
				authCode = status.AuthCode
			}
		}
		//如果已经获取到了 authCode, 就跳出循环
		if authCode != "" {
			fmt.Println("authCode: ", authCode)
			break
		}
	}

	//4. 登录授权
	//通过 authCode 获取 access_token 和 refresh_token
	authorize, err := Auth(authCode)
	if err != nil {
		log.Printf("登录授权失败: %s\n", err)
		return result, err
	}

	log.Printf("登录授权成功\naccess_token: %s\n\nrefresh_token: %s\n driver_id: %s\n, 过期时间: %s\n", authorize.AccessToken, authorize.RefreshToken, authorize.DriveID, authorize.ExpiresTime.String())

	log.Printf("稍等 3 秒钟, 刷新 access_token\n")

	time.Sleep(3 * time.Second)

	//5. 刷新 access_token
	//通过 refresh_token 刷新 access_token
	authorize, err = RefreshToken(authorize.RefreshToken)
	if err != nil {
		log.Printf("刷新Token失败: %s\n", err)
		return result, err
	}

	log.Printf("刷新Token成功\naccess_token: %s\n\nrefresh_token: %s, 过期时间: %s\n", authorize.AccessToken, authorize.RefreshToken, authorize.ExpiresTime.String())
	return authorize, nil
}

// GetDriveInfo 通过 access_token 获取云盘信息
// 这里将获取到该用户云盘的 drive_id, 以后的每一个操作得需要
// 其实 Auth 方法已经集成 DriveInfo 方法得到了 drive_id
// 这里仅做示例
func GetDriveInfo(authorize aliyundrive_open.Authorize) {
	driveInfo, err := authorize.DriveInfo()
	if err != nil {
		log.Printf("获取云盘信息失败: %s\n", err)
		return
	}

	log.Printf("获取云盘信息成功, drive_id: %s\n", driveInfo.DefaultDriveId)
}

// GetDrivesSpace 获取空间使用情况
func GetDrivesSpace(authorize aliyundrive_open.Authorize) {
	driveSpace, err := authorize.DriveSpace()
	if err != nil {
		log.Printf("获取云盘空间使用情况失败: %s\n", err)
		return
	}

	log.Printf("获取云盘空间使用情况成功, 总空间: %d GB, 已使用空间: %d GB\n", driveSpace.PersonalSpaceInfo.TotalSize/1024/1024/1024/1024, driveSpace.PersonalSpaceInfo.UsedSize/1024/1024/1024/1024)
}

// GetFileList 获取文件列表
// 这里我们先获取根目录下的文件列表
// 这里我们使用了 NewFileListOption 方法来生成请求选项. 其余个性化参数可以通过 Set 方法设置
func GetFileList(authorize aliyundrive_open.Authorize, parentID string) {
	if parentID == "" {
		parentID = "root"
	}
	option := aliyundrive_open.NewFileListOption(authorize.DriveID, parentID, "")
	fileList, err := authorize.FileList(option)
	if err != nil {
		log.Printf("获取文件列表情况失败: %s\n", err)
		return
	}

	log.Printf("获取文件列表成功, 文件数量: %d\n", len(fileList.Items))
}

// GetFileInfo 获取文件信息, 目录和文件都支持
func GetFileInfo(authorize aliyundrive_open.Authorize, fileID string) (file aliyundrive_open.FileInfo, err error) {
	option := aliyundrive_open.NewFileOption(authorize.DriveID, fileID)
	file, err = authorize.File(option)
	if err != nil {
		log.Printf("获取文件信息失败: %s\n", err)
		return
	}

	log.Printf("文件ID: %s 名称: %s, 类型: %s 大小: %d\n", file.FileId, file.Name, file.Type, file.Size)
	return file, err
}

// GetFilesInfo 批量获取文件信息
// 这里我们使用了 NewFilesOption 方法来生成请求选项. 其余个性化参数可以通过 Set 方法设置
func GetFilesInfo(authorize aliyundrive_open.Authorize, ids []string) {

	bOption := aliyundrive_open.NewFilesOption(authorize.DriveID, ids)
	bFiles, err := authorize.Files(bOption)
	if err != nil {
		log.Printf("批量获取文件信息失败: %s\n", err)
		return
	}

	log.Printf("批量获取文件信息成功, 文件数量: %d\n", len(bFiles.Items))

	for _, file := range bFiles.Items {
		log.Printf("文件ID: %s 名称: %s, 类型: %s 大小: %d\n", file.FileId, file.Name, file.Type, file.Size)
	}

}

// 获取文件下载地址
func GetDownloadURL(authorize aliyundrive_open.Authorize, fileID string) {
	option := aliyundrive_open.NewFileDownloadURLOption(authorize.DriveID, fileID)
	downInfo, err := authorize.FileDownloadURL(option)
	if err != nil {
		log.Printf("获取文件下载地址失败: %s\n", err)
		return
	}

	log.Println("文件下载地址: ", downInfo.URL)
}

// RenameFile 重命名文件
func RenameFile(authorize aliyundrive_open.Authorize, fileID string, newName string) {
	option := aliyundrive_open.NewFileRenameOption(authorize.DriveID, fileID, newName)
	result, err := authorize.FileRename(option)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("修改名字成功: %s\n", result.Name)
}

// GetVideoPlayURL 获取视频播放地址
func GetVideoPlayInfo(authorize aliyundrive_open.Authorize, fileID string) {
	option := aliyundrive_open.NewFileVideoPlayInfoOption(authorize.DriveID, fileID)
	result, err := authorize.FileVideoPlayInfo(option)
	if err != nil {
		log.Println(err)
		return
	}

	for _, playInfo := range result.VideoPreviewPlayInfo.LiveTranscodingTaskList {
		log.Printf("视频播放地址: %s\n", playInfo.Url)
	}
}

// MoveFile 移动文件
func MoveFile(authorize aliyundrive_open.Authorize, fileID, parentID string) {

	file, err := GetFileInfo(authorize, fileID)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("移动前文件 %s(%s) 父目录ID: %s\n", file.Name, file.FileId, file.ParentFileId)

	option := aliyundrive_open.NewFileMoveAndCopyOption(authorize.DriveID, fileID, parentID)
	_, err = authorize.FileMove(option)
	if err != nil {
		log.Println(err)
		return
	}

	file, err = GetFileInfo(authorize, fileID)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("移动后文件 %s(%s) 父目录ID: %s\n", file.Name, file.FileId, file.ParentFileId)
}

// CopyFile 复制文件
func CopyFile(authorize aliyundrive_open.Authorize, fileID, toParentID string) {
	option := aliyundrive_open.NewFileMoveAndCopyOption(authorize.DriveID, fileID, toParentID)
	_, err := authorize.FileCopy(option)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("复制文件成功: %s\n", fileID)
}

// CreateFolder 创建目录
func CreateFolder(authorize aliyundrive_open.Authorize, parentID, folderName string) (result aliyundrive_open.FileCreate, err error) {
	option := aliyundrive_open.NewFileCreateOption(authorize.DriveID, "root", "新目录")
	result, err = authorize.FolderCreate(option)
	if err != nil {
		log.Printf("Token 刷新失败: %s\n", err)
	}
	return result, err
}

// TrashFile 将文件移入回收站
func TrashFile(authorize aliyundrive_open.Authorize, fileID string) (result aliyundrive_open.FileMoveCopyDelTask, err error) {
	option := aliyundrive_open.NewFileTrashAndDeleteOption(authorize.DriveID, fileID)
	result, err = authorize.FileTrash(option)
	if err != nil {
		log.Printf("Token 刷新失败: %s\n", err)
	}
	return result, err
}

// 彻底删除文件
func DeleteFile(authorize aliyundrive_open.Authorize, fileID string) (result aliyundrive_open.FileMoveCopyDelTask, err error) {
	option := aliyundrive_open.NewFileTrashAndDeleteOption(authorize.DriveID, fileID)
	result, err = authorize.FileDelete(option)
	if err != nil {
		log.Printf("Token 刷新失败: %s\n", err)
	}
	return result, err
}

// UploadFile 上传文件
func UploadFile(authorize aliyundrive_open.Authorize, filePath string) (uploadResult aliyundrive_open.FileInfo, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("打开文件失败: %s\n", err)
		return uploadResult, err
	}

	_, name := filepath.Split(file.Name())

	// 上传文件
	option := aliyundrive_open.NewFileUploadOption(authorize.DriveID, "root", name, file)
	option.SetParallelUpload(true)
	uploadResult, err = authorize.FileUpload(option)
	if err != nil {
		log.Printf("上传文件失败: %s\n", err)
	}

	return uploadResult, err
}

func main() {
	LoginQRCode()
}
