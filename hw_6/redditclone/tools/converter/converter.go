package converter

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

func AnyToBytesBuffer(i interface{}) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(i)
	if err != nil {
		return buf, err
	}
	return buf, nil
}

func ReadBytes(r io.Reader) []byte {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println(err)
	}
	return bytes
}

func AnyBytesToString(i interface{}) *strings.Reader {
	anyJSON, err := AnyToBytesBuffer(i)
	if err != nil {
		log.Println(err)
	}
	return strings.NewReader(anyJSON.String())
}