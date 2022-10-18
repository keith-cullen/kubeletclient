package unixsock

import (
	"fmt"
	"testing"
)

func TestUnixsockGetPodResources(t *testing.T) {
	resources, err := GetPodResources()
	if err != nil {
		t.Errorf("%v\n", err)
	}
	for i, pod := range resources {
		fmt.Printf("pod[%d]: %s\n", i, pod.Name)
		for j, cont := range pod.Containers {
			fmt.Printf("    container[%d]: %s\n", j, cont.Name)
		}
	}
}
