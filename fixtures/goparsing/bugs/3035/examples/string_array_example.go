//xgo:build goparsing

package examples

// Ntp servers
//
// swagger:response getNtpServersResponse
//
// Example: ["10.10.10.10","20.20.20.20"]
type A []string

// swagger:response getNtpServerList
//
// Example: ["x","y"]
type B struct {
	List []string `json:"list"`
}

// Error from the API.
//
// swagger:response
type Error struct {
	Code int `json:"code"`
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
//	500: Error
func GetNtpServerHandler() {}
