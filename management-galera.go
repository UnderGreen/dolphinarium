package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli"
)

type action string

// Possible actions for kubectl command
const (
	actionCreate action = "create"
	actionDelete action = "delete"
	actionApply  action = "apply"
)

// Resources path map
var resourcePath = [3]string{
	"galera/k8s/secrets.yaml",
	"galera/k8s/galera-daemonset.yaml",
	"galera/k8s/galera-service.yaml",
}

func kubeCommand(args ...string) *exec.Cmd {
	return exec.Command("/usr/local/bin/kubectl", args...)
}

func fromFile(act action, path string) []byte {
	cmd := kubeCommand(string(act), "-f", path)
	out, _ := cmd.CombinedOutput()
	return out
}

func main() {

	app := cli.NewApp()
	app.Name = "management-galera"
	app.Usage = "Utility for management Galera Cluster inside Kubernetes"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:  "create",
			Usage: "Create Kubernetes resources for Galera Cluster",
			Action: func(c *cli.Context) error {
				for _, path := range resourcePath {
					out := fromFile(actionCreate, path)
					fmt.Print(string(out))
				}
				return nil
			},
		},
		{
			Name:  "delete",
			Usage: "Delete Kubernetes resources for Galera Cluster",
			Action: func(c *cli.Context) error {
				for _, path := range resourcePath {
					out := fromFile(actionDelete, path)
					fmt.Print(string(out))
				}
				return nil
			},
		},
	}

	_ = app.Run(os.Args)
}
