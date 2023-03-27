package aliyundrive_open

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Scope AuthorizeOption 授权类型
type Scope string

const (
	ScopeBase  Scope = "user:base"                                         // 授权获取用户ID, 头像, 昵称
	ScopePhone Scope = "user:phone"                                        // 获取手机号
	ScopeRead  Scope = "file:all:read"                                     // 所有文件读取权限
	ScopeWrite Scope = "file:all:write"                                    // 所有文件写入权限
	ScopeAll   Scope = "file:all:read,file:all:write,user:base,user:phone" // 所有权限
)

func (s Scope) String() string {
	return string(s)
}

// AuthorizeURL 构建 H5前端 授权页面. 需要一个回调地址接收 code
// 拼接示例 https://openapi.aliyundrive.com/oauth/authorize?client_id=xxx&redirect_uri=xxx&scope=user:base,user:phone,file:all:read,file:all:write&state=xxx
func (c *Client) AuthorizeURL(option *AuthorizeOption) (authURL string, err error) {
	if option == nil {
		err = fmt.Errorf("option is nil")
		return
	}

	values := make(url.Values)
	values.Set("client_id", c.ClientId)
	values.Set("redirect_uri", option.RedirectUri)
	values.Set("scope", joinCustomString(option.Scopes, ","))

	u, _ := url.Parse(APIAuthorizeMultiple)
	u.RawQuery = values.Encode()

	return u.String(), nil
}

// ReceiveAuthorizeCode 接收前端授权 code, 并获得授权
func (c *Client) ReceiveAuthorizeCode(req *http.Request) (result Authorize, err error) {
	queryParams := req.URL.Query()

	code := queryParams.Get("code")
	if code == "" {
		err = fmt.Errorf("code 为空")
		return result, err
	}
	return c.Authorize(code)
}

// AuthorizeQRCode 授权二维码数据
type AuthorizeQRCode struct {
	QrCodeUrl string `json:"qrCodeUrl"`
	Sid       string `json:"sid"`
	ErrorInfo
}

// QRCode  获取登录二维码信息
func (c *Client) QRCode(option *AuthorizeOption) (result AuthorizeQRCode, err error) {
	req := map[string]interface{}{
		"client_id":     c.ClientId,
		"client_secret": c.ClientSecret,
		"scopes":        option.Scopes,
	}

	err = HttpPost(APIAuthorizeQrCode, nil, req, &result)
	if err != nil {
		return result, err
	}

	if result.Code != "" {
		err = fmt.Errorf("获取二维码失败: %s", result.Message)
	}

	return result, err
}

type AuthorizeQRCodeStatus struct {
	Status   string `json:"status"` //状态有三种: waiting, success, failed
	AuthCode string `json:"authCode"`
	ErrorInfo
}

// QrCodeStatus 获取二维码状态
func (c *Client) QrCodeStatus(sid string) (result AuthorizeQRCodeStatus, err error) {
	if sid == "" {
		err = fmt.Errorf("需要传入 QRCode 方法返回 sid 值")
		return result, err
	}

	_, err = RestyHttpClient.R().SetResult(&result).Get(fmt.Sprintf(APIAuthorizeQrCodeStatus, sid))
	if err != nil {
		return result, err
	}

	if result.Code != "" {
		err = fmt.Errorf("获取二维码状态失败: %s", result.Message)
	}

	return result, err
}

// Authorize 登录授权信息
type Authorize struct {
	TokenType    string    `json:"token_type"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"`
	ExpiresTime  time.Time `json:"expires_time"`
	DriveID      string    `json:"drive_id"`
	ErrorInfo
}

//

// Authorize 授权登录
func (c *Client) Authorize(authCode string) (result Authorize, err error) {
	if authCode == "" {
		err = fmt.Errorf("需要传入 QrCodeStatus 方法返回 authCode 值")
		return result, err
	}

	req := map[string]string{
		"client_id":     c.ClientId,
		"client_secret": c.ClientSecret,
		"code":          authCode,
		"grant_type":    "authorization_code",
	}

	err = HttpPost(APIRefreshToken, nil, req, &result)
	if err != nil {
		return result, err
	}

	if result.Code != "" {
		err = fmt.Errorf("授权失败: %s", result.Message)
	}

	result.ExpiresTime = time.Now().Add(time.Duration(result.ExpiresIn-60) * time.Second)

	info, err := result.DriveInfo()
	if err != nil {
		return result, err
	}

	result.DriveID = info.DefaultDriveId
	c.DriveID = info.DefaultDriveId

	return result, err
}

// RefreshToken 刷新 token
func (c *Client) RefreshToken(refreshToken string) (result Authorize, err error) {
	req := map[string]string{
		"client_id":     c.ClientId,
		"client_secret": c.ClientSecret,
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	}

	err = HttpPost(APIRefreshToken, nil, req, &result)
	if err != nil {
		return result, err
	}

	if result.Code != "" {
		err = fmt.Errorf("刷新授权失败: %s", result.Message)
	}

	result.ExpiresTime = time.Now().Add(time.Duration(result.ExpiresIn-60) * time.Second)
	result.DriveID = c.DriveID
	return result, err
}
