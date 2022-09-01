package main

import (
	"flag"
	"log"

	"github.com/gregtwallace/goracadm/idrac"
)

func main() {
	// config options
	hostname := ""
	username := ""
	password := ""

	// parse command line
	flag.StringVar(&hostname, "r", "", "idrac hostname or ip address (and port)")
	flag.StringVar(&username, "u", "", "idrac username")
	flag.StringVar(&password, "p", "", "idrac password")

	flag.Parse()

	// make idrac
	rac := idrac.NewIdrac(hostname, username, password)

	// do discover (confirm hostname is actually an idrac)
	_, err := rac.Discover()
	if err != nil {
		log.Fatal(err)
	}

	// login to idrac and save the sid cookie
	_, err = rac.Login()
	if err != nil {
		log.Fatal(err)
	}

	// execute the subcommand
	_, err = rac.Exec()
	if err != nil {
		// not fatal, continue to logout
		log.Printf("execution error: %s", err)
	}

	// logout of the idrac
	_, err = rac.Logout()
	if err != nil {
		log.Fatal(err)
	}
}
