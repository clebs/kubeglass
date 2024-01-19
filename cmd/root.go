package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/clebs/kubeglass/pkg/api"
	"github.com/clebs/kubeglass/pkg/api/json"
	"github.com/clebs/kubeglass/pkg/fetch"
	"github.com/clebs/kubeglass/pkg/format"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

var (
	from string
	to   string
	out  string
)

func New() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          filepath.Base(os.Args[0]),
		SilenceUsage: true,
		Short:        "Compare the APIs of 2 kubernetes versions.",
		Long:         `Given 2 kubernetes versions, compares the APIs of both showing all group, resource and version changes.`,
		Example:      "kubeglass -f 1.28 -t 1.29",
		Args:         cobra.MinimumNArgs(0),
		PreRunE:      validate,
		RunE:         run,
	}

	rootCmd.PersistentFlags().StringVarP(&from, "from", "f", "", "Sets the starting kuberntes version for the comparison")
	rootCmd.PersistentFlags().StringVarP(&to, "to", "t", "", "Sets the target kuberntes version for the comparison")
	rootCmd.PersistentFlags().StringVarP(&out, "out", "o", "", "Sets the output format (stdout, json, yaml)")

	return rootCmd
}

func validate(_ *cobra.Command, _ []string) error {
	if !strings.HasPrefix(from, "v") {
		from = fmt.Sprintf("%c%s", 'v', from)
	}
	if !strings.HasPrefix(to, "v") {
		to = fmt.Sprintf("%c%s", 'v', to)
	}

	if !semver.IsValid(from) {
		return fmt.Errorf("the given from version (%s) is not a valid semantic version", from)
	}
	if !semver.IsValid(to) {
		return fmt.Errorf("the given to version (%s) is not a valid semantic version", to)
	}

	if semver.Compare(from, "v1.28") == -1 {
		return fmt.Errorf("kubeglass only supports kubernetes versions v1.28+")
	}

	if semver.Compare(from, to) >= 0 {
		return fmt.Errorf("the given from version (%s) must be lower than the to version (%s)", from, to)
	}

	return nil
}

func run(_ *cobra.Command, _ []string) error {
	fromJSON, err := fetch.Version(from)
	if err != nil {
		return err
	}

	toJSON, err := fetch.Version(to)
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

	output, err := format.NewFormatter(out).Format(changes)
	if err != nil {
		return err
	}

	fmt.Print(string(output))
	return nil
}
