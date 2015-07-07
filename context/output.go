// Copyright 2014 beego Author. All Rights Reserved.
// 2015 (c) Dmitriy Blokhin (sv.dblokhin@gmail.com), www.webjinn.ru

package context

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

// Output does work for sending response header.
type Output struct {
	response http.ResponseWriter
	Status   int
}

// NewOutput returns new Output.
// it contains nothing now.
func NewOutput(resp http.ResponseWriter) *Output {
	return &Output{resp, http.StatusOK}
}

// Response возвращает http.ResponseWriter
func (output *Output) Response() http.ResponseWriter {
	return output.response
}

// Header sets response header item string via given key.
func (output *Output) Header(key, val string) {
	output.response.Header().Set(key, val)
}

// Cookie sets cookie value via given key.
// others are ordered as cookie's max age time, path,domain, secure and httponly.
func (output *Output) Cookie(name string, value string, others ...interface{}) {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s=%s", sanitizeName(name), sanitizeValue(value))
	if len(others) > 0 {
		switch v := others[0].(type) {
		case int:
			if v > 0 {
				fmt.Fprintf(&b, "; Max-Age=%d", v)
			} else if v < 0 {
				fmt.Fprintf(&b, "; Max-Age=0")
			}
		case int64:
			if v > 0 {
				fmt.Fprintf(&b, "; Max-Age=%d", v)
			} else if v < 0 {
				fmt.Fprintf(&b, "; Max-Age=0")
			}
		case int32:
			if v > 0 {
				fmt.Fprintf(&b, "; Max-Age=%d", v)
			} else if v < 0 {
				fmt.Fprintf(&b, "; Max-Age=0")
			}
		}
	}

	// the settings below
	// Path, Domain, Secure, HttpOnly
	// can use nil skip set

	// default "/"
	if len(others) > 1 {
		if v, ok := others[1].(string); ok && len(v) > 0 {
			fmt.Fprintf(&b, "; Path=%s", sanitizeValue(v))
		}
	} else {
		fmt.Fprintf(&b, "; Path=%s", "/")
	}

	// default empty
	if len(others) > 2 {
		if v, ok := others[2].(string); ok && len(v) > 0 {
			fmt.Fprintf(&b, "; Domain=%s", sanitizeValue(v))
		}
	}

	// default empty
	if len(others) > 3 {
		var secure bool
		switch v := others[3].(type) {
		case bool:
			secure = v
		default:
			if others[3] != nil {
				secure = true
			}
		}
		if secure {
			fmt.Fprintf(&b, "; Secure")
		}
	}

	// default false. for session cookie default true
	httponly := false
	if len(others) > 4 {
		if v, ok := others[4].(bool); ok && v {
			// HttpOnly = true
			httponly = true
		}
	}

	if httponly {
		fmt.Fprintf(&b, "; HttpOnly")
	}

	output.response.Header().Add("Set-Cookie", b.String())
}

var cookieNameSanitizer = strings.NewReplacer("\n", "-", "\r", "-")

func sanitizeName(n string) string {
	return cookieNameSanitizer.Replace(n)
}

var cookieValueSanitizer = strings.NewReplacer("\n", " ", "\r", " ", ";", " ")

func sanitizeValue(v string) string {
	return cookieValueSanitizer.Replace(v)
}

// SetStatus sets response status code.
// It writes response header directly.
func (output *Output) SetStatus(status int) {
	output.Status = status
}

// IsCachable returns boolean of this request is cached.
// HTTP 304 means cached.
func (output *Output) IsCachable(status int) bool {
	return output.Status >= 200 && output.Status < 300 || output.Status == 304
}

// IsEmpty returns boolean of this request is empty.
// HTTP 201，204 and 304 means empty.
func (output *Output) IsEmpty(status int) bool {
	return output.Status == 201 || output.Status == 204 || output.Status == 304
}

// IsOk returns boolean of this request runs well.
// HTTP 200 means ok.
func (output *Output) IsOk(status int) bool {
	return output.Status == 200
}

// IsSuccessful returns boolean of this request runs successfully.
// HTTP 2xx means ok.
func (output *Output) IsSuccessful(status int) bool {
	return output.Status >= 200 && output.Status < 300
}

// IsRedirect returns boolean of this request is redirection header.
// HTTP 301,302,307 means redirection.
func (output *Output) IsRedirect(status int) bool {
	return output.Status == 301 || output.Status == 302 || output.Status == 303 || output.Status == 307
}

// IsForbidden returns boolean of this request is forbidden.
// HTTP 403 means forbidden.
func (output *Output) IsForbidden(status int) bool {
	return output.Status == 403
}

// IsNotFound returns boolean of this request is not found.
// HTTP 404 means forbidden.
func (output *Output) IsNotFound(status int) bool {
	return output.Status == 404
}

// IsClient returns boolean of this request client sends error data.
// HTTP 4xx means forbidden.
func (output *Output) IsClientError(status int) bool {
	return output.Status >= 400 && output.Status < 500
}

// IsServerError returns boolean of this server handler errors.
// HTTP 5xx means server internal error.
func (output *Output) IsServerError(status int) bool {
	return output.Status >= 500 && output.Status < 600
}
