package eks

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
)

type node struct {
	Metadata nodeMetadata `json:"metadata"`
}

type nodeMetadata struct {
	Labels map[string]string `json:"labels"`
}

type kubeClient struct {
	client      *http.Client
	host        string
	bearerToken string
}

func newKubeClient() (*kubeClient, error) {
	const (
		tokenFile  = "/var/run/secrets/kubernetes.io/serviceaccount/token"
		rootCAFile = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	)
	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	if len(host) == 0 || len(port) == 0 {
		return nil, errors.New("unable to load in-cluster configuration, KUBERNETES_SERVICE_HOST and KUBERNETES_SERVICE_PORT must be defined")
	}

	token, err := os.ReadFile(tokenFile)
	if err != nil {
		return nil, err
	}

	caCert, err := os.ReadFile(rootCAFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, errors.New("could load certs from CA file")
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}

	return &kubeClient{client: client, host: "https://" + net.JoinHostPort(host, port), bearerToken: string(token)}, nil
}

func (c *kubeClient) getNodeLabels(name string) (map[string]string, error) {
	endpoint, err := url.JoinPath(c.host, "api/v1/nodes", name)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.bearerToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d %s", resp.StatusCode, resp.Status)
	}

	var n node
	if err := json.NewDecoder(resp.Body).Decode(&n); err != nil {
		return nil, err
	}

	return n.Metadata.Labels, nil
}
