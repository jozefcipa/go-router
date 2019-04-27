package router

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type ResponseBuilder struct {
	ctx  *Context
	data ResponseData
}

type ResponseData interface{}

func getResponseBuilder(ctx *Context, data ResponseData) *ResponseBuilder {
	return &ResponseBuilder{ctx: ctx, data: data}
}

// Serialize data structure for response
func (builder ResponseBuilder) Serialize() ([]byte, error) {
	switch builder.GetContentType() {
	case ContentType["JSON"]:
		return jsonSerialize(builder.data)
	default:
		return textSerialize(builder.data)
	}
}

// GetContentType returns Content-Type
func (builder ResponseBuilder) GetContentType() string {
	// If content type has been already set, don't change it
	if contentType := builder.ctx.Res.Header().Get("Content-Type"); contentType != "" {
		return contentType
	}
	// Otherwise try to guess proper content type by provided data
	t := reflect.ValueOf(builder.data).Kind()
	switch t {
	case reflect.Map, reflect.Struct, reflect.Array:
		return ContentType["JSON"]
	default:
		return ContentType["text"]
	}
}

func jsonSerialize(data ResponseData) ([]byte, error) {
	json, err := json.Marshal(data)
	if err != nil {
		return []byte(""), err
	}
	return json, nil
}

func textSerialize(data ResponseData) ([]byte, error) {
	return []byte(fmt.Sprintf("%v", data)), nil
}
