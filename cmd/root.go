package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/clebs/kubeglass/pkg/api"
	"github.com/clebs/kubeglass/pkg/api/json"
	"github.com/clebs/kubeglass/pkg/fetch"
	"github.com/clebs/kubeglass/pkg/format"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
	"sigs.k8s.io/release-utils/version"
)

const (
	sourceFile source = iota
	sourceSemVer
)

type source uint8

type root struct {
	source source
	from   string
	to     string
	out    string
}

func New() *cobra.Command {
	cmd := root{}

	cobraCmd := &cobra.Command{
		Use:          "kubeglass",
		SilenceUsage: true,
		Short:        "Compare the APIs of 2 kubernetes versions.",
		Long:         `Given 2 kubernetes versions, compares the APIs of both showing all group, resource and version changes.`,
		Example:      "kubeglass -f 1.28 -t 1.29",
		Args:         cobra.MinimumNArgs(0),
		PreRunE:      cmd.validate,
		RunE:         cmd.run,
	}

	cobraCmd.PersistentFlags().StringVarP(&cmd.from, "from", "f", "", "Sets the starting kuberntes version for the comparison")
	cobraCmd.PersistentFlags().StringVarP(&cmd.to, "to", "t", "", "Sets the target kuberntes version for the comparison")
	cobraCmd.PersistentFlags().StringVarP(&cmd.out, "out", "o", "", "Sets the output format (stdout, json, yaml)")

	cobraCmd.AddCommand(version.WithFont("starwars"))

	return cobraCmd
}

func (r *root) validate(cmd *cobra.Command, _ []string) error {
	if r.from == "" && r.to == "" {
		if err := cmd.Help(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	// r.from flag
	if stats, err := os.Stat(r.from); err == nil {
		if stats.IsDir() {
			// not a file
			return fmt.Errorf("the given r.from source (%s) is not a valid local file", r.from)
		}
		r.source = sourceFile
	} else {
		// not a path, maybe a semver
		if !strings.HasPrefix(r.from, "v") {
			r.from = fmt.Sprintf("%c%s", 'v', r.from)
		}

		if !semver.IsValid(r.from) {
			return fmt.Errorf("the given r.from source (%s) is not a valid semantic version", r.from)
		}

		if semver.Compare(r.from, "v1.28") == -1 {
			return fmt.Errorf("kubeglass only supports kubernetes versions v1.28+")
		}
		r.source = sourceSemVer
	}

	// r.to flag
	if stats, err := os.Stat(r.to); err == nil {
		if stats.IsDir() {
			// not a file
			return fmt.Errorf("the given r.to source (%s) is not a valid local file", r.to)
		}
	} else {
		// not a path, maybe a semver
		if !strings.HasPrefix(r.to, "v") {
			r.to = fmt.Sprintf("%c%s", 'v', r.to)
		}

		if !semver.IsValid(r.to) {
			return fmt.Errorf("the given r.to source (%s) is not a valid semantic version", r.to)
		}

		if semver.Compare(r.from, r.to) >= 0 {
			return fmt.Errorf("the given r.from version (%s) must be lower than the r.to version (%s)", r.from, r.to)
		}
	}

	return nil
}

func (r *root) run(_ *cobra.Command, _ []string) error {
	var fromJSON []byte
	var err error
	switch r.source {
	case sourceFile:
		fromJSON, err = fetch.File(r.from)
	case sourceSemVer:
		fromJSON, err = fetch.SemVer(r.from)
	}
	if err != nil {
		return err
	}

	var toJSON []byte
	switch r.source {
	case sourceFile:
		toJSON, err = fetch.File(r.to)
	case sourceSemVer:
		toJSON, err = fetch.SemVer(r.to)
	}
	if err != nil {
		return err
	}

	fromMap, err := json.ToAPIMap(fromJSON)
	if err != nil {
		return err
	}

	toMap, err := json.ToAPIMap(toJSON)
	if err != nil {
		return err
	}

	changes, err := api.Report(fromMap, toMap)
	if err != nil {
		return err
	}

	output, err := format.NewFormatter(r.out).Format(changes)
	if err != nil {
		return err
	}

	fmt.Print(string(output))
	return nil
}
