package http

func newRequest(client *Client) *Request {
	return &Request{Request: client.R()}
}
