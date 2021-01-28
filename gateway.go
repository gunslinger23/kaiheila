package kaiheila

type gatewayResp struct {
	URL string
}

// GetGateway Get websocket url from gateway
func (c *Client) GetGateway() (url string, err error) {
	gateway := gatewayResp{}
	err = c.request("GET", 3, "gateway/index", nil, &gateway)
	return gateway.URL, err
}
