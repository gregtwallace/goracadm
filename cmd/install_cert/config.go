package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/peterbourgon/ff/v4"
)

var (
	ErrExtraArgs = errors.New("extra args present")

	environmentVarPrefix = "GORACADM_CERT"
)

// keyCertPemCfg contains values common to subcommands that need to use key
// and cert pem
type keyCertPemCfg struct {
	keyPemFilePath  *string
	certPemFilePath *string
	keyPem          *string
	certPem         *string
}

// app's config options from user
type config struct {
	hostname *string
	username *string
	password *string
	keyCertPemCfg
	insecure     *bool
	debugLogging *bool
}

// getConfig returns the app's configuration from either command line args,
// or environment variables
func (app *app) getConfig() error {
	// make config
	cfg := &config{}

	// goracadm-cert -- root command
	rootFlags := ff.NewFlagSet("goracadm-cert")

	cfg.hostname = rootFlags.StringLong("hostname", "", "the hostname of the remote idrac")
	cfg.username = rootFlags.StringLong("username", "", "the username to login to the remote idrac")
	cfg.password = rootFlags.StringLong("password", "", "the password to login to the remote idrac")
	cfg.keyPemFilePath = rootFlags.StringLong("keyfile", "", "path and filename of the rsa-2048 key in pem format")
	cfg.certPemFilePath = rootFlags.StringLong("certfile", "", "path and filename of the certificate in pem format")
	cfg.keyPem = rootFlags.StringLong("keypem", "", "string of the rsa-2048 key in pem format")
	cfg.certPem = rootFlags.StringLong("certpem", "", "string of the certificate in pem format")
	cfg.insecure = rootFlags.BoolLong("insecure", "disable https certificate validation (DANGEROUS)")
	cfg.debugLogging = rootFlags.BoolLong("debug", "set this flag to enable additional debug logging messages")

	rootCmd := &ff.Command{
		Name:      "goracadm-cert",
		Usage:     "goracadm-cert --hostname idrac.example.com --username someone --password secret --keyfile key.pem --certfile cert.pem [FLAGS]",
		ShortHelp: "install the specified key and cert pem files on an idrac and reset the idrac to load the new key/cert",
		Flags:     rootFlags,
		Exec:      app.cmdInstallCertAndReset,
	}

	// set cfg & parse
	app.config = cfg
	app.cmd = rootCmd
	err := app.cmd.Parse(os.Args[1:], ff.WithEnvVarPrefix(environmentVarPrefix))
	if err != nil {
		return err
	}

	return nil
}

// GetPemBytes returns the key and cert pem bytes as specified in keyCertPemCfg
// or an error if it cant get the bytes of both
func (kcCfg *keyCertPemCfg) GetPemBytes(subcommand string) (keyPem, certPem []byte, err error) {
	// key pem (from arg or file)
	if kcCfg.keyPem != nil && *kcCfg.keyPem != "" {
		// error if filename is also set
		if kcCfg.keyPemFilePath != nil && *kcCfg.keyPemFilePath != "" {
			return nil, nil, fmt.Errorf("%s: failed, both key pem and key file specified", subcommand)
		}

		// use pem
		keyPem = []byte(*kcCfg.keyPem)
	} else {
		// pem wasn't specified, try reading file
		if kcCfg.keyPemFilePath == nil || *kcCfg.keyPemFilePath == "" {
			return nil, nil, fmt.Errorf("%s: failed, neither key pem nor key file specified", subcommand)
		}

		// read file to get pem
		keyPem, err = os.ReadFile(*kcCfg.keyPemFilePath)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: failed to read key file (%w)", subcommand, err)
		}
	}

	// cert pem (repeat same process)
	if kcCfg.certPem != nil && *kcCfg.certPem != "" {
		// error if filename is also set
		if kcCfg.certPemFilePath != nil && *kcCfg.certPemFilePath != "" {
			return nil, nil, fmt.Errorf("%s: failed, both cert pem and cert file specified", subcommand)
		}

		// use pem
		certPem = []byte(*kcCfg.certPem)
	} else {
		// pem wasn't specified, try reading file
		if kcCfg.certPemFilePath == nil || *kcCfg.certPemFilePath == "" {
			return nil, nil, fmt.Errorf("%s: failed, neither cert pem nor cert file specified", subcommand)
		}

		// read file to get pem
		certPem, err = os.ReadFile(*kcCfg.certPemFilePath)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: failed to read cert file (%w)", subcommand, err)
		}
	}

	return keyPem, certPem, nil
}
