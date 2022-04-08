package main

import (
	"context"
	"fmt"
	"os"
	"time"

	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/operator-framework/operator-lifecycle-manager/test/e2e"
	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	var (
		ns          string
		catalogName string
	)
	cmd := &cobra.Command{
		Use:   "evaluate",
		Args:  cobra.ExactArgs(1),
		Short: "Instantiate a file-based catalog (FBC) out of thin air",
		Long: `
A kubectl plugin that's responsible for taking an input FBC YAML or JSON file, and creating
a Operator installation using OLM.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fbcPath := args[0]
			provider, err := e2e.NewFileBasedFiledBasedCatalogProvider(fbcPath)
			if err != nil {
				return err
			}

			scheme := runtime.NewScheme()
			if err := operatorsv1alpha1.AddToScheme(scheme); err != nil {
				return err
			}
			if err := appsv1.AddToScheme(scheme); err != nil {
				return err
			}
			if err := corev1.AddToScheme(scheme); err != nil {
				return err
			}

			config := ctrl.GetConfigOrDie()
			client, err := client.New(config, client.Options{
				Scheme: scheme,
			})
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()

			magicCatalog := e2e.NewMagicCatalog(client, ns, catalogName, provider)
			if err := magicCatalog.DeployCatalog(ctx); err != nil {
				return err
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&ns, "namespace", "default", "Configures the namespace to find the Bundle underlying resources")
	cmd.Flags().StringVar(&catalogName, "catalog-name", "magiccatalog", "Configures the metadata.Name for the generated ConfigMap resource")

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
