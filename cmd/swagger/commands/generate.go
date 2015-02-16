package commands

import "github.com/casualjim/go-swagger/cmd/swagger/commands/generate"

// Generate command to group all generator commands together
type Generate struct {
	Model     *generate.Model     `command:"model"`
	Operation *generate.Operation `command:"operation"`
	// All       *generate.All       `command:"all"`
}
