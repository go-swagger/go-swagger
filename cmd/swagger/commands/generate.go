package commands

import "github.com/casualjim/go-swagger/cmd/swagger/commands/generate"

// Generate command to group all generator commands together
type Generate struct {
	Model     *generate.Model     `command:"model"`
	Operation *generate.Operation `command:"operation"`
	Support   *generate.Support   `command:"support"`
	Server    *generate.Server    `command:"server"`
	Spec      *generate.SpecFile  `command:"spec"`
}
