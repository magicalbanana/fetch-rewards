package handler

import (
	"encoding/json"
)

func NewJSONResponse() *Response {
	r := Response{
		Body: &ResponseBody{Errors: make(Errors, 0)},
	}
	return &r
}

// Response ...
type Response struct {
	// AppCode is the code that will be used on the ResponseBody.Code
	AppCode string
	// HTTPStatus is the HTTP Status Code that will be used in the
	// response
	HTTPStatus int
	// Body is the response body that will be marshaled in to a []byte
	Body *ResponseBody
}

// Marshal ...
func (r *Response) Marshal() ([]byte, error) {
	appCode := "app.undefined_app_code"
	if r.AppCode != "" {
		appCode = r.AppCode
	}
	r.Body.Code = appCode
	b, err := json.Marshal(r.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// AddErrors adds Errors
func (r *Response) AddErrors(ee Errors) {
	if r.Body.Errors == nil {
		r.Body.Errors = make(Errors, 0)
	}

	for _, e := range ee {
		if r.AppCode != "" && e.Code == "" {
			e.Code = r.AppCode
		}
	}
	r.Body.Errors = append(r.Body.Errors, ee...)
}

// AddError adds a single error
func (r *Response) AddError(e *Error) {
	if r.Body.Errors == nil {
		r.Body.Errors = make(Errors, 0)
	}
	if e.Code == "" && r.AppCode != "" {
		e.Code = r.AppCode
	}
	r.Body.Errors = append(r.Body.Errors, *e)
}

// ResponseBody ...
type ResponseBody struct {
	// Code is the application specific code
	Code string `json:"code"`
	// Data is the data of the resposne body (main response)
	Data *json.RawMessage `json:"data,omitempty"`
	// Errors is the errors that have been collected during the request
	Errors []Error `json:"errors,omitempty"`
}
