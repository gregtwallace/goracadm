package main

import (
	"context"
	"fmt"

	"github.com/gregtwallace/goracadm/pkg/idrac"
)

// cmdInstallCertAndReset executes a series of commands against an idrac to install
// the specified ssl key and cert. it then resets the idrac so it will load the
// newly installed key/cert
func (app *app) cmdInstallCertAndReset(_ context.Context, args []string) error {
	// extra args == error
	if len(args) != 0 {
		return fmt.Errorf("main: failed, %w (%d)", ErrExtraArgs, len(args))
	}

	// load key and cert
	keyPem, certPem, err := app.config.keyCertPemCfg.GetPemBytes("main")
	if err != nil {
		return err
	}

	// validate ssl?
	strictCerts := true
	if app.config.insecure != nil && *app.config.insecure {
		app.stdLogger.Println("WARNING: --insecure flag set, https certificate will not be validated")
		strictCerts = false
	}

	// make idrac
	rac, err := idrac.NewIdrac(*app.config.hostname, *app.config.username, *app.config.password, strictCerts)
	if err != nil {
		return err
	}

	// do discover (confirm hostname is actually an idrac)
	_, err = rac.Discover()
	if err != nil {
		return err
	}

	// login to idrac and save the sid cookie
	_, err = rac.Login()
	if err != nil {
		return fmt.Errorf("login error: %w", err)
	}

	// execute 3 commands: sslkeyupload, sslcertupload, racreset
	// sslkeyupload
	_, err = rac.Exec("sslkeyupload", []string{"-t", "1", "-f", string(keyPem)})
	if err != nil {
		return fmt.Errorf("failed to upload key (%w)", err)
	}

	// sslcertupload
	_, err = rac.Exec("sslcertupload", []string{"-t", "1", "-f", string(certPem)})
	if err != nil {
		return fmt.Errorf("failed to upload cert (%w)", err)
	}

	// racreset
	_, err = rac.Exec("racreset", nil)
	if err != nil {
		return fmt.Errorf("failed to reset (%w)", err)
	}

	// logout of the idrac
	_, _ = rac.Logout()
	// don't worry about error
	// an error isn't too concerning as rac may reset before logout actually processes

	return nil
}
