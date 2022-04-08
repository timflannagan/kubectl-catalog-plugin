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

type options struct {
	namespace   string
	catalogName string
	client      client.Client
}

func main() {
	o := &options{}

	rootCmd := &cobra.Command{
		Use: "catalog",
	}
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

			if err := o.commonSetup(); err != nil {
				return err
			}
			if err := o.run(ctx, provider); err != nil {
				return err
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&o.namespace, "namespace", "default", "Configures the namespace to find the Bundle underlying resources")
	cmd.Flags().StringVar(&o.catalogName, "catalog-name", "magiccatalog", "Configures the metadata.Name for the generated ConfigMap resource")

	rootCmd.AddCommand(cmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (o *options) commonSetup() error {
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
	o.client = client

	return nil
}

func (o *options) run(ctx context.Context, provider e2e.FileBasedCatalogProvider) error {
	magicCatalog := e2e.NewMagicCatalog(o.client, o.namespace, o.catalogName, provider)
	if err := magicCatalog.DeployCatalog(ctx); err != nil {
		return err
	}
	return nil
}
