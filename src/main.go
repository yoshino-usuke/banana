package main

import (
	"flag"
	"fmt"
	"sync"
	"log"
	"./helper"
)

func main()  {
	var updateType string
	flag.StringVar(&updateType,"term",helper.UPDATE_DAILY,"term[ all | realtime | daily | weekly | monthly ]")
	flag.Parse()
	if ! helper.ValidateArgs(updateType) {
		log.Fatalf("invalid args : %s",updateType)
	}

	sheetService, d := helper.GetSheet()
	dbList := d.GetDb(updateType)

	fnc := func(db helper.Db, wait *sync.WaitGroup) {
		res := helper.GetRes(db)

		results := helper.Extract(*res)
		err := helper.Insert(results, *sheetService, db, d.Dbid)
		defer wait.Done()
		if err != nil {
			log.Fatalf("update failed %s", err.Error())
		}
		log.Printf("update success %s!", db.Sheetname)
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(dbList))
	for _, db := range dbList {

		go fnc(db, &waitGroup)
	}

	waitGroup.Wait()
	fmt.Printf("completed!")

}