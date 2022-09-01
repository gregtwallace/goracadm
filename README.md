# Go Racadm
Go Racadm is a recreation of Dell's racadm admin tool used with their
idrac devices. Not nearly all functions are implemented. If you need
a specific function, please submit a request or pull request.

The tool is implemented using packet captures of racadm 9.1.2 and the
cloning the functionality in Go.

The idrac package can also be imported into other Go programs if
interfacing with an idrac is needed, as opposed to a racadm
reimplementation.

Disclaimer: I only have idrac 6 & 7 to test with.

## Implemented (so far)
Subcommands:
sslcertdownload
