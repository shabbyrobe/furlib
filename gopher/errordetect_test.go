package gopher

import (
	"fmt"
	"testing"
)

func TestErrorDetect(t *testing.T) {
	for idx, tc := range []struct {
		in string
	}{
		{`3'/caps.txt' does not exist (no handler found)		error.host	1`},
		{"3 '/caps.txt' doesn't exist!		error.host	1"},
		{"3`caps.txt' invalid.	Error	Error	0"},
		{`3'caps.txt' No such file or directory (2)`},
		{`3caps.txt NOT FOUND`},
		{`3Error accessing /caps.txt.		error.host	1`},
		{`3 File not found`},
		{`3file not found	fake	(NULL)	0`},
		{`3Happy helping ☃ here: Sorry, your selector contains '..'. That's illegal here.	Err	localhost	70`},
		{`3Happy helping ☃ here: Sorry, your selector does not start with / or contains '..'. That's illegal here.	Err	localhost	70`},
		{`3Malformed request	fakeselector	fakeserver	70`},
		{`3not found	(NULL)	error.host	0`},
		{`3open path/to/caps.txt: no such file or directory		error.host	1`},
		{`3Sorry, but the requested token 'caps.txt' could not be found.	Err	localhost	70`},
		{`3Sorry! I could not find caps.txt`},
		{`3The provided selector is invalid.		example.com	70`},
		{`3"/usr/stevie/gophercaps.txt" not found	error.file	error.host	0`},
		{`An error occurred: Resource not found.`},
		{`Error: 404 Not Found`},
		{`Error: File or directory not found!`},
		{`Error: Page not found	example.com	70`},
		{`Error: resource caps.txt does not exist on example.com`},
		{`File: '/caps.txt' not found.`},
		{`Error: 404 Not Found\n\nThe requested URL was not found on the server. If you entered the URL\nmanually please check your spelling and try again.`},
		{"--404\r\nNot Found\r\n.\r\n"},

		{`` +
			`i   ____            _       ____      _ _` + "\n" +
			`i  |  _ \ _   _ ___| |_ ___|  _ \  __| | | __` + "\n" +
			"i  | | | | | | / __| __/ _ | | | |/ _` | |/ /" + "\n" +
			`i  | |_| | |_| \__ \ ||  __/ |_| | (_| |   < ` + "\n" +
			`i  |____/ \__,_|___/\__\___|____(_)__,_|_|\_\` + "\n" +
			`i                    - a strange place indeed` + "\n" +
			`i ` + "\n" +
			`i ` + "\n" +
			`3Sorry! I could not find caps.txt`,
		},

		{`` + // gopher://mozz.us:7005/1/error/403/menu
			`iError: 403 Forbidden	fake	example.com	0` + "\r\n" +
			`i	fake	example.com	0` + "\r\n" +
			`iYou don't have the permission to access the requested resource. It is	fake	example.com	0` + "\r\n" +
			`ieither read-protected or not readable by the server.	fake	example.com	0` + "\r\n",
		},
	} {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			fn := func(status Status, msg string, confidence float64) *Error {
				return NewError(URL{}, status, msg, confidence)
			}
			if DetectError([]byte(tc.in), fn) == nil {
				t.Fatal(tc.in)
			}
		})
	}
}

func TestExtractGopherII(t *testing.T) {
	for idx, tc := range []struct {
		in     string
		status Status
		msg    string
		found  bool
	}{
		{"--404\r\nNot Found\r\n.\r\n", StatusNotFound, "Not Found", true},
	} {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			s, m, f := extractGopherIIError([]byte(tc.in))
			if f != tc.found {
				t.Fatal()
			}
			if s != tc.status {
				t.Fatal(s)
			}
			if m != tc.msg {
				t.Fatal()
			}
		})
	}
}
