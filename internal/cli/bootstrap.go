package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/tagesspiegel/helm-plugin-bootstrap/internal/bootstrap"
	"helm.sh/helm/v3/cmd/helm/require"
)

const bootstrapDesc = `
This command modifies an existing chart to add additional files for common use cases like PodDisruptionBudget and NetworkPolicy.

For example, 'helm create foo' was used to create a chart named 'foo':

    foo/
    ├── .helmignore   # Contains patterns to ignore when packaging Helm charts.
    ├── Chart.yaml    # Information about your chart
    ├── values.yaml   # The default values for your templates
    ├── charts/       # Charts that this chart depends on
    └── templates/    # The template files
        └── tests/    # The test files

'helm bootstrap ./foo' takes a path for an argument. If Chart.yaml, values.yaml or templates folder
do not exist, we will return an error. If the given destination exists and there are files in that directory,
conflicting files will be overwritten, but other files will be left alone. The resulting chart will look like this:

	foo/
	├── .helmignore   # Contains patterns to ignore when packaging Helm charts.
	├── Chart.yaml    # Information about your chart
	├── values.yaml   # The default values for your templates
	├── charts/       # Charts that this chart depends on
	├── templates/    # The template files
	│   ├── tests/    # The test files
	│   ├── pdb.yaml  # PodDisruptionBudget file
	│   └── networkpolicy.yaml  # NetworkPolicy file
`

func NewBootstrapCmd(out io.Writer) *cobra.Command {
	flagForce := false
	cmd := &cobra.Command{
		Use:   "bootstrap PATH-TO-CHART",
		Short: "Modify an existing chart to add additional files for common use cases like PodDisruptionBudget and NetworkPolicy",
		Long:  bootstrapDesc,
		Args:  require.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			chartsFolder := args[0]
			fmt.Fprintf(out, "Bootstraping additional files for chart folder %s\n", chartsFolder)
			return bootstrap.Bootstrap(chartsFolder, flagForce)
		},
	}
	flagForce = *cmd.Flags().BoolP("force", "f", false, "Force overwriting existing files")
	return cmd
}
