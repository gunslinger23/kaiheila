package kaiheila

const (
	apiUrlDefault  = "https://www.kaiheila.cn/api"
	timeoutDefault = 1
	TokenBot       = "Bot"
	TokenOauth2    = "Bearer"
)

type Client struct {
	Url       string
	Token     string
	TokenType string
	// API request timeout (sec)
	Timeout int
	// Proxy
	HttpProxy string
}

type APIRequest struct {
	Version int
	Method  string
	Path    string
}

func NewClient(url, tokenType, token string, timeout int) *Client {
	if len(url) == 0 {
		url = apiUrlDefault
	}
	if timeout == 0 {
		timeout = timeoutDefault
	}

	return &Client{
		Url:       url,
		Token:     token,
		TokenType: tokenType,
		Timeout:   timeout,
	}
}
