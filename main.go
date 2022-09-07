package main

import (
	"flag"
	"log"
	"os"

	"github.com/gregtwallace/goracadm/idrac"
)

const version = "0.1.0"

func main() {
	log.Printf("goracadm v.%s", version)

	// exit code
	exitCode := 0

	// config options
	hostname := ""
	username := ""
	password := ""
	strictCerts := false

	// parse command line
	flag.StringVar(&hostname, "r", "", "idrac hostname or ip address (and port)")
	flag.StringVar(&username, "u", "", "idrac username")
	flag.StringVar(&password, "p", "", "idrac password")
	flag.BoolVar(&strictCerts, "S", false, "strictly require validated certs")

	flag.Parse()

	// make idrac
	rac, err := idrac.NewIdrac(hostname, username, password, strictCerts)
	if err != nil {
		log.Fatal(err)
	}

	// do discover (confirm hostname is actually an idrac)
	_, err = rac.Discover()
	if err != nil {
		log.Fatal(err)
	}

	// login to idrac and save the sid cookie
	_, err = rac.Login()
	if err != nil {
		log.Fatalf("login error: %s", err)
	}

	// get subcommand and flags
	cmd := flag.Args()[0]
	flags := flag.Args()[1:]

	// execute the subcommand
	_, err = rac.Exec(cmd, flags)
	if err != nil {
		// not fatal, continue to logout and change exit code to error
		log.Printf("exec error: %s", err)
		exitCode = 1
	}

	// logout of the idrac
	_, err = rac.Logout()
	if err != nil {
		// not fatal, nor change exit code (could result from things like
		// success of racreset)
		log.Printf("logout error: %s", err)
	}

	// exit with appropriate code
	log.Printf("goracadm exit code: %d", exitCode)
	os.Exit(exitCode)
}
