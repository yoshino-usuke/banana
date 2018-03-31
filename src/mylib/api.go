package mylib

import (
	"net/http"
	"log"
)

func GetRes(db Db) *http.Response {
	res, err := http.Get(db.GetUrl())

	// defer res.Body.Close()

	if err != nil {
		log.Fatalf(err.Error())
	}
	return res
}

