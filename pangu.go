package http

import (
	"github.com/pangum/pangu"
)

func init() {
	pangu.New().Dependencies(
		newClient,
		newRequest,
	)
}
