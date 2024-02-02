package idrac

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
)

// idracClient is a custom http.Client designed to interface
// with idrac in a way similar to racadm
type idracClient struct {
	http      http.Client
	userAgent string
}

// New creates a new http client using the custom struct. It also
// specifies timeout options so the client behaves sanely.
func newIdracClient(strictCerts bool) (client *idracClient, err error) {
	// make timeouts
	clientTimeout := 60 * time.Second

	// make transport based on strictCerts
	transport, err := newIdracAiaTransport(strictCerts)
	if err != nil {
		return nil, err
	}

	// create *Client
	client = new(idracClient)
	client.http.Timeout = clientTimeout
	client.http.Transport = transport

	// cookie jar to save login (sid) cookie
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client.http.Jar = jar

	// based on racadm 9.1.2
	client.userAgent = "SSLClient"

	return client, nil
}

// newRequest creates an http request for the client to later do
func (client *idracClient) newRequest(method string, url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// based on racadm 9.1.2
	request.Header.Set("Cache-Control", "no-cache")

	// set user agent
	request.Header.Set("User-Agent", client.userAgent)

	return request, nil
}

// do does the specified request
func (client *idracClient) do(request *http.Request) (*http.Response, error) {
	response, err := client.http.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Get does a get request to the specified url
func (client *idracClient) Get(url string) (*http.Response, error) {
	request, err := client.newRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return client.do(request)
}

// // Head does a head request to the specified url
// // a head request is the same as Get but without the body
// func (client *Client) Head(url string) (*http.Response, error) {
// 	request, err := client.newRequest(http.MethodHead, url, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return client.do(request)
// }

// Post does a post request using the specified url, content type, and
// body
func (client *idracClient) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	request, err := client.newRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", contentType)

	return client.do(request)
}
