package main

import (
	"context"
	"errors"
	"io"
	"log"
	"os"

	"github.com/gregtwallace/goracadm/pkg/idrac"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
)

// struct for receivers to use common app pieces
type app struct {
	stdLogger   *log.Logger
	debugLogger *log.Logger
	errLogger   *log.Logger
	cmd         *ff.Command
	config      *config
}

// a binary that accepts args or environment variables and executes the
// sequential commands to install an ssl key and certificate and then
// restart the idrac
func main() {
	// make app w/ logger
	app := &app{
		stdLogger:   log.New(os.Stdout, "", 0),
		debugLogger: log.New(io.Discard, "", 0), // discard debug logging by default
		errLogger:   log.New(os.Stderr, "", 0),
	}

	// log start
	app.stdLogger.Printf("goracadm-cert v%s", idrac.Version)

	// get & parse config
	err := app.getConfig()

	// if debug logging, make real debug logger
	if app.config.debugLogging != nil && *app.config.debugLogging {
		app.debugLogger = log.New(os.Stdout, "debug: ", 0)
	}

	// deal with config err (after logger re-init)
	if err != nil {
		exitCode := 0

		if errors.Is(err, ff.ErrHelp) {
			// help explicitly requested
			app.stdLogger.Printf("\n%s\n", ffhelp.Command(app.cmd))

		} else if errors.Is(err, ff.ErrDuplicateFlag) ||
			errors.Is(err, ff.ErrUnknownFlag) ||
			errors.Is(err, ff.ErrNoExec) ||
			errors.Is(err, ErrExtraArgs) {
			// other error that suggests user needs to see help
			exitCode = 1
			app.errLogger.Print(err)
			app.stdLogger.Printf("\n%s\n", ffhelp.Command(app.cmd))

		} else {
			// any other error
			exitCode = 1
			app.errLogger.Print(err)
		}

		os.Exit(exitCode)
	}

	// run it
	exitCode := 0
	err = app.cmd.Run(context.Background())
	if err != nil {
		exitCode = 1
		app.errLogger.Print(err)

		// if extra args, show help
		if errors.Is(err, ErrExtraArgs) {
			app.stdLogger.Printf("\n%s\n", ffhelp.Command(app.cmd))
		}
	}

	app.stdLogger.Print("goracadm-cert done")
	os.Exit(exitCode)
}
