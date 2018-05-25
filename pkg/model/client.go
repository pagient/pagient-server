package model

type Client struct {
	ID   int64
	Name string
}

func GetClients() ([]*Client, error) {
	// TODO: retrieve clients from config

	return nil, nil
}
