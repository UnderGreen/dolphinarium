package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"

	"github.com/Jeffail/gabs"
	"github.com/urfave/cli"
)

// PodsRequired - required count of worker pods for Galera Cluster
var PodsRequired = 3

type action string

// Possible actions for kubectl command
const (
	actionCreate action = "create"
	actionDelete action = "delete"
	actionGet    action = "get"
)

type resourcetype string

// Possible resource types for kubectl get command
const (
	resourcePod       resourcetype = "pod"
	resourceDaemonset resourcetype = "daemonset"
	resourceService   resourcetype = "service"
	resourceSecret    resourcetype = "secret"
)

// Resources path map
var resourcePath = map[resourcetype]string{
	resourceSecret:    "galera/k8s/secrets.yaml",
	resourceDaemonset: "galera/k8s/galera-daemonset.yaml",
	resourceService:   "galera/k8s/galera-service.yaml",
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

func fromCmd(act action, resource resourcetype, jsonOutput bool, label string) []byte {
	var cmd *exec.Cmd
	if jsonOutput {
		cmd = kubeCommand(string(act), "-o", "json", string(resource), "-l", label)
	} else {
		cmd = kubeCommand(string(act), string(resource), "-l", label)
	}
	out, _ := cmd.CombinedOutput()
	return out
}

func getInfoJSON(resourcePath string, resource resourcetype, jsonPath string, label string) interface{} {
	var (
		jsonParsed *gabs.Container
		err        error
	)
	if len(resourcePath) > 0 {
		jsonParsed, err = gabs.ParseJSON(fromFile(actionGet, resourcePath, true))
	} else {
		jsonParsed, err = gabs.ParseJSON(fromCmd(actionGet, resource, true, label))
	}
	if err != nil {
		return nil
	}
	value := jsonParsed.Path(jsonPath).Data()

	return value
}

func checkDaemonSet() error {
	currentNumberScheduled := getInfoJSON(resourcePath[resourceDaemonset], "", "status.currentNumberScheduled", "")
	desiredNumberScheduled := getInfoJSON(resourcePath[resourceDaemonset], "", "status.desiredNumberScheduled", "")

	if currentNumberScheduled == nil || desiredNumberScheduled == nil {
		return errors.New("===> DaemonSet has't a valid state. Check JSON output from kubectl")
	}

	if currentNumberScheduled != desiredNumberScheduled || currentNumberScheduled != float64(PodsRequired) {
		return errors.New("===> DaemonSet has't a valid state. Number of current nodes not equal desired")
	}

	return nil
}

func checkPods(label string) error {
	jsonParsed, err := gabs.ParseJSON(fromCmd(actionGet, resourcePod, true, "name="+label))
	if err != nil {
		return nil
	}

	runningStatuses := reflect.ValueOf(jsonParsed.Path("items.status.phase").Data())

	for i := 0; i < runningStatuses.Len(); i++ {
		podStatus := runningStatuses.Index(i)
		if podStatus.Interface().(string) != "Running" {
			return errors.New("One or more pods is not running")
		}
	}

	readyStatuses := jsonParsed.Path("items.status.containerStatuses.ready").Bytes()
	rs, _ := gabs.ParseJSON(readyStatuses)
	for i := 0; i < PodsRequired; i++ {
		readiness := reflect.ValueOf(rs.Index(i).Index(0).String())
		if readiness.Interface().(string) != "true" {
			return errors.New("One or more pods is not ready")
		}
	}

	return nil
}

func checkResources(c *cli.Context) (err error) {
	err = checkDaemonSet()
	if err != nil {
		fmt.Fprintln(c.App.Writer, err)
	} else {
		fmt.Fprintln(c.App.Writer, "DaemonSet has a valid state")
	}

	label := getInfoJSON(resourcePath[resourceDaemonset], "", "metadata.labels.name", "")
	if label == nil {
		fmt.Fprintln(c.App.Writer, "===> Cannot determine label from DaemonSet output")
	}

	err = checkPods(label.(string))
	if err != nil {
		fmt.Fprintln(c.App.Writer, err)
	} else {
		fmt.Fprintln(c.App.Writer, "All pods have a valid state")
	}

	return nil
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
			Name:   "validate",
			Usage:  "Validate readiness of pods for all labeled nodes and readiness of service",
			Action: checkResources,
		},
	}

	_ = app.Run(os.Args)
}
