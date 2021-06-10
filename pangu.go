package http

import (
	`github.com/storezhang/pangu`
)

func init() {
	if err := pangu.New().Provides(newClient, newRequest); nil != err {
		panic(err)
	}
}
