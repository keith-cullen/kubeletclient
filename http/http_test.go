package http

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

const (
	port = "10250"
	caPath = "/var/lib/kubelet/pki/kubelet.crt"
	nodePath = "./node.txt"
	tokenPath = "./token.txt"
)

var (
	node string
	token string
)

func init() {
	b, err := os.ReadFile(nodePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Missing node file: %s\n", nodePath)
		os.Exit(1)
	}
	node = strings.Trim(string(b), " \t\r\n")

	b, err = os.ReadFile(tokenPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Missing token file: %s\n", tokenPath);
		fmt.Fprintf(os.Stderr, "To take a service account token from an existing container:\n")
		fmt.Fprintf(os.Stderr, "    kubectl exec -n <namespace> <podname> -c <containername> -- cat /var/run/secrets/kubernetes.io/serviceaccount/token > token.txt\n")
		os.Exit(1)
	}
	token = strings.Trim(string(b), " \t\r\n")
}

func TestHttpGetPodList(t *testing.T) {
	client, err := NewClient(node, port, caPath, token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	podList, err := client.GetPodList()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	for _, pod := range podList.Items {
		fmt.Printf("pod.Spec.NodeName: %s, pod.Status.PodIP: %s\n", pod.Spec.NodeName, pod.Status.PodIP)
	}
}

func TestHttpGetMetrics(t *testing.T) {
	client, err := NewClient(node, port, caPath, token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	metrics, err := client.GetStr("/metrics")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", metrics)
}
