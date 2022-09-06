package idrac

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/fcjr/aia-transport-go"
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
	// make default timeouts
	tlsTimeout := 5 * time.Second
	clientTimeout := 10 * time.Second

	// make transport based ignoreCertErrors
	transport := &http.Transport{}
	if strictCerts {
		// make transport with AIA support (many (all?) idracs do not
		// supply intermediate cert and linux does not perform AIA by
		// default)
		transport, err = aia.NewTransport()
		if err != nil {
			return nil, err
		}
	} else {
		// not strict
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	// configure transport
	transport.TLSHandshakeTimeout = tlsTimeout
	// default connections per host
	transport.MaxConnsPerHost = 2
	transport.MaxIdleConnsPerHost = 2

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
