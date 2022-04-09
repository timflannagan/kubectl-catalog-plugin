package dance

import (
	"context"
	"fmt"
	"time"

	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/spf13/cobra"
	"github.com/timflannagan/kubectl-catalog-plugin/cmd/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dance",
		Args:  cobra.ExactArgs(1),
		Short: "Perform the OLM dance for a namespace that contains an existing magic catalog",
		RunE: func(cmd *cobra.Command, args []string) error {
			namespace := args[0]
			o, err := util.CommonSetup(cmd)
			if err != nil {
				return err
			}
			c := o.Client

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()

			// TODO: handle weird edge case where there's a --namespace flag + a namespace parameter
			// TODO: handle edge case for a partial package uninstall (e.g. no subscription but there's CSVs)
			subs := &operatorsv1alpha1.SubscriptionList{}
			if err := c.List(ctx, subs, client.InNamespace(namespace)); err != nil {
				return err
			}
			if len(subs.Items) == 0 {
				return nil
			}
			if len(subs.Items) > 1 {
				return fmt.Errorf("found more than one subscription in the %s namespace", namespace)
			}
			existingSub := subs.Items[0]

			// TODO: filter out copied CSVs?
			if err := c.DeleteAllOf(ctx, &operatorsv1alpha1.Subscription{}, client.InNamespace(namespace)); err != nil {
				return fmt.Errorf("failed to delete all the subscriptions in the %s namespace: %v", namespace, err)
			}
			if err := c.DeleteAllOf(ctx, &operatorsv1alpha1.ClusterServiceVersion{}, client.InNamespace(namespace)); err != nil {
				return fmt.Errorf("failed to delete all the csvs in the %s namespace: %v", namespace, err)
			}

			newSub := &operatorsv1alpha1.Subscription{
				ObjectMeta: metav1.ObjectMeta{
					Name:      existingSub.GetName(),
					Namespace: existingSub.GetNamespace(),
				},
				Spec: existingSub.Spec,
			}
			if err := c.Create(ctx, newSub); err != nil {
				return fmt.Errorf("failed to recreate the existing %s subscription: %v", newSub.GetName(), err)
			}
			return nil
		},
	}
	return cmd
}
