package idrac

import "flag"

// racresetcfg resets the idrac to factory settings.
// https://www.dell.com/support/manuals/en-us/integrated-dell-remote-access-cntrllr-8-with-lifecycle-controller-v2.00.00.00/racadm_idrac_pub-v1/racresetcfg?guid=guid-bf4676bd-f885-4e20-a7e6-875751246867&lang=en-us
func (rac *idrac) racresetcfg(flags []string) (execResp execResponse, err error) {
	// parse command flags (options)
	fs := flag.NewFlagSet("sslresetcfg", flag.ExitOnError)

	// no flags should be present

	// parse and check for basic errors
	err = parseFlags(fs, flags)
	if err != nil {
		return execResponse{}, err
	}

	// build payload to post to drac
	payload := execPayload{}
	payload.Request.CommandInput = "racadm racresetcfg"
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
