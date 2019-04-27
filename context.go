package router

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

// ContextData is map of key-value data that can be stored in global context
type ContextData map[string]interface{}

// URLParams is map of key-value parameters that can be part of url as placeholders
type URLParams map[string]string

// Context type type
type Context struct {
	Res          http.ResponseWriter
	Req          *http.Request
	Params       URLParams
	Query        interface{}
	Method       string
	URI          string
	data         ContextData
	responseSent bool // Flag is response has been already written, used to avoid write headers multiple times
}

type response struct {
	Status      int
	ContentType string
	Data        []byte
}

func createContext(response http.ResponseWriter, request *http.Request) (ctx *Context) {
	request.ParseMultipartForm(32 << 20) // parse form data

	// prepare url
	url := request.URL.Path
	url = strings.TrimSuffix(url, "/")       // trim trailing /
	url = "/" + strings.TrimPrefix(url, "/") // make sure every route starts with /

	ctx = &Context{
		Res:          response,
		Req:          request,
		Params:       make(URLParams),
		Query:        request.URL.Query(),
		Method:       request.Method,
		URI:          url,
		data:         make(ContextData),
		responseSent: false,
	}
	return
}

func (ctx *Context) JSON(parsed interface{}) error {
	// check if content-type is json
	if ctx.Req.Header.Get("Content-Type") != ContentType["JSON"] {
		return errors.New("Content-Type is not JSON")
	}

	// read body
	body, err := ioutil.ReadAll(ctx.Req.Body)
	defer ctx.Req.Body.Close()
	if err != nil {
		return err
	}

	// parse json
	json.Unmarshal(body, parsed)
	return nil
}

func (ctx *Context) Error(err *HTTPError) {
	errResponse := &response{
		Status:      err.Code,
		ContentType: ContentType["JSON"],
		Data:        err.GetResponseData(),
	}
	ctx.writeResponse(errResponse)
}

// Send response to client
func (ctx *Context) Send(params ...interface{}) {
	if len(params) < 1 {
		panic("Provide data to send")
	}

	// set HTTP response status if not provided
	var status = http.StatusOK
	if len(params) == 2 {
		status = params[1].(int)
	}

	// Build response
	var data ResponseData = params[0]
	responseBuilder := getResponseBuilder(ctx, data)
	serialized, err := responseBuilder.Serialize()
	if err != nil {
		ctx.Error(InternalError("Failed to serialize output data"))
		return
	}
	ctx.writeResponse(&response{
		Status:      status,
		ContentType: responseBuilder.GetContentType(),
		Data:        serialized,
	})
}

func (ctx *Context) writeResponse(res *response) {
	if !ctx.responseSent {
		ctx.Res.Header().Add("Content-Type", res.ContentType)
		ctx.Res.WriteHeader(res.Status)
		ctx.Res.Write(res.Data)
		ctx.responseSent = true
	}
}

// Set value to context
func (ctx *Context) Set(key string, data interface{}) {
	ctx.data[key] = data
}

// Get value stored in context by key
func (ctx *Context) Get(key string) (interface{}, bool) {
	data, ok := ctx.data[key]
	if !ok {
		return nil, false
	}
	return data, true
}
