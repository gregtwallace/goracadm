package idrac

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net/http"
	"strconv"
)

const endpointPutfile = "/cgi-bin/putfile"

// putfilePayload is a struct to represent the payload
// the idrac expects when putting a file.
type putfilePayload struct {
	filename string // file name on idrac
	flags    int    // file flags (unsure of purpose)
	content  []byte // the actual file content being put
}

// Bytes() translates the putfilePayload into the byte slice
// that will be posted to the idrac
// composition: [32]byte filename, [4]byte file content length,
// [4]byte flags, []byte file content
// props to https://github.com/KraudSecurity/Exploits/blob/master/CVE-2018-1207/CVE-2018-1207.py
// for helping me to understand the putfile byte data
func (payload *putfilePayload) bytes() []byte {
	// name and content as bytes
	name := []byte(payload.filename)

	// normalize new line style in file content (racadm doesn't do
	// this but doing for consistency)
	// windows
	payload.content = bytes.Replace(payload.content, []byte{13, 10}, []byte{10}, -1)
	// mac
	payload.content = bytes.Replace(payload.content, []byte{13}, []byte{10}, -1)

	// calc file content len and encode to little endian
	fileContentLen := len(payload.content)
	len := make([]byte, 4)
	binary.LittleEndian.PutUint32(len, uint32(fileContentLen))

	// flags encoded to little endian
	flags := make([]byte, 4)
	binary.LittleEndian.PutUint32(flags, uint32(payload.flags))

	// total length of the byte payload
	payloadLen := 32 + 4 + 4 + fileContentLen

	// allocate the byte slice
	data := make([]byte, payloadLen)

	// fill the data slice
	_ = copy(data, name)
	_ = copy(data[32:], len)
	_ = copy(data[36:], flags)
	_ = copy(data[40:], payload.content)

	return data
}

// putfile sends the specified data as an octetstream to the rac's
// putfile endpoint
func (rac *idrac) putfile(payload putfilePayload) (err error) {
	// post
	payloadBytes := payload.bytes()
	// log.Println(string(payloadBytes))

	resp, err := rac.client.Post(rac.url()+endpointPutfile, "application/octet-stream", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	// check status code 200 (OK)
	if resp.StatusCode != http.StatusOK {
		return errors.New("error: http status code " + strconv.Itoa(resp.StatusCode))
	}

	return nil
}
