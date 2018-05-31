package model

import (
	"strconv"
	"strings"

	"github.com/pagient/pagient-easy-call-go/easycall"
	"github.com/rs/zerolog/log"
)

// Pager struct
type Pager struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Call calls the pager and let it vibrate
func (pager *Pager) Call() error {
	log.Debug().
		Str("pager", pager.Name).
		Msg("pager has been called")

	client := easycall.NewClient(cfg.EasyCall.URL)
	client.SetCredentials(cfg.EasyCall.User, cfg.EasyCall.Password)

	err := client.Send(&easycall.SendOptions{
		Receiver: pager.ID,
		Message:  "",
	})

	return err
}

// GetPagers returns all available pagers
func GetPagers() ([]*Pager, error) {
	pagers := []*Pager{}
	for _, pagerInfo := range cfg.General.Pagers {
		pair := strings.SplitN(pagerInfo, ":", 2)

		id, err := strconv.Atoi(pair[0])
		if err != nil {
			return nil, err
		}

		pagers = append(pagers, &Pager{ID: id, Name: pair[1]})
	}

	return pagers, nil
}

// GetPagerByID returns a single pager by ID
func GetPagerByID(id int) (*Pager, error) {
	// TODO: retrieve pager by id from Pager Webapp

	return nil, nil
}
