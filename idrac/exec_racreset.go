package idrac

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

var (
	errInvalidFirmness = errors.New("invalid first option (must either be soft, hard, -m, of -f)")
	errInvalidModule   = errors.New("invalid module (-m) option")
)

// racreset resets the idrac using the specified flags.
// https://www.dell.com/support/manuals/en-us/poweredge-m630/idrac8_2.70.70.70_racadm/racreset?guid=guid-7866bef3-f5c4-4c8b-b2e3-ce22d6332ddc&lang=en-us
// https://www.dell.com/support/manuals/en-us/idrac9-lifecycle-controller-v5.x-series/idrac9_5.xx_racadm_pub/racreset?guid=guid-a5b943ea-b4b5-415a-bd3c-09a02dfed465&lang=en-us
func (rac *idrac) racreset(flags []string) (execResp execResponse, err error) {
	// check for non-flag firmness before parsing flags
	firmnessParam := ""
	if len(flags) > 0 {
		if flags[0] == "soft" || flags[0] == "hard" {
			firmnessParam = " " + flags[0]
			flags = flags[1:]
		} else if !strings.HasPrefix(flags[0], "-m") && !strings.HasPrefix(flags[0], "-f") {
			return execResponse{}, errInvalidFirmness
		}
	}

	// parse command flags (options)
	force := false
	// TODO: Implement support for multiple modules to be specified
	module := ""

	fs := flag.NewFlagSet("racreset", flag.ExitOnError)

	fs.BoolVar(&force, "f", false, "This option is used to force the reset.")
	fs.StringVar(&module, "m", "", "server-<n> — where n=1-16	-or- server-<nx> — where n=1-8; x = a, b, c, d (lower case)")

	// parse and check for basic errors
	err = parseFlags(fs, flags)
	if err != nil {
		return execResponse{}, err
	}

	// validate command flags and build command
	forceParam := ""
	if force {
		forceParam = " -f"
	}

	moduleParam := ""
	if module != "" {
		// must start with server
		if !strings.HasPrefix(module, "server-") {
			return execResponse{}, errInvalidModule
		}

		// remove server prefix
		module = strings.TrimPrefix(module, "server-")

		// valid remaining is either 1 or 2 chars
		// for 1 char, must be number 1-9
		if len(module) == 1 {
			srvNumb, err := strconv.Atoi(module)
			if err != nil {
				return execResponse{}, errInvalidModule
			}
			if srvNumb >= 1 && srvNumb <= 9 {
				// valid
				moduleParam = " -m server-" + module
			}

		} else if len(module) == 2 {
			// for two char, can either be 10-16, or 1-8 a,b,c,d
			// address a,b,c,d
			if string(module[1]) == "a" || string(module[1]) == "b" || string(module[1]) == "c" || string(module[1]) == "d" {
				srvNumb, err := strconv.Atoi(string(module[0]))
				if err != nil {
					return execResponse{}, errInvalidModule
				}
				if srvNumb >= 1 && srvNumb <= 8 {
					// valid
					moduleParam = " -m server-" + module
				} else {
					return execResponse{}, errInvalidModule
				}
			} else {
				// address 10-16
				srvNumb, err := strconv.Atoi(module)
				if err != nil {
					return execResponse{}, errInvalidModule
				}
				if srvNumb >= 10 && srvNumb <= 16 {
					// valid
					moduleParam = " -m server-" + module
				}
			}
		} else {
			return execResponse{}, errInvalidModule
		}
	}

	// build payload to post to drac
	cmdInput := fmt.Sprintf("racadm racreset%s%s%s", firmnessParam, forceParam, moduleParam)

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
