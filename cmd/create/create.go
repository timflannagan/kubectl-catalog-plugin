package create

import (
	"context"
	"time"

	"github.com/operator-framework/operator-lifecycle-manager/test/e2e"
	"github.com/spf13/cobra"
	"github.com/timflannagan/kubectl-magic-catalog-plugin/cmd/util"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Args:  cobra.ExactArgs(1),
		Short: "Instantiate a file-based catalog (FBC) out of thin air",
		Long: `A kubectl plugin that's responsible for taking an input FBC YAML or JSON file, and creating
a Operator installation using OLM.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fbcPath := args[0]
			provider, err := e2e.NewFileBasedFiledBasedCatalogProvider(fbcPath)
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()

			o, err := util.CommonSetup(cmd)
			if err != nil {
				return err
			}
			magicCatalog := e2e.NewMagicCatalog(o.Client, o.Namespace, o.CatalogName, provider)
			if err := magicCatalog.DeployCatalog(ctx); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
