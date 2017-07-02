package commands

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
)

// MixinSpec holds command line flag definitions specific to the mixin
// command. The flags are defined using struct field tags with the
// "github.com/jessevdk/go-flags" format.
type MixinSpec struct {
	ExpectedCollisionCount uint `short:"c" description:"expected # of rejected mixin paths, defs, etc due to existing key. Non-zero exit if does not match actual."`
}

// Execute runs the mixin command which merges Swagger 2.0 specs into
// one spec
//
// Use cases include adding independently versioned metadata APIs to
// application APIs for microservices.
//
// Typically, multiple APIs to the same service instance is not a
// problem for client generation as you can create more than one
// client to the service from the same calling process (one for each
// API).  However, merging clients can improve clarity of client code
// by having a single client to given service vs several.
//
// Server skeleton generation, ie generating the model & marshaling
// code, http server instance etc. from Swagger, becomes easier with a
// merged spec for some tools & target-languages.  Server code
// generation tools that natively support hosting multiple specs in
// one server process will not need this tool.
func (c *MixinSpec) Execute(args []string) error {

	if len(args) < 2 {
		log.Fatalln("Nothing to do. Need some swagger files to merge.\nUSAGE: swagger mixin [-c <expected#Collisions>] <primary-swagger-file> <mixin-swagger-file>...")
	}

	log.Printf("args[0] = %v\n", args[0])
	log.Printf("args[1:] = %v\n", args[1:])
	collisions, err := MixinFiles(args[0], args[1:], os.Stdout)

	for _, warn := range collisions {
		log.Println(warn)
	}

	if err != nil {
		log.Fatalln(err)
	}

	if len(collisions) != int(c.ExpectedCollisionCount) {
		if len(collisions) != 0 {
			// use bash $? to get actual # collisions
			// (but has to be non-zero)
			os.Exit(len(collisions))
		}
		os.Exit(254)
	}
	return nil
}

// MixinFiles is a convenience function for Mixin that reads the given
// swagger files, adds the mixins to primary, calls
// FixEmptyResponseDescriptions on the primary, and writes the primary
// with mixins to the given writer in JSON.  Returns the warning
// messages for collsions that occured during mixin process and any
// error.
func MixinFiles(primaryFile string, mixinFiles []string, w io.Writer) ([]string, error) {

	primaryDoc, err := loads.Spec(primaryFile)
	if err != nil {
		return nil, err
	}
	primary := primaryDoc.Spec()

	var mixins []*spec.Swagger
	for _, mixinFile := range mixinFiles {
		mixin, err := loads.Spec(mixinFile)
		if err != nil {
			return nil, err
		}
		mixins = append(mixins, mixin.Spec())
	}

	collisions := analysis.Mixin(primary, mixins...)
	analysis.FixEmptyResponseDescriptions(primary)

	bs, err := json.MarshalIndent(primary, "", "  ")
	if err != nil {
		return collisions, err
	}

	_, _ = w.Write(bs)

	return collisions, nil
}
