package main

import (
	"fmt"
	"os"

	"github.com/timflannagan/kubectl-catalog-plugin/cmd/root"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Options struct {
	namespace   string
	catalogName string
	client      client.Client
}

func main() {
	cmd := root.NewCmd()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
