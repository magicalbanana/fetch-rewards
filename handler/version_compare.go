package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/magicalbanana/fetch-rewards/semver"
	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/net/webdav"

	uuid "github.com/satori/go.uuid"
)

type VersionCompare struct {
	Data struct {
		CompareFrom string `json:"compare_from,omitempty"`
		CompareTo   string `json:"compare_to,omitempty"`
	} `json:"data"`
}

func (v *VersionCompare) Compare(r *http.Request) (*Response, error) {
	// validate JSON schema
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	defer func() {
		if r.Body != nil {
			r.Body.Close()
		}
	}()

	result, err := gojsonschema.Validate(
		gojsonschema.NewBytesLoader([]byte(compareVersionSchema)),
		gojsonschema.NewBytesLoader(reqBody),
	)
	if err != nil {
		return nil, err
	}

	resp := NewJSONResponse()

	// check if the JSON schema validation returned validation errors
	if !result.Valid() {
		var ee []Error
		for _, resErr := range result.Errors() {
			var property string
			v, ok := resErr.Details()["property"]
			if ok || v == nil {
				property = fmt.Sprintf("%v", v)
			}
			e := Error{
				ID:     uuid.NewV4().String(),
				Status: http.StatusText(webdav.StatusUnprocessableEntity),
				Code:   "app.invalid_payload",
				Title:  "Invalid Request Payload",
				Detail: resErr.Description(),
				Source: &ErrorSource{
					Pointer: resErr.Context().String(),
					Value:   fmt.Sprintf("%v", resErr.Details()[property]),
				},
			}
			ee = append(ee, e)
		}

		resp.HTTPStatus = http.StatusBadRequest
		resp.AddErrors(ee)

		return resp, nil
	}

	err = json.Unmarshal(reqBody, &v)
	if err != nil {
		return nil, err
	}

	compareFrom, err := semver.NewVersion(v.Data.CompareFrom)
	if err != nil {
		return nil, err
	}
	compareTo, err := semver.NewVersion(v.Data.CompareTo)
	if err != nil {
		return nil, err
	}

	var msg string
	switch compareFrom.Compare(compareTo) {
	case semver.LessThan:
		msg = fmt.Sprintf(`%v is "before" %v`, v.Data.CompareFrom, v.Data.CompareTo)
	case semver.Equal:
		msg = fmt.Sprintf(`%v is "equal" to %v`, v.Data.CompareFrom, v.Data.CompareTo)
	case semver.GreaterThan:
		msg = fmt.Sprintf(`%v is "after" %v`, v.Data.CompareFrom, v.Data.CompareTo)
	}

	b, err := json.Marshal(map[string]string{
		"result": msg,
	})
	if err != nil {
		return nil, err
	}
	resp.AppCode = "app.compare_success"
	resp.HTTPStatus = http.StatusOK
	rm := json.RawMessage(b)
	resp.Body.Data = &rm
	return resp, nil
}
