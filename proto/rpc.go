package proto

import "encoding/json"

// Request object
type Request struct {
	Version string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      int             `json:"id"`
}

// Marshall returns the json encoding of Request
func (s Request) Marshal() json.RawMessage {
	b, _ := json.Marshal(s)
	return json.RawMessage(b)
}

// Response object
type Response struct {
	Result json.RawMessage `json:"result"`
	Error  *Error          `json:"error"`
}

// Error object
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
