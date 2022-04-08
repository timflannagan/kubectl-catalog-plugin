package catalog

import (
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

// waitFor wraps wait.Pool with default polling parameters
func waitFor(fn func() (bool, error)) error {
	return wait.Poll(1*time.Second, 5*time.Minute, func() (bool, error) {
		return fn()
	})
}
