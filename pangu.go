package http

import (
	"github.com/pangum/http/internal/plugin"
	"github.com/pangum/pangu"
)

func init() {
	creator := new(plugin.Creator)
	pangu.New().Dependencies(
		creator.New,
	)
}
