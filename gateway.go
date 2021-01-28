package kaiheila

type gatewayResp struct {
	URL string
}

// GetGateway get websocket url from gateway
func (c *Client) GetGateway() (string, error) {
	gateway := gatewayResp{}
	err := c.request("GET", 3, "gateway/index", nil, &gateway)
	return gateway.URL, err
}
