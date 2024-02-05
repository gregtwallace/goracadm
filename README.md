# Go Racadm
Go Racadm is a recreation of Dell's racadm admin tool used with their
idrac devices. Not nearly all functions are implemented. The posted
binaries are specifically goracadm-cert which is my implementation of
a tool to install key/certificate pem to an idrac and then reset the 
idrac so it loads the new cert.

The tool is implemented using packet captures of racadm 9.1.2 and the
cloning the functionality in Go.

The idrac package can also be imported into other Go programs if
interfacing with an idrac is needed, as opposed to a racadm
reimplementation.

## Compatibility Notice

I only have an idrac 7 to test with. I previously had an idrac 6 and it also
was working with this tool.

Your mileage may vary.

## Subcommands Implemented in the IDRAC package (so far)
Subcommands:
racreset,
racresetcfg,
sslcertdownload,
sslcertupload,
sslkeyupload,
sslresetcfg

## Usage

Run the tool as:

`./goracadm-cert --hostname idrac.example.com --username someone --password secret --keyfile key.pem --certfile cert.pem [FLAGS]`

Help can be viewed with:

`./goracadm-cert --help`

## Note About Install Automation

The application supports passing all args instead as environment 
variables by prefixing the flag name with `GORACADM_CERT`. 

e.g. `GORACADM_CERT_KEYPEM`

There are mutually exclusive flags that allow specifying the pem 
as either filenames or directly as strings. The strings are useful 
for passing the pem content from another application without having 
to save the pem files to disk.

Putting all of this together, you can combine the install binary with 
a tool like LeGo CertHub (https://www.legocerthub.com/) to call the 
install binary, with environment variables, to directly upload new 
certificates as they're issued by LeGo, without having to write a 
separate script.

![LeGo CertHub with GoRacAdm Cert](https://raw.githubusercontent.com/gregtwallace/goracadm/main/img/goracadm-cert.png)
