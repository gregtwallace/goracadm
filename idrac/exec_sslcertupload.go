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
	file := ""
	certType := 0

	fs := flag.NewFlagSet("sslcertupload", flag.ExitOnError)
	fs.StringVar(&file, "f", "", "local filename to upload or pem string of cert (required)")
	fs.IntVar(&certType, "t", 0, "certificate type (required - int - see Dell docs)")

	// parse and check for basic errors
	err = parseFlags(fs, flags)
	if err != nil {
		// TODO: implement p, k, and i
		log.Println("-p, -k, and -i are not currently supported by goracadm")
		return execResponse{}, err
	}

	// validate command flags
	if file == "" {
		return execResponse{}, errors.New("filename (-f) must be specified")
	}
	if (certType < 1 || certType == 5 || certType > 10) && certType != 16 {
		return execResponse{}, errors.New("cert type must be between 1 and 4, 6 and 10, or 16")
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
		if pemBlock == nil || pemBlock.Type != "CERTIFICATE" {
			return execResponse{}, errors.New("file is not a pem encoded certificate")
		}
	}

	// TODO: maybe verify rest of pem chain

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
