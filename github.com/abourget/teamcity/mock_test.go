package teamcity

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func NewTestClient(replyResp *http.Response, err error) *Client {
	client := &Client{
		username: "username",
		password: "password",
		host:     "host.example.com",
	}
	httpClient := &http.Client{}
	httpClient.Transport = &MockTransport{
		resp: replyResp,
		err:  err,
	}
	client.HTTPClient = httpClient
	return client
}

type MockTransport struct {
	req  *http.Request
	resp *http.Response
	err  error
}

func (b *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	b.req = req
	return b.resp, b.err
}

func newResponse(body string) *http.Response {
	resp := &http.Response{Body: ioutil.NopCloser(bytes.NewBuffer([]byte(body)))}
	return resp
}
