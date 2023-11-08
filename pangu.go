package http

import (
	"github.com/pangum/http/internal/plugin"
	"github.com/pangum/pangu"
)

func init() {
	creator := new(plugin.Constructor)
	pangu.New().Get().Dependency().Put(
		creator.New,
	).Build().Build().Apply()
}
