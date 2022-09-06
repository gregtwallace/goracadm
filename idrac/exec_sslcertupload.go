package idrac

import (
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

// sslkeyupload executes the sslkeyupload subcommand using
// the specified flags.
// https://www.dell.com/support/manuals/en-us/poweredge-m630/idrac8_2.70.70.70_racadm/sslcertupload?guid=guid-c1610ee7-2216-4f05-904c-50ae536e8412&lang=en-us
// https://www.dell.com/support/manuals/en-us/idrac9-lifecycle-controller-v5.x-series/idrac9_5.xx_racadm_pub/sslcertupload?guid=guid-4c93d9c0-ec1f-42a3-b746-67d980819ba7&lang=en-us
func (rac *idrac) sslcertupload(flags []string) (execResp execResponse, err error) {
	// parse command flags (options)
	filename := ""
	certType := 0

	fs := flag.NewFlagSet("sslcertupload", flag.ExitOnError)
	fs.StringVar(&filename, "f", "", "local filename to upload (required)")
	fs.IntVar(&certType, "t", 0, "certificate type (required - int - see Dell docs)")

	// parse and check for basic errors
	err = parseFlags(fs, flags)
	if err != nil {
		// TODO: implement p, k, and i
		log.Println("-p, -k, and -i are not currently supported by goracadm")
		return execResponse{}, err
	}

	// validate command flags
	if filename == "" {
		return execResponse{}, errors.New("filename (-f) must be specified")
	}
	if (certType < 1 || certType == 5 || certType > 10) && certType != 16 {
		return execResponse{}, errors.New("cert type must be between 1 and 4, 6 and 10, or 16")
	}

	// open specified file
	cert, err := os.ReadFile(filename)
	if err != nil {
		return execResponse{}, err
	}

	// very basic pem check
	// TODO: maybe verify rest of pem chain
	// This is already doing more than racadm which will allow the upload
	// of ANY file, it seems.
	block, _ := pem.Decode(cert)
	if block == nil || block.Type != "CERTIFICATE" {
		return execResponse{}, errors.New("file is not a pem encoded certificate")
	}
	// validation done

	// file put payload
	filePayload := putfilePayload{
		filename: "RACSSLCERT1",
		flags:    0,
		content:  cert,
	}

	// put the file on the rac
	err = rac.putfile(filePayload)
	if err != nil {
		return execResponse{}, err
	}

	// TODO: racadm executes getconfig here, unsure why

	// exec payload
	cmdInput := fmt.Sprintf("racadm sslcertupload -f sslcertfile -t %d", certType)

	payload := execPayload{}
	payload.Request.CommandInput = cmdInput
	payload.Request.MaxOutputLen = "0x0fff"
	payload.Request.Capability = "0x1"
	payload.Request.UserPrivilege = 0

	// execute payload
	execResp, err = rac.executePayload(payload)
	if err != nil {
		return execResponse{}, err
	}

	return execResp, nil
}
