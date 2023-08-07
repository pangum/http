package http

import (
	"github.com/pangum/http/internal/plugin"
	"github.com/pangum/pangu"
)

func init() {
	pangu.New().Dependencies(
		plugin.NewClient,
		plugin.NewRequest,
	)
}
