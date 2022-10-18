package http

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	v1 "k8s.io/api/core/v1"
)

const (
	timeout = 30 * time.Second
)

type Client struct {
	serverAddr string
	httpClient http.Client
	token      string
}

func NewClient(node, port, caPath, token string) (*Client, error) {
	caCert, err := os.ReadFile(caPath)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	client := &Client{
		serverAddr: "https://" + node + ":" + port,
		httpClient: http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: caCertPool,
				}},
			Timeout: timeout,
		},
		token: token,
	}
	return client, nil
}

func (client *Client) GetPodList() (*v1.PodList, error) {
	url := client.serverAddr + "/pods"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	header := http.Header{}
	header.Add("Authorization", "bearer " + client.token)
	req.Header = header
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP response status: %d, req URL: %s", resp.StatusCode, req.URL.String())
	}
	podList := &v1.PodList{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(podList)
	if err != nil {
		return nil, err
	}
	return podList, nil
}

func (client *Client) GetStr(path string) (string, error) {
	url := client.serverAddr + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	header := http.Header{}
	header.Add("Authorization", "bearer " + client.token)
	req.Header = header
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP response status: %d, req URL: %s", resp.StatusCode, req.URL.String())
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
