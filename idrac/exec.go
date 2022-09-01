package idrac

import (
	"bytes"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

const endpointExec = "/cgi-bin/exec"

var (
	errInvalidSubCommand = errors.New("subcommand is either invalid or not implemented")
)

type execPayload struct {
	XMLName xml.Name `xml:"EXEC"`
	Request struct {
		XMLName       xml.Name `xml:"REQ"`
		CommandInput  string   `xml:"CMDINPUT"`
		MaxOutputLen  string   `xml:"MAXOUTPUTLEN"`
		Capability    string   `xml:"CAPABILITY"`
		UserPrivilege int      `xml:"USERPRIV"`
	}
}

type execResponse struct {
	XMLName  xml.Name `xml:"EXEC"`
	Response struct {
		XMLName           xml.Name `xml:"RESP"`
		ReturnCode        string   `xml:"RC"`
		OutputLen         string   `xml:"OUTPUTLEN"`
		CommandReturnCode string   `xml:"CMDRC"`
		Capability        string   `xml:"CAPABILITY"`
		CommandOutput     string   `xml:"CMDOUTPUT"`
	}
}

// TODO: https://www.digitalocean.com/community/tutorials/how-to-use-the-flag-package-in-go
type SubCommand interface {
}

// Exec executes the specified command against the idrac. To
// avoid unexpected behavior, error if specified command has
// not been specifically implemented and tested.
func (rac *idrac) Exec() (execResp execResponse, err error) {
	// check subcommand is implemented and parse flags accordingly
	// subcommands: https://www.dell.com/support/manuals/en-us/idrac9-lifecycle-controller-v5.x-series/idrac9_5.xx_racadm_pub/racadm-subcommand-details?guid=guid-3e09aba8-6e2c-4fd9-9a17-d05f2596dbac&lang=en-us
	switch flag.Args()[0] {
	case "sslcertdownload":
		return rac.sslcertdownload()
	default:
		// break for error
	}

	return execResponse{}, errInvalidSubCommand
}

// sslcertdownload executes the sslcertdownload subcommand using
// the specified flags.
// https://www.dell.com/support/manuals/en-us/oth-r6415/idrac9_5.xx_racadm_pub/sslcertdownload?guid=guid-33c6a0ac-ee43-4bb6-9413-1e83e359144a&lang=en-us
func (rac *idrac) sslcertdownload() (execResp execResponse, err error) {
	// drop the subcommand to leave the flags
	flags := flag.Args()[1:]

	// parse flags for sslcertdownload command
	filename := ""
	certType := 0
	instance := 0

	fs := flag.NewFlagSet("sslcertdownload", flag.ExitOnError)
	fs.StringVar(&filename, "f", "", "filename to save the cert locally (required)")
	fs.IntVar(&certType, "t", 0, "certificate type (required - int - see Dell docs)")
	fs.IntVar(&instance, "i", 0, "instance (1 or 2) (optional)")

	fs.Parse(flags)

	// validate command flags
	if filename == "" {
		return execResponse{}, errors.New("sslcertdownload: filename (f) must be specified")
	}
	if certType == 0 {
		return execResponse{}, errors.New("sslcertdownload: cert type (t) must be specified")
	}
	if certType < 1 || certType > 11 {
		return execResponse{}, errors.New("sslcertdownload: cert type must be between 1 and 11, inclusive")
	}

	// confirm file doesn't already exist
	_, err = os.Stat(filename)
	if err == nil {
		return execResponse{}, errors.New("filename already exists")
	} else if !errors.Is(err, os.ErrNotExist) {
		return execResponse{}, err
	}

	// optional, validate and make param if appropriate
	instanceParam := ""
	if instance == 0 {
		// no-op
	} else if instance == 1 || instance == 2 {
		// add -i param
		instanceParam = fmt.Sprintf(" -i %d", instance)
	} else {
		// error, invalid -i
		return execResponse{}, errors.New("sslcertdownload: instance must be 1 or 2, if specified")
	}

	// build payload to post to drac
	cmdInput := fmt.Sprintf("racadm sslcertdownload -f sslcertfile -t %d%s", certType, instanceParam)

	payload := execPayload{}
	payload.Request.CommandInput = cmdInput
	payload.Request.MaxOutputLen = "0x0fff"
	payload.Request.Capability = "0x1"
	payload.Request.UserPrivilege = 0

	// marshal payload
	payloadXml, err := xml.Marshal(payload)
	if err != nil {
		return execResponse{}, err
	}

	// post
	resp, err := rac.client.Post(rac.url()+endpointExec, "application/xml", bytes.NewBuffer(payloadXml))
	if err != nil {
		return execResponse{}, err
	}

	// read and unmarshal body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return execResponse{}, err
	}

	err = xml.Unmarshal(body, &execResp)
	if err != nil {
		return execResponse{}, err
	}

	// save certificate to specified file
	f, err := os.Create(filename)
	if err != nil {
		return execResponse{}, err
	}

	// write cert to file
	_, err = f.WriteString(execResp.Response.CommandOutput)
	if err != nil {
		return execResponse{}, err
	}

	return execResp, nil
}
