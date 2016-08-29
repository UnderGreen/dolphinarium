package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/Jeffail/gabs"
	"github.com/urfave/cli"
)

type action string

// Possible actions for kubectl command
const (
	actionCreate action = "create"
	actionDelete action = "delete"
	actionApply  action = "apply"
	actionGet    action = "get"
)

type resourcetype string

// Possible resource types for kubectl get command
const (
	resourcePod       resourcetype = "pod"
	resourceDaemonset resourcetype = "daemonset"
	resourceService   resourcetype = "service"
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

func fromFile(act action, path string, jsonOutput bool) []byte {
	var cmd *exec.Cmd
	if jsonOutput {
		cmd = kubeCommand(string(act), "-f", path, "-o", "json")
	} else {
		cmd = kubeCommand(string(act), "-f", path)
	}
	out, _ := cmd.CombinedOutput()
	return out
}

func getInfoJSON(resourcePath string, jsonPath string) interface{} {
	jsonParsed, err := gabs.ParseJSON(fromFile(actionGet, resourcePath, true))
	if err != nil {
		return nil
	}
	value := jsonParsed.Path(jsonPath).Data()

	return value
}

func checkDaemonSet() (err error) {
	currentNumberScheduled, ok := getInfoJSON("galera/k8s/galera-daemonset.yaml", "status.currentNumberScheduled").(float64)
	if !ok {
		return errors.New("DaemonSet has't a valid state. Check the JSON output.")
	}
	desiredNumberScheduled, ok := getInfoJSON("galera/k8s/galera-daemonset.yaml", "status.desiredNumberScheduled").(float64)
	if !ok {
		return errors.New("DaemonSet has't a valid state. Check the JSON output.")
	}

	if currentNumberScheduled != desiredNumberScheduled {
		return errors.New("DaemonSet has't a valid state. Number of current nodes not equal desired.")
	}

	return err
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
					out := fromFile(actionCreate, path, false)
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
					out := fromFile(actionDelete, path, false)
					fmt.Print(string(out))
				}
				return nil
			},
		},
		{
			Name:  "validate",
			Usage: "Validate readiness of pods for all labeled nodes and readiness of service",
			Action: func(c *cli.Context) error {
				err := checkDaemonSet()
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("DaemonSet has a valid state.")
				}
				return nil
			},
		},
	}

	_ = app.Run(os.Args)
}
