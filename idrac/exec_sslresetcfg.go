package idrac

import (
	"flag"
)

// racreset resets the idrac using the specified flags.
// https://www.dell.com/support/manuals/en-us/poweredge-m630/idrac8_2.70.70.70_racadm/racreset?guid=guid-7866bef3-f5c4-4c8b-b2e3-ce22d6332ddc&lang=en-us
// https://www.dell.com/support/manuals/en-us/idrac9-lifecycle-controller-v5.x-series/idrac9_5.xx_racadm_pub/racreset?guid=guid-a5b943ea-b4b5-415a-bd3c-09a02dfed465&lang=en-us
func (rac *idrac) sslresetcfg(flags []string) (execResp execResponse, err error) {
	// parse command flags (options)
	fs := flag.NewFlagSet("sslresetcfg", flag.ExitOnError)

	// no flags should be present

	// parse and check for basic errors
	err = parseFlags(fs, flags)
	if err != nil {
		return execResponse{}, err
	}

	// TODO: racadm executes getconfig here, unsure why

	// build payload to post to drac
	payload := execPayload{}
	payload.Request.CommandInput = "racadm sslresetcfg"
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
