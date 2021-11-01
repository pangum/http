package http

import (
	`github.com/pangum/pangu`
)

func init() {
	if err := pangu.New().Provides(newClient, newRequest); nil != err {
		panic(err)
	}
}
