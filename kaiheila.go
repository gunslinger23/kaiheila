package kaiheila

const (
	apiURLDefault  = "https://www.kaiheila.cn/api"
	timeoutDefault = 1
	// TokenBot Type: bot
	TokenBot = "Bot"
	// TokenOAuth2 Type: OAuth2
	TokenOAuth2 = "Bearer"
)

// Client sdk client
type Client struct {
	URL       string
	Token     string
	TokenType string
	// API request timeout (sec)
	Timeout int
}

// NewClient create a new client for access
func NewClient(url, tokenType, token string, timeout int) *Client {
	if len(url) == 0 {
		url = apiURLDefault
	}
	if timeout == 0 {
		timeout = timeoutDefault
	}

	return &Client{
		URL:       url,
		Token:     token,
		TokenType: tokenType,
		Timeout:   timeout,
	}
}
