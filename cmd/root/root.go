package root

import (
	"github.com/spf13/cobra"
	"github.com/timflannagan/kubectl-magic-catalog-plugin/cmd/create"
	"github.com/timflannagan/kubectl-magic-catalog-plugin/cmd/delete"
)

func NewCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "catalog",
		Short: `A kubectl plugin for creating and managing FBC catalogs in dev environments`,
		Args:  cobra.ExactArgs(0),
	}
	rootCmd.AddCommand(create.NewCmd(), delete.NewCmd())

	rootCmd.PersistentFlags().String("namespace", "default", "Configures the namespace to find the Bundle underlying resources")
	rootCmd.PersistentFlags().String("catalog-name", "magiccatalog", "Configures the metadata.Name for the generated ConfigMap resource")

	return rootCmd
}
