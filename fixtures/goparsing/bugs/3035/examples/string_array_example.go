//xgo:build goparsing

package examples

// Ntp servers
//
// swagger:response getNtpServersResponse
//
// In: body
// Example: ["10.10.10.10","20.20.20.20"]
type A []string

// B yields a list of ntp servers.
//
// swagger:response getNtpServerList
//
// Example: ["x","y"]
type B struct {
	// In: body
	//
	// Example: ["a","b"]
	List []string `json:"list"`
}

// Error from the API.
type Error struct {
	// swagger:response
	// In: body
	Code uint64 `json:"code"`
}

type response struct {
	// swagger:response
	// In: body
	Body A
	// swagger:response
	// In: headers
	Headers B
}

// swagger:route GET /ntp-server opID
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Responses:
//
//	200: getNtpServersResponse
//	201: getNtpServerList
//	203: response
//	500: Error
func GetNtpServerHandler() {}
