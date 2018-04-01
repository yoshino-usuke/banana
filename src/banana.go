package main

import (
	"./helper"
	"log"
	"sync"
	"fmt"
)

const PROJECT = "banana"

func main() {
	project := PROJECT
	updateType := helper.UPDATE_MONTHLY
	sheetService, d := helper.GetSheet(project)
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
