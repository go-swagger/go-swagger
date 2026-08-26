// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

// Command swagger is a CLI tool to work with Swagger specifications.
package main

import (
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"slices"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

type docCommand struct {
	Destination string `default:"./docs" description:"Output destination folder" long:"dest" short:"d"`
	parser      *flags.Parser
}

func (d *docCommand) Execute(_ []string) error {
	return d.gendoc()
}

func (d *docCommand) gendoc() (err error) {
	if d.Destination == "" {
		const defaultFolder = "docs"
		d.Destination = defaultFolder
	}

	const writableMode = 0o755
	if err = os.MkdirAll(d.Destination, writableMode); err != nil {
		return err
	}

	for documented := range documentCLI(d.parser) {
		file, err := os.Create(filepath.Join(d.Destination, documented.Target))
		if err != nil {
			return err
		}

		// front matter
		fmt.Fprintln(file, "---")
		fmt.Fprintf(file, "title: %q\n", documented.Title)
		fmt.Fprintf(file, "description: %q\n", documented.Description)
		fmt.Fprintf(file, "weight: %d\n", documented.Index)
		fmt.Fprintln(file, "---")
		fmt.Fprintln(file, "")

		// markdown formatting
		fmt.Fprintf(file, "## %s\n", documented.Title)
		fmt.Fprintln(file, "")
		fmt.Fprintln(file, "```cmd")

		// go-flags formats a help message by polling the width of the stdin terminal: we give it a pseudo-tty, sized at 132 cols.
		const desiredTermWidth = 132
		restore, err := setTermsize(desiredTermWidth)
		if err != nil {
			return err
		}
		defer restore()

		// go-flags writes help messages to stdout: redirecting here
		stdout := os.Stdout
		stderr := os.Stderr
		defer func() {
			os.Stdout = stdout
			os.Stderr = stderr
		}()

		os.Stdout = file
		os.Stderr = file

		// run the command with --help flag; swallow the error systematically returned by go-flag on --help
		_ = run(d.parser, documented.Args)

		fmt.Fprintln(file, "```")
		err = file.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

type doc struct {
	Title       string
	Description string
	Target      string
	Args        []string
	Index       int // the index in the tree of commands: used to weight the produced pages
}

// documentCLI yields an iterator over the tree of commands and their subcommands.
func documentCLI(parser *flags.Parser) iter.Seq[doc] {
	return func(yield func(doc) bool) {
		i := 0
		rootDoc := doc{
			Title:       "Commands",
			Description: "All swagger commands",
			Index:       i,
			Target:      "commands.md",
			Args:        []string{"--help"},
		}
		if !yield(rootDoc) {
			return
		}

		i++
		for cmdDoc := range documentCommands([]string{}, &i, parser.Commands()) {
			if !yield(cmdDoc) {
				return
			}
			i++
		}
	}
}

func documentCommands(parents []string, index *int, commands []*flags.Command) iter.Seq[doc] {
	var prefix, parent string
	if len(parents) > 0 {
		prefix = strings.Join(parents, "_") + "_"
		parent = strings.Join(parents, " ") + " "
	}

	return func(yield func(doc) bool) {
		for _, cmd := range commands {
			const addedArgs = 2
			args := make([]string, 0, len(parents)+addedArgs)
			args = append(args, parents...)
			args = append(args, cmd.Name, "--help")

			cmdDoc := doc{
				Title:       parent + cmd.Name,
				Description: cmd.ShortDescription,
				Target:      prefix + cmd.Name + ".md",
				Args:        args,
				Index:       *index,
			}
			if !yield(cmdDoc) {
				return
			}
			*index++

			children := slices.Clone(parents)
			children = slices.Grow(children, 1)
			children = append(children, cmd.Name)
			for subCmd := range documentCommands(children, index, cmd.Commands()) {
				if !yield(subCmd) {
					return
				}
				*index++
			}
		}
	}
}
