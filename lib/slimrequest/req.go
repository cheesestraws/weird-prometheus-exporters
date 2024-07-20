package slimrequest

import (
	_ "encoding/json"
)

// The same type is used for requests and responses: only difference
// is requests don't have a results field.
type req[A any] struct {
	ID     int    `json:"id"`
	Method string `json:"method"`
	Params []any  `json:"params"`
	Result A      `json:"result,omitempty" `
}

type Request req[*struct{}]

func NewRequest(playerID string, params []string) Request {
	return Request{
		ID:     1,
		Method: "slim.request",
		Params: []any{
			playerID,
			params,
		},
		Result: nil,
	}
}
