package json

import (
	"encoding/json"
	"errors"
)

type obj map[string]any

type resp struct {
	Name               string `form:"name" json:"name" bson:"name"`
	DisplayScrollTypes bool   `form:"displayScrollTypes" json:"displayScrollTypes" bson:"displayScrollTypes"`
	StopTransaction    bool   `form:"stopTransaction" json:"stopTransaction" bson:"stopTransaction"`
}

type respAPI struct {
	resp
	Error string `form:"error" json:"error" bson:"error"`
}

func testRespAPI() ([]byte, error) {
	v := resp{
		DisplayScrollTypes: false,
		StopTransaction:    true,
		Name:               "test name",
	}

	vAPI := respAPI{
		resp:  v,
		Error: errors.New("test error").Error(),
	}

	return json.MarshalIndent(vAPI, "", "    ")
}

func testObj() ([]byte, error) {
	v := obj{
		"displayScrollTypes": false,
		"stopTransaction":    true,
		"name":               "test name",
		"error":              errors.New("test error").Error(),
	}

	return json.MarshalIndent(v, "", "    ")
}
