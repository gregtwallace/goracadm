package idrac

import (
	"encoding/xml"
	"errors"
)

const endpointExec = "/cgi-bin/exec"

var (
	errInvalidSubCommand = errors.New("subcommand is either invalid or not implemented")
)

// execPayload is the payload to execute on idrac
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

// execResponse is the idrac's response to an execution
type execResponse struct {
	XMLName  xml.Name `xml:"EXEC"`
	Response struct {
		XMLName           xml.Name   `xml:"RESP"`
		ReturnCode        ReturnCode `xml:"RC"`
		OutputLen         string     `xml:"OUTPUTLEN"`
		CommandReturnCode ReturnCode `xml:"CMDRC"`
		Capability        string     `xml:"CAPABILITY"`
		CommandOutput     string     `xml:"CMDOUTPUT"`
	}
}

// Exec executes the specified command against the idrac. To
// avoid unexpected behavior, error if specified command has
// not been specifically implemented and tested.
func (rac *idrac) Exec(command string, flags []string) (execResp execResponse, err error) {
	// check subcommand is implemented and parse flags accordingly
	// subcommands: https://www.dell.com/support/manuals/en-us/idrac9-lifecycle-controller-v5.x-series/idrac9_5.xx_racadm_pub/racadm-subcommand-details?guid=guid-3e09aba8-6e2c-4fd9-9a17-d05f2596dbac&lang=en-us
	switch command {
	case "sslcertdownload":
		execResp, err = rac.sslcertdownload(flags)
	default:
		// error, unsupported
		return execResponse{}, errInvalidSubCommand
	}

	// debugging
	// log.Println(execResp)

	return execResp, err
}
