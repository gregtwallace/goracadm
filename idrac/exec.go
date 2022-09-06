package idrac

import (
	"bytes"
	"encoding/xml"
	"errors"
	"flag"
	"io"
	"log"
)

const endpointExec = "/cgi-bin/exec"

var (
	errInvalidSubCommand      = errors.New("subcommand is either invalid or not implemented")
	errInvalidOrMalpositioned = errors.New("invalid or malpositioned param or flag")
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
	// subcommands:
	// https://www.dell.com/support/manuals/en-us/poweredge-m630/idrac8_2.70.70.70_racadm/racadm-subcommand-details?guid=guid-cd4e81e6-818c-44fb-9e7a-82950425fbbb&lang=en-us
	// https://www.dell.com/support/manuals/en-us/idrac9-lifecycle-controller-v5.x-series/idrac9_5.xx_racadm_pub/racadm-subcommand-details?guid=guid-3e09aba8-6e2c-4fd9-9a17-d05f2596dbac&lang=en-us
	switch command {
	case "racreset":
		execResp, err = rac.racreset(flags)
	case "sslcertdownload":
		execResp, err = rac.sslcertdownload(flags)
	case "sslkeyupload":
		execResp, err = rac.sslkeyupload(flags)
	default:
		// error, unsupported
		return execResponse{}, errInvalidSubCommand
	}

	// debugging
	// log.Println(execResp)

	return execResp, err
}

// executePayload executes the specified payload against
// the idrac and returns the response or an error.
func (rac *idrac) executePayload(payload execPayload) (execResp execResponse, err error) {
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

	// check return codes for errors
	if execResp.Response.ReturnCode != RcOK {
		return execResponse{}, execResp.Response.ReturnCode
	}
	if execResp.Response.CommandReturnCode != RcOK {
		return execResponse{}, execResp.Response.CommandReturnCode
	}

	// success - write command output
	log.Println(execResp.Response.CommandOutput)

	return execResp, nil
}

// parseFlags parses the flag set and returns an error if there are
// any extraneous / leftover bits after the flags are parsed.
func parseFlags(fs *flag.FlagSet, flags []string) (err error) {
	fs.Parse(flags)

	// check for leftovers
	if len(fs.Args()) > 0 {
		return errInvalidOrMalpositioned
	}

	return nil
}
