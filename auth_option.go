package aliyundrive_open

// H5页多种登录方式选项
type AuthorizeOption struct {
	ClientID     string  `json:"client_id"`              // 开放平台应用ID
	ClientSecret string  `json:"client_secret"`          // 开放平台应用密钥
	Scopes       []Scope `json:"scopes"`                 // 授权范围
	RedirectUri  string  `json:"redirect_uri,omitempty"` // 回调地址
	State        string  `json:"state,omitempty"`        // 防止CSRF攻击
}

// NewDefaultMultipleAuthorizeOption 创建默认授权选项
// 网页多种登录方式, 回调地址为必传,否则无法接收code.如无法满足条件,请使用使用 "单一扫码方式"
func NewDefaultMultipleAuthorizeOption(redirectUri string) *AuthorizeOption {
	return &AuthorizeOption{
		Scopes: []Scope{
			ScopeBase,
			ScopePhone,
			ScopeRead,
			ScopeWrite,
		},
		State:       randomString(8),
		RedirectUri: redirectUri,
	}
}

// NewDefaultSingleAuthorizeOption 创建默认单一扫码方式授权选项
func NewDefaultSingleAuthorizeOption() *AuthorizeOption {
	return &AuthorizeOption{
		Scopes: []Scope{
			ScopeBase,
			ScopePhone,
			ScopeRead,
			ScopeWrite,
		},
	}
}

// NewMultipleAuthorizeOption 创建H5页多登陆方式授权选项
func NewMultipleAuthorizeOption(redirectUri string) *AuthorizeOption {
	return &AuthorizeOption{
		RedirectUri: redirectUri,
	}
}

// NewSingleAuthorizeOption 创建单一扫码方式授权选项
func NewSingleAuthorizeOption() *AuthorizeOption {
	return &AuthorizeOption{}
}

// SetScopes 设置授权范围
func (option *AuthorizeOption) SetScopes(scopes []Scope) *AuthorizeOption {
	option.Scopes = scopes
	return option
}

// SetState 设置防止CSRF攻击
func (option *AuthorizeOption) SetState(state string) *AuthorizeOption {
	option.State = state
	return option
}
