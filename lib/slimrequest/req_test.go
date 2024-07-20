package slimrequest

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestNewRequest(t *testing.T) {
	expected := Request{
		ID:     1,
		Method: "slim.request",
		Params: []any{"0", []string{"serverstatus", "0"}},
	}
	expectedJSON := `{"id":1,"method":"slim.request","params":["0",["serverstatus","0"]]}`

	made := NewRequest("0", []string{"serverstatus", "0"})

	if !reflect.DeepEqual(expected, made) {
		t.Errorf("expected %+v but got %+v", expected, made)
	}

	bs, err := json.Marshal(made)
	if err != nil {
		t.Errorf("JSON marshalling failed: %v", err)
	}
	s := string(bs)

	if s != expectedJSON {
		t.Errorf("expected %+v but got %+v", expectedJSON, s)
	}
}
