package client

import (
	objectscale "github.com/vangork/objectscale-client/golang/pkg"
)

type Client struct {
	ManagementClient *objectscale.ManagementClient
}

func NewClient(
	endpoint string,
	insecure bool,
	username string,
	password string,
) (*Client, error) {
	managementClient, err := objectscale.NewManagementClient(endpoint, username, password, insecure)
	if err != nil {
		return nil, err
	}

	client := Client{
		ManagementClient: managementClient,
	}

	return &client, nil
}
