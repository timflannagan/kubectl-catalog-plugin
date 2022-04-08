package util

import (
	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Options struct {
	Client      client.Client
	Namespace   string
	CatalogName string
}

func CommonSetup(cmd *cobra.Command) (*Options, error) {
	namespace, err := cmd.Flags().GetString("namespace")
	if err != nil {
		return nil, err
	}
	catalogName, err := cmd.Flags().GetString("catalog-name")
	if err != nil {
		return nil, err
	}

	scheme := runtime.NewScheme()
	if err := operatorsv1alpha1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := appsv1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := corev1.AddToScheme(scheme); err != nil {
		return nil, err
	}

	config := ctrl.GetConfigOrDie()
	client, err := client.New(config, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, err
	}
	o := &Options{
		Client:      client,
		Namespace:   namespace,
		CatalogName: catalogName,
	}

	return o, nil
}
