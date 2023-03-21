package aliyundrive_open

import "log"

type Client struct {
	ClientId     string //开放平台应用ID
	ClientSecret string //开放平台应用密钥
	DriveID      string //阿里云盘ID
}

type ErrorInfo struct {
	Code      string `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	RequestId string `json:"requestId,omitempty"`
}

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

func NewClient(clientID, secret string) *Client {
	return &Client{
		ClientId:     clientID,
		ClientSecret: secret,
	}
}
