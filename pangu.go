package http

import (
	"github.com/pangum/http/internal/plugin"
	"github.com/pangum/pangu"
)

func init() {
	pangu.New().Get().Dependency().Build().Provide(new(plugin.Creator).New)
}
