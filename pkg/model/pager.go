package model

import "github.com/rs/zerolog/log"

type Pager struct {
	ID   int64
	Name string
}

func (pager *Pager) Call() {
	log.Debug().
		Str("pager", pager.Name).
		Msg("pager has been called")

	// TODO: make call to Pager Webapp
}

func GetPagers() ([]*Pager, error) {
	// TODO: retrieve pagers from Pager Webapp

	return nil, nil
}
