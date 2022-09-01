package idrac

import (
	"encoding/xml"
	"fmt"
	"io"
)

const endpointLogout = "/cgi-bin/logout"

type LogoutResponse struct {
	XMLName  xml.Name `xml:"LOGOUT"`
	Response struct {
		XMLName    xml.Name `xml:"RESP"`
		ReturnCode string   `xml:"RC"`
		SessionID  string   `xml:"SID"`
	}
}

// login logs out of the idrac
func (rac *idrac) Logout() (logoutResp LogoutResponse, err error) {
	// GET (not post) logout
	resp, err := rac.client.Get(rac.url() + endpointLogout)
	if err != nil {
		return LogoutResponse{}, err
	}

	// read and unmarshal body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LogoutResponse{}, err
	}

	err = xml.Unmarshal(body, &logoutResp)
	if err != nil {
		return LogoutResponse{}, err
	}

	// verify logout success
	if logoutResp.Response.ReturnCode != "0x0" {
		return LogoutResponse{}, fmt.Errorf("logout failed (code: %s)", logoutResp.Response.ReturnCode)
	}

	return logoutResp, nil
}
