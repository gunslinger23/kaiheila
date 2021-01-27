package kaiheila

type gateway struct {
	Url string
}

func (c *Client) GetGateway() (string, error) {
	gateway := gateway{}
	err := c.request("GET", 3, "gateway/index", nil, &gateway)
	return gateway.Url, err
}
