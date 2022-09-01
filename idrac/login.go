package idrac

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
)

const endpointLogin = "/cgi-bin/login"

type loginPayload struct {
	XMLName xml.Name `xml:"LOGIN"`
	Request struct {
		XMLName  xml.Name `xml:"REQ"`
		Username string   `xml:"USERNAME"`
		Password string   `xml:"PASSWORD"`
		// CommandInput string `xml:"CMDINPUT"`
	}
}

type LoginResponse struct {
	XMLName  xml.Name `xml:"LOGIN"`
	Response struct {
		XMLName           xml.Name   `xml:"RESP"`
		ReturnCode        ReturnCode `xml:"RC"`
		SessionID         string     `xml:"SID"`
		State             string     `xml:"STATE"`
		StateName         string     `xml:"STATENAME"`
		DefaultCredential string     `xml:"DEFCRED"`
	}
}

// login logs into the idrac
func (rac *idrac) Login() (loginResp LoginResponse, err error) {
	// make login payload and marshal it
	payload := loginPayload{}
	payload.Request.Username = rac.username
	payload.Request.Password = rac.password

	payloadXml, err := xml.Marshal(payload)
	if err != nil {
		return LoginResponse{}, err
	}

	// post the login
	resp, err := rac.client.Post(rac.url()+endpointLogin, "application/xml", bytes.NewBuffer(payloadXml))
	if err != nil {
		return LoginResponse{}, err
	}

	// read and unmarshal body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LoginResponse{}, err
	}

	err = xml.Unmarshal(body, &loginResp)
	if err != nil {
		return LoginResponse{}, err
	}

	// verify login success
	if loginResp.Response.ReturnCode != RcOK {
		return LoginResponse{}, loginResp.Response.ReturnCode
	}

	// save login cookie to jar
	url, err := url.Parse("https://" + rac.hostname)
	if err != nil {
		return LoginResponse{}, err
	}
	loginCookie := &http.Cookie{
		Name:  "sid",
		Value: loginResp.Response.SessionID,
		Path:  "/cgi-bin/",
	}
	rac.client.http.Jar.SetCookies(url, []*http.Cookie{loginCookie})

	return loginResp, nil
}
