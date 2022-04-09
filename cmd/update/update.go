package update

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/timflannagan/kubectl-catalog-plugin/cmd/util"
	catalog "github.com/timflannagan/kubectl-catalog-plugin/internal"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Args:  cobra.ExactArgs(1),
		Short: "Update an existing FBC magic catalog",
		RunE: func(cmd *cobra.Command, args []string) error {
			fbcPath := args[0]
			provider, err := catalog.NewFileBasedFiledBasedCatalogProvider(fbcPath)
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()

			o, err := util.CommonSetup(cmd)
			if err != nil {
				return err
			}
			magicCatalog := catalog.NewMagicCatalog(o.Client, o.Namespace, o.CatalogName, provider)
			if err := magicCatalog.UpdateCatalog(ctx, provider); err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}
