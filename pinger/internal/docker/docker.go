package docker

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
)

type container struct {
	NetworkSettings struct {
		Networks map[string]struct {
			IPAddress string
		}
	}
}

func GetRunningContainerIps() ([]string, error) {
	tr := &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return net.Dial("unix", "/var/run/docker.sock")
		},
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", "http://localhost/containers/json", nil)
	if err != nil {
		return nil, err
	}

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var ips []string
	var containers []container

	err = json.NewDecoder(r.Body).Decode(&containers)
	if err != nil {
		return nil, err
	}

	for _, c := range containers {
		for _, net := range c.NetworkSettings.Networks {
			if len(net.IPAddress) != 0 {
				ips = append(ips, net.IPAddress)
			}
		}
	}

	return ips, nil
}
