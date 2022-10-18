package unixsock

import (
	"context"
	"time"
	podresourcesapi "k8s.io/kubelet/pkg/apis/podresources/v1"
	"k8s.io/kubernetes/pkg/kubelet/apis/podresources"
	"k8s.io/kubernetes/pkg/kubelet/util"
)

const (
	kubeletPodResourcePath = "/var/lib/kubelet/pod-resources"
	connectionTimeout = 10 * time.Second
	maxMsgSize = 16 * 1024 * 1024
)

var (
	kubeletSocket string
)

func init() {
	// func LocalEndpoint(path, file string) (string, error)
	// returns the full path to a unix socket
	var err error
	kubeletSocket, err = util.LocalEndpoint(kubeletPodResourcePath, podresources.Socket)
	if err != nil {
		panic(err)
	}
}

func GetPodResources() ([]*podresourcesapi.PodResources, error) {
	// func GetV1Client(socket string, connectionTimeout time.Duration, maxMsgSize int) (v1.PodResourcesListerClient, *grpc.ClientConn, error)
	client, conn, err := podresources.GetV1Client(kubeletSocket, connectionTimeout, maxMsgSize)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	// List(ctx context.Context, in *ListPodResourcesRequest, opts ...grpc.CallOption) (*ListPodResourcesResponse, error)
	resp, err := client.List(ctx, &podresourcesapi.ListPodResourcesRequest{})
	if err != nil {
		return nil, err
	}

	return resp.PodResources, nil
}
