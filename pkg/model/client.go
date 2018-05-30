package model

import (
	"strconv"
	"strings"
)

// Client struct
type Client struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetClients returns all configured clients
func GetClients() ([]*Client, error) {
	clients := []*Client{}
	for _, clientInfo := range cfg.General.Clients {
		pair := strings.SplitN(clientInfo, ":", 2)

		id, err := strconv.Atoi(pair[1])
		if err != nil {
			return nil, err
		}

		clients = append(clients, &Client{ID: id, Name: pair[0]})
	}

	return clients, nil
}

// GetClient returns a client by name
func GetClient(name string) (*Client, error) {
	clientID, err := cfg.General.GetClientID(name)
	if err != nil {
		return nil, err
	}

	client := &Client{
		ID:   clientID,
		Name: name,
	}

	return client, nil
}
