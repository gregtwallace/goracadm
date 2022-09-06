package idrac

import "log"

// idrac contains details about a specific idrac
type idrac struct {
	hostname string
	username string
	password string
	client   *idracClient
}

// NewIdrac creates an Idrac and client to access it
func NewIdrac(hostname, username, password string, strictCerts bool) *idrac {
	// make http client for idrac
	idracClient, err := newIdracClient(strictCerts)
	if err != nil {
		log.Fatal(err)
	}

	return &idrac{
		hostname: hostname,
		username: username,
		password: password,
		client:   idracClient,
	}
}

// url returns the base url to access the idrac
func (rac *idrac) url() string {
	return "https://" + rac.hostname
}
