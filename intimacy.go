package kaiheila

// GetIntimacyReq Get intimacy request struct
type GetIntimacyReq struct {
	UserID string `json:"user_id"` // User's ID
}

// GetIntimacyResp Get intimacy respone struct
type GetIntimacyResp struct {
	ImgURL     string      `json:"img_url"`     // URL of image
	SocialInfo string      `json:"social_info"` // Social Info
	LastRead   int64       `json:"last_read"`   // Timestamp of last read
	Score      string      `json:"score"`       // Score of intimacy(0-2200)
	ImgList    []ImageList `json:"img_list"`    // List of images
}

// ImageList List of images
type ImageList struct {
	ID  string `json:"id"`  // Image ID
	URL string `json:"url"` // Image URL
}

// GetIntimacy Get user's intimacy
func (c *Client) GetIntimacy(req GetIntimacyReq) (GetIntimacyResp, error) {
	resp := GetIntimacyResp{}
	err := c.request("GET", 3, "intimacy/index", &req, &resp)
	return resp, err
}

// UpdateIntimacyReq Update intimacy request struct
type UpdateIntimacyReq struct {
	UserID     string `json:"user_id"`               // User's ID
	Score      int    `json:"score,omitempty"`       // (Optional) Score of intimacy(0-2200)
	SocialInfo string `json:"social_info,omitempty"` // (Optional) Social Info
	ImgID      int    `json:"img_id,omitempty"`      // (Optional) Image ID
}

// UpdateIntimacy Update user's intimacy
func (c *Client) UpdateIntimacy(req UpdateIntimacyReq) error {
	err := c.request("POST", 3, "intimacy/update", &req, nil)
	return err
}
