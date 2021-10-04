package cli

import (
	"github.com/robinfoe/kuberparse/internal/task"
	"github.com/spf13/cobra"
)

// NewTestserverCommand ...
func ParseCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "parse [OPTIONS]",
		Short: "parse openshift construct",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {

			dc, _ := cmd.Flags().GetString("path-dc")
			cm, _ := cmd.Flags().GetString("path-cm")
			fo, _ := cmd.Flags().GetString("path-csv")

			opts := &task.ParseOptions{

				Path: &task.PathOptions{
					DeploymentConfig: dc,
					ConfigMap:        cm,
					FileOutput:       fo,
				},
			}

			parser := &task.Parser{
				Opts: opts,
			}

			return parser.Parse()
		},
	}

	cmd.Flags().String("path-dc", "/home/robin/workspace/source/golang/kubeparse/asset/dc-trim.yaml", "Path to deployment config")
	cmd.Flags().String("path-cm", "/home/robin/workspace/source/golang/kubeparse/asset/cm.yaml", "Path to configmap")
	cmd.Flags().String("path-csv", "/home/robin/workspace/source/golang/kubeparse/asset/output.csv", "Path to output csv")

	// cmd.Flags().String("path-dc", "dc.yaml", "Path to deployment config")
	// cmd.Flags().String("path-cm", "cm.yaml", "Path to configmap")
	// cmd.Flags().String("path-csv", "result.csv", "Path to output csv")

	return &cmd
}
