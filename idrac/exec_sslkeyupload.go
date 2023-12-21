package idrac

import (
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"os"
)

// sslkeyupload executes the sslkeyupload subcommand using
// the specified flags.
// https://www.dell.com/support/manuals/en-us/poweredge-m630/idrac8_2.70.70.70_racadm/sslkeyupload?guid=guid-293e0da4-1ed3-4ed3-9363-f3091c0ecd1c&lang=en-us
// https://www.dell.com/support/manuals/en-us/idrac9-lifecycle-controller-v5.x-series/idrac9_5.xx_racadm_pub/sslkeyupload?guid=guid-27e877c9-5ede-41c5-975f-497bc7443555&lang=en-us
func (rac *idrac) sslkeyupload(flags []string) (execResp execResponse, err error) {
	// parse command flags (options)
	file := ""
	certType := 0

	fs := flag.NewFlagSet("sslkeyupload", flag.ExitOnError)
	fs.StringVar(&file, "f", "", "local filename to upload or pem string of cert (required)")
	fs.IntVar(&certType, "t", 0, "certificate type (required) (only 1 is valid)")

	// parse and check for basic errors
	err = parseFlags(fs, flags)
	if err != nil {
		return execResponse{}, err
	}

	// validate command flags
	if file == "" {
		return execResponse{}, errors.New("file (-f) must be specified")
	}
	if certType == 0 {
		return execResponse{}, errors.New("cert type (-t) must be specified")
	}
	if certType != 1 {
		return execResponse{}, errors.New("cert type (-t) must be 1")
	}

	// MODIFIED BEHAVIOR FROM racadm, though still fully compliant with spec
	// try to parse file as pem content
	pemBlock, _ := pem.Decode([]byte(file))
	if pemBlock != nil {
		// file input is valid pem string (as opposed to file name)
		// no-op
	} else {
		// if failed to parse file as pem content, do normal behavior of trying to open the filename and read it
		keyFileBytes, err := os.ReadFile(file)
		if err != nil {
			return execResponse{}, err
		}

		// confirm file content is valid pem (discards any "extra" content after key block)
		pemBlock, _ = pem.Decode(keyFileBytes)
		if pemBlock == nil || (pemBlock.Type != "PRIVATE KEY" && pemBlock.Type != "RSA PRIVATE KEY") {
			return execResponse{}, errors.New("file is not a pem encoded private key")
		}
	}

	// file put payload
	filePayload := putfilePayload{
		filename: "RACSSLCERT1",
		flags:    0,
		content:  pem.EncodeToMemory(pemBlock),
	}

	// put the file on the rac
	err = rac.putfile(filePayload)
	if err != nil {
		return execResponse{}, err
	}

	// exec payload
	cmdInput := fmt.Sprintf("racadm sslkeyupload -f sslcertfile -t %d", certType)

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
