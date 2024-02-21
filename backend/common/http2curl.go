package common

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

// CurlCommand contains exec.Command compatible slice + helpers
type CurlCommand struct {
	Slice []string
}

// append appends a string to the CurlCommand
func (c *CurlCommand) append(newSlice ...string) {
	c.Slice = append(c.Slice, newSlice...)
}

// String returns a ready to copy/paste command
func (c *CurlCommand) String() string {
	return strings.Join(c.Slice, " ")
}

// nopCloser is used to create a new io.ReadCloser for req.Body
type nopCloser struct {
	io.Reader
}

func bashEscape(str string) string {
	return `$'` + strings.Replace(str, `'`, `'\''`, -1) + `'`
}

func bodyEscape(str string) string {
	return `$'` + strings.Replace(str, "\r\n", "\\r\\n", -1) + `'`
}

func (nopCloser) Close() error { return nil }

// GetCurlCommand returns a CurlCommand corresponding to an http.Request
func GetCurlCommand(req *http.Request) (*CurlCommand, error) {
	if req.URL == nil {
		return nil, fmt.Errorf("getCurlCommand: invalid request, req.URL is nil")
	}

	command := CurlCommand{}

	command.append("curl")

	command.append("-i")

	command.append("-s")

	command.append("-k")

	// command.append("-v")

	command.append("--compressed")

	command.append("-m 10")

	command.append("-X", bashEscape(req.Method))

	var keys []string

	for k := range req.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		command.append("-H", bashEscape(fmt.Sprintf("%s: %s", k, strings.Join(req.Header[k], " "))))
	}

	if req.Body != nil {
		// TODO unexpected EOF
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body = nopCloser{bytes.NewBuffer(body)}

		bodyEscaped := bodyEscape(string(body))
		command.append("-d", bodyEscaped)
	}

	command.append(bashEscape(req.URL.String()))

	return &command, nil
}
