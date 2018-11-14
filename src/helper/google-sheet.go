package helper

import (
	"log"

	"errors"
	"fmt"
	"google.golang.org/api/sheets/v4"
	"strconv"
)

const URL = "http://www.dmm.co.jp/digital/videoc/-/ranking/=/term=%s/"
const (
	UPDATE_ALL = "all"
	UPDATE_RANKING = "realtime"
	UPDATE_WEEK = "weekly"
	UPDATE_DAILY = "daily"
	UPDATE_MONTHLY = "monthly"

)

type Db struct {
	Sheetname  string
	Columns    []string
	ColumnRows []string
}

func (db Db) SheetId(sheet sheets.Spreadsheet) (int64, error) {
	for _, val := range sheet.Sheets {
		if val.Properties.Title == db.Sheetname {
			return val.Properties.SheetId, nil
		}

	}
	return 0, errors.New("not found sheet")
}

func (db Db) GetUrl() string{
	var replace string
	if db.Sheetname != UPDATE_ALL{
		replace = db.Sheetname
	}
	return  fmt.Sprintf(URL,replace)
}

func (d Db) ColumnRange(idx int) string {
	var res string
	kF, vL := d.ColumnRows[0], d.ColumnRows[len(d.ColumnRows)-1]
	no := strconv.Itoa(idx + 1)
	res = kF + "2:" + vL + no
	return d.Sheetname + "!" + res
}

func ValueRange(result Result) []interface{} {

	return []interface{}{result.rank, result.name, result.img, result.url,result.information}
}

func ValidateArgs(val string) bool {
	var list = []string{UPDATE_ALL, UPDATE_RANKING, UPDATE_MONTHLY, UPDATE_WEEK, UPDATE_DAILY}
	for _,uType := range list{
		if uType == val {
			return  true
		}
	}
	return false
}

func (c *Config) GetDb(updateType string) []Db{
	var typeList []string
	switch updateType {
	case UPDATE_ALL:
		typeList = append(typeList,UPDATE_RANKING,UPDATE_WEEK,UPDATE_DAILY,UPDATE_MONTHLY)

	case UPDATE_RANKING,UPDATE_WEEK,UPDATE_DAILY,UPDATE_MONTHLY:
		typeList = append(typeList,updateType)
	}
	var results []Db
	for _,d := range c.Meta{
		for _,t := range typeList{
			if d.Sheetname == t{
				results = append(results,d)
			}
		}
	}
	return results
}

func GetSheet() (*sheets.Service, *Config) {
	client := GetClient()
	c := GetConfig()
	sheetService, err := sheets.New(client)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return sheetService, c
}

func Insert(results []Result, s sheets.Service,d Db,sheetId string) error {

	var res [][]interface{}
	var rangeIdx int
	for _, result := range results {
		res = append(res, ValueRange(result))
		rangeIdx++
	}

	vR := sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         res,
	}
	_, err := s.Spreadsheets.Values.Update(sheetId, d.ColumnRange(rangeIdx), &vR).ValueInputOption("USER_ENTERED").Do()

	if err != nil {
		return err
	}
	return nil
}
