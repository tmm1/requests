package requests

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// Configure creates a new Builder suitable for method chaining by applying the specified Configs.
func Configure(cfgs ...Config) *Builder {
	return (&Builder{}).Config(cfgs...)
}

// Head sets HTTP method to HEAD.
func (rb *Builder) Head() *Builder {
	return rb.Method(http.MethodHead)
}

// Put sets HTTP method to PUT.
func (rb *Builder) Put() *Builder {
	return rb.Method(http.MethodPut)
}

// Patch sets HTTP method to PATCH.
func (rb *Builder) Patch() *Builder {
	return rb.Method(http.MethodPatch)
}

// Delete sets HTTP method to DELETE.
func (rb *Builder) Delete() *Builder {
	return rb.Method(http.MethodDelete)
}

// Hostf calls Host with fmt.Sprintf.
func (rb *Builder) Hostf(format string, a ...interface{}) *Builder {
	return rb.Host(fmt.Sprintf(format, a...))
}

// Pathf calls Path with fmt.Sprintf.
//
// Note that for security reasons, you must not use %s
// with a user provided string!
func (rb *Builder) Pathf(format string, a ...interface{}) *Builder {
	return rb.Path(fmt.Sprintf(format, a...))
}

// ParamInt converts value to a string and calls Param.
func (rb *Builder) ParamInt(key string, value int) *Builder {
	return rb.Param(key, strconv.Itoa(value))
}

// Accept sets the Accept header for a request.
func (rb *Builder) Accept(contentTypes string) *Builder {
	return rb.Header("Accept", contentTypes)
}

// CacheControl sets the client-side Cache-Control directive for a request.
func (rb *Builder) CacheControl(directive string) *Builder {
	return rb.Header("Cache-Control", directive)
}

// ContentType sets the Content-Type header on a request.
func (rb *Builder) ContentType(ct string) *Builder {
	return rb.Header("Content-Type", ct)
}

// UserAgent sets the User-Agent header.
func (rb *Builder) UserAgent(s string) *Builder {
	return rb.Header("User-Agent", s)
}

// BasicAuth sets the Authorization header to a basic auth credential.
func (rb *Builder) BasicAuth(username, password string) *Builder {
	auth := username + ":" + password
	v := base64.StdEncoding.EncodeToString([]byte(auth))
	return rb.Header("Authorization", "Basic "+v)
}

// Bearer sets the Authorization header to a bearer token.
func (rb *Builder) Bearer(token string) *Builder {
	return rb.Header("Authorization", "Bearer "+token)
}

// BodyReader sets the Builder's request body to r.
func (rb *Builder) BodyReader(r io.Reader) *Builder {
	return rb.Body(BodyReader(r))
}

// BodyWriter pipes writes from w to the Builder's request body.
func (rb *Builder) BodyWriter(f func(w io.Writer) error) *Builder {
	return rb.Body(BodyWriter(f))
}

// BodyBytes sets the Builder's request body to b.
func (rb *Builder) BodyBytes(b []byte) *Builder {
	return rb.Body(BodyBytes(b))
}

// BodyJSON sets the Builder's request body to the marshaled JSON.
// It also sets ContentType to "application/json".
func (rb *Builder) BodyJSON(v interface{}) *Builder {
	return rb.
		Body(BodyJSON(v)).
		ContentType("application/json")
}

// BodyForm sets the Builder's request body to the encoded form.
// It also sets the ContentType to "application/x-www-form-urlencoded".
func (rb *Builder) BodyForm(data url.Values) *Builder {
	return rb.
		Body(BodyForm(data)).
		ContentType("application/x-www-form-urlencoded")
}

// BodyFile sets the Builder's request body to read from the given file path.
func (rb *Builder) BodyFile(name string) *Builder {
	return rb.Body(BodyFile(name))
}

// CheckStatus adds a validator for status code of a response.
func (rb *Builder) CheckStatus(acceptStatuses ...int) *Builder {
	return rb.AddValidator(CheckStatus(acceptStatuses...))
}

// CheckContentType adds a validator for the content type header of a response.
func (rb *Builder) CheckContentType(cts ...string) *Builder {
	return rb.AddValidator(CheckContentType(cts...))
}

// CheckPeek adds a validator that peeks at the first n bytes of a response body.
func (rb *Builder) CheckPeek(n int, f func([]byte) error) *Builder {
	return rb.AddValidator(CheckPeek(n, f))
}

// ToJSON sets the Builder to decode a response as a JSON object
func (rb *Builder) ToJSON(v interface{}) *Builder {
	return rb.Handle(ToJSON(v))
}

// ToString sets the Builder to write the response body to the provided string pointer.
func (rb *Builder) ToString(sp *string) *Builder {
	return rb.Handle(ToString(sp))
}

// ToBytesBuffer sets the Builder to write the response body to the provided bytes.Buffer.
func (rb *Builder) ToBytesBuffer(buf *bytes.Buffer) *Builder {
	return rb.Handle(ToBytesBuffer(buf))
}

// ToWriter sets the Builder to copy the response body into w.
func (rb *Builder) ToWriter(w io.Writer) *Builder {
	return rb.Handle(ToWriter(w))
}

// ToFile sets the Builder to write the response body to the given file name.
// The file and its parent directories are created automatically.
// For more advanced use cases, use ToWriter.
func (rb *Builder) ToFile(name string) *Builder {
	return rb.Handle(ToFile(name))
}

// ToHeaders sets the method to HEAD and adds a handler which copies the response headers to h.
func (rb *Builder) ToHeaders(h map[string][]string) *Builder {
	return rb.
		Head().
		Handle(ChainHandlers(ToHeaders(h), consumeBody))
}
