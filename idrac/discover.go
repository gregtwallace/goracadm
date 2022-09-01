package idrac

import (
	"encoding/xml"
	"errors"
	"io"
	"strings"
)

const endpointDiscover = "/cgi-bin/discover"

var errBadResp = errors.New("discovery response unmarshalled but is incorrect")

type DiscoverResponse struct {
	XMLName  xml.Name `xml:"DISCOVER"`
	Response struct {
		XMLName         xml.Name `xml:"RESP"`
		ReturnCode      string   `xml:"RC"`
		EndpointType    string   `xml:"ENDPOINTTYPE"`
		EndpointVersion string   `xml:"ENDPOINTVER"`
		ProtocolType    string   `xml:"PROTOCOLTYPE"`
		ProtocolVersion string   `xml:"PROTOCOLVER"`
	}
}

// Discover inquires for basic information from the idrac
func (rac *idrac) Discover() (discResp DiscoverResponse, err error) {
	// do discover
	resp, err := rac.client.Get(rac.url() + endpointDiscover)
	if err != nil {
		return DiscoverResponse{}, err
	}

	// read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return DiscoverResponse{}, err
	}

	// unmarshal body
	err = xml.Unmarshal(body, &discResp)
	if err != nil {
		return DiscoverResponse{}, err
	}

	// verify response code is good, endpoint type is correct, and
	// https
	if discResp.Response.ReturnCode != "0x0" ||
		!strings.Contains(strings.ToLower(discResp.Response.EndpointType), "idrac") ||
		strings.ToLower(discResp.Response.ProtocolType) != "https" {

		return DiscoverResponse{}, errBadResp
	}

	// good response
	return discResp, nil
}
