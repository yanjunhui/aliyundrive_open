package aliyundrive_open

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

const (
	APIBase = "https://openapi.aliyundrive.com"

	//用户权限相关
	APIAuthorizeMultiple     = APIBase + "/oauth/authorize"               //多种授权方式
	APIAuthorizeQrCode       = APIBase + "/oauth/authorize/qrcode"        //获取二维码, 仅支持扫码登录
	APIAuthorizeQrCodeStatus = APIBase + "/oauth/qrcode/%s/status"        //获取扫码结果
	APIRefreshToken          = APIBase + "/oauth/access_token"            //刷新 access_token
	APIDriveInfo             = APIBase + "/adrive/v1.0/user/getDriveInfo" //获取用户云盘信息
	APISpaceInfo             = APIBase + "/adrive/v1.0/user/getSpaceInfo" //获取空间大小信息

	//文件操作相关
	APIList              = APIBase + "/adrive/v1.0/openFile/list"                    //获取文件列表
	APIFile              = APIBase + "/adrive/v1.0/openFile/get"                     //获取文件信息
	APIFiles             = APIBase + "/adrive/v1.0/openFile/batch/get"               //批量获取文件信息
	APIFileTrash         = APIBase + "/adrive/v1.0/openFile/recyclebin/trash"        //移动文件到垃圾箱
	APIFileDelete        = APIBase + "/adrive/v1.0/openFile/delete"                  //彻底删除文件
	APIFileCreate        = APIBase + "/adrive/v1.0/openFile/create"                  //创建目录/文件
	APIFileComplete      = APIBase + "/adrive/v1.0/openFile/complete"                //创建文件完成
	APIFileDownload      = APIBase + "/adrive/v1.0/openFile/getDownloadUrl"          //获取下载链接
	APIFileVideoPlayInfo = APIBase + "/adrive/v1.0/openFile/getVideoPreviewPlayInfo" //获取视频转码播放信息
	APIFileMove          = APIBase + "/adrive/v1.0/openFile/move"                    //移动文件
	APIFileCopy          = APIBase + "/adrive/v1.0/openFile/copy"                    //复制文件
	APIFileUpdate        = APIBase + "/adrive/v1.0/openFile/update"                  //更新文件

)

var RestyHttpClient = NewRestyClient()
var UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36 Edg/108.0.1462.54"
var DefaultTimeout = time.Second * 30

func NewRestyClient() *resty.Client {
	return resty.New().
		SetHeader("user-agent", UserAgent).
		SetRetryCount(3).
		SetTimeout(DefaultTimeout)
}

// HttpPost 请求
func (a *Authorize) HttpPost(url string, reqData interface{}, result interface{}) error {
	header := http.Header{}
	header.Set("Authorization", "Bearer "+a.AccessToken)
	return HttpPost(url, header, reqData, result)
}

func HttpPost(url string, header http.Header, reqData interface{}, result interface{}) error {
	r := RestyHttpClient.R()
	if reqData != nil {
		dataJson, err := json.Marshal(reqData)
		if err != nil {
			return err
		}
		r.SetBody(dataJson)
	}

	if header == nil {
		header = http.Header{}
	}

	header.Set("Content-Type", "application/json;charset=UTF-8")
	resp, err := r.SetHeaderMultiValues(header).Post(url)
	if err != nil {
		return errors.New("请求失败: " + err.Error())
	}

	err = json.Unmarshal(resp.Body(), result)
	if err != nil {
		return errors.New("解析数据失败: " + err.Error())
	}
	return err
}
