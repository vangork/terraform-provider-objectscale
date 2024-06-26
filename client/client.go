package client

import (
	objectscale "github.com/vangork/objectscale-client/golang/pkg"
)

type Client struct {
	ObsClient *objectscale.Client
}

func NewClient(
	endpoint string,
	insecure bool,
	username string,
	password string,
) (*Client, error) {
	obsClient, err := objectscale.NewClient(endpoint, username, password, insecure)
	if err != nil {
		return nil, err
	}

	client := Client{
		ObsClient: obsClient,
	}

	return &client, nil
}
