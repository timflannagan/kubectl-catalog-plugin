package delete

import (
	"context"
	"time"

	"github.com/operator-framework/operator-lifecycle-manager/test/e2e"
	"github.com/spf13/cobra"
	"github.com/timflannagan/kubectl-magic-catalog-plugin/cmd/util"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Args:  cobra.ExactArgs(1),
		Short: "Delete an existing FBC magic catalog",
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
			if errors := magicCatalog.UndeployCatalog(ctx); len(errors) != 0 {
				return utilerrors.NewAggregate(errors)
			}

			return nil
		},
	}
	return cmd
}
