package idrac

import (
	"errors"
)

// idrac contains details about a specific idrac
type idrac struct {
	hostname string
	username string
	password string
	client   *idracClient
}

// NewIdrac creates an Idrac and client to access it
func NewIdrac(hostname, username, password string, strictCerts bool) (*idrac, error) {
	if hostname == "" {
		return nil, errors.New("hostname (-r) must be specified")
	} else if username == "" {
		return nil, errors.New("username (-u) must be specified")
	} else if password == "" {
		return nil, errors.New("password (-p) must be specified")
	}

	// make http client for idrac
	idracClient, err := newIdracClient(strictCerts)
	if err != nil {
		return nil, err
	}

	return &idrac{
		hostname: hostname,
		username: username,
		password: password,
		client:   idracClient,
	}, nil
}

// url returns the base url to access the idrac
func (rac *idrac) url() string {
	return "https://" + rac.hostname
}
