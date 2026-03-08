package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type rpcResponse struct {
	Result *json.RawMessage
	Error  *struct {
		Code    int
		Message string
	}
}

func doJSONRPCQuery(url string, method string, result any) error {
	// Create request
	rpcreq := map[string]string{
		"jsonrpc": "2.0",
		"method":  method,
		"id":      "0",
	}
	reqText, err := json.Marshal(rpcreq)
	if err != nil {
		return fmt.Errorf("creating json-rpc query: %w", err)
	}

	// giss an http client
	cli := http.Client{
		Timeout: 100 * time.Millisecond,
	}

	resp, err := cli.Post(url, "text/json", bytes.NewReader(reqText))
	if err != nil {
		return fmt.Errorf("json-rpc post: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("json-rpc readall: %w", err)
	}

	// Parse response
	var rpcResp rpcResponse
	err = json.Unmarshal(respBytes, &rpcResp)
	if err != nil {
		return fmt.Errorf("json-rpc unmarshal response: %w", err)
	}

	// is it an error?
	if rpcResp.Error != nil {
		return fmt.Errorf("json-rpc error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}

	if rpcResp.Result == nil {
		return fmt.Errorf("json-rpc result missing")
	}

	err = json.Unmarshal([]byte(*rpcResp.Result), result)
	if err != nil {
		return fmt.Errorf("json-rpc unmarshal result: %w", err)
	}

	return nil
}
