package mylib

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

const CONF_FILE = "../conf/%s.yml"
const URL = "http://www.dmm.co.jp/digital/videoc/-/ranking/=/term=%s/"
const (
	UPDATE_ALL = "all"
	UPDATE_RANKING = "realtime"
	UPDATE_WEEK = "weekly"
	UPDATE_DAILY = "daily"
	UPDATE_MONTHLY = "monthly"

)

type Data struct {
	CreadentialPath string
	Baseurl         string
	Url             string
	Dburl           string
	Dbid       string
	Meta            []Db
	Authurl         string
}

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
	return res
}

func ValueRange(result Result) []interface{} {

	return []interface{}{result.rank, result.name, result.img, result.url,result.information}
}

func (d Data) GetClient() (*http.Client, error) {
	// 認証情報
	credential, err := ioutil.ReadFile(d.CreadentialPath)
	if err != nil {
		return nil, err
	}

	gConf, err := google.JWTConfigFromJSON(credential, d.Authurl)
	if err != nil {
		return nil, err
	}

	client := gConf.Client(oauth2.NoContext)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return client, nil
}

func (d Data) GetDb(updateType string) []Db{
	var typeList []string
	switch updateType {
	case UPDATE_ALL:
		typeList = append(typeList,UPDATE_RANKING,UPDATE_WEEK,UPDATE_DAILY,UPDATE_MONTHLY)

	case UPDATE_RANKING,UPDATE_WEEK,UPDATE_DAILY,UPDATE_MONTHLY:
		typeList = append(typeList,updateType)
	}
	var results []Db
	for _,d := range d.Meta{
		for _,t := range typeList{
			if d.Sheetname == t{
				results = append(results,d)
			}
		}
	}
	return results
}

func getData(project string) (*Data, error) {
	path := fmt.Sprintf(CONF_FILE, project)
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var d Data
	err = yaml.Unmarshal(file, &d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func GetSheet(project string) (*sheets.Service, *Data) {
	d, err := getData(project)
	if err != nil {
		log.Fatalf(err.Error())
	}

	client, err := d.GetClient()
	if err != nil {
		log.Fatalf(err.Error())
	}
	sheetService, err := sheets.New(client)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return sheetService, d
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
