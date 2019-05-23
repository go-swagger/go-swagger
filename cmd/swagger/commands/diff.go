package commands

import (
	"log"
	"errors"

	"github.com/go-openapi/loads"
	"github.com/go-swagger/go-swagger/cmd/swagger/commands/diff"
)

// DiffCommand is a command that generates the diff of two swagger specs.
//
// There are no specific options for this expansion.
type DiffCommand struct {
	OnlyBreakingChanges bool `long:"break" short:"b" description:"When present, only shows incompatible changes"`
	Format          string `long:"format" short:"f" description:"When present, writes output as json" default:"txt" choice:"txt" choice:"json"`
	IgnoreFile string `long:"ignore" short:"i" description:"Exception file of diffs to ignore"  default:"none specified"`
	Destination string `long:"dest" short:"d" description:"Output destination file or stdout" default:"stdout"`
}

// Execute diffs the two specs provided
func (c *DiffCommand) Execute(args []string) error {
	if len(args) != 2 {
		msg := `The diff command requires an old and new swagger document urls to be specified.
		diff will fail if any breaking changes are detected.
		eg 
		   swagger diff spec1Path spec2Path
		
		Breaking changes are:
		  * all path and verb combinations in the old specification are present in the new one
		  * no request parameters are required in the new specification that were not required in the old one
		  * all request parameters in the old specification are present in the new one
		  * all request parameters in the old specification have the same type in the new one
		  * all response attributes in the old specification are present in the new one
		  * all response attributes in the old specification have the same type in the new one
		  * all enum values accepted by the old specification are accepted by the new one
`
		return errors.New(msg)
	}

	log.Println("Run Config:")
	log.Printf("Spec1: %s",args[0])
	log.Printf("Spec2: %s",args[1])
	log.Printf("ReportOnlyBreakingChanges (-c) :%v",c.OnlyBreakingChanges)
	log.Printf("OutputFormat (-f) :%s",c.Format)
	log.Printf("IgnoreFile (-i) :%s",c.IgnoreFile)
	log.Printf("Diff Report Destination (-d) :%s",c.Destination)

	diffs, err := getDiffs(args[0], args[1])

	if err != nil {
		return err
	}

	if c.Format=="json"{
		diffs.ReportAllDiffs(true)

	}else{
		if c.OnlyBreakingChanges {
			err = diffs.ReportCompatability()
		} else {
			diffs.ReportAllDiffs(false)
		}
	}

	return err
}

func getIgnores()[]diff.SpecDifference{
	ignores := []diff.SpecDifference{}
	return ignores
}

func getDiffs(oldSpecPath, newSpecPath string) (*diff.SpecDiff, error) {
	swaggerDoc1 := oldSpecPath
	specDoc1, err := loads.Spec(swaggerDoc1)

	if err != nil {
		return nil, err
	}

	swaggerDoc2 := newSpecPath
	specDoc2, err := loads.Spec(swaggerDoc2)
	if err != nil {
		return nil, err
	}

	diffs := diff.Compare(specDoc1.Spec(), specDoc2.Spec())
	return diffs, nil
}
