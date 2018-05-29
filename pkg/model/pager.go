package model

import "github.com/rs/zerolog/log"

// Pager struct
type Pager struct {
	ID   int
	Name string
}

// Call calls the pager and let it vibrate
func (pager *Pager) Call() error {
	log.Debug().
		Str("pager", pager.Name).
		Msg("pager has been called")

	// TODO: make call to Pager Webapp

	return nil
}

// GetPagers returns all available pagers
func GetPagers() ([]*Pager, error) {
	// TODO: retrieve pagers from Pager Webapp

	return nil, nil
}

// GetPagerByID returns a single pager by ID
func GetPagerByID(id int) (*Pager, error) {
	// TODO: retrieve pager by id from Pager Webapp

	return nil, nil
}
