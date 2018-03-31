package main

import (
	"./mylib"
	"log"
	"sync"
	"fmt"
)

const PROJECT = "banana"

func main() {
	project := PROJECT
	updateType := mylib.UPDATE_MONTHLY
	sheetService, d := mylib.GetSheet(project)
	dbList := d.GetDb(updateType)

	fnc := func(db mylib.Db, wait *sync.WaitGroup) {
		res := mylib.GetRes(db)

		results := mylib.Extract(*res)
		err := mylib.Insert(results, *sheetService, db, d.Dbid)
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
