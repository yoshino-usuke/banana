package helper

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"golang.org/x/net/context"
)

const CONF = "../conf/banana.yml"
const HOST_CONF = "../conf/hosts.yml"

type Config struct {
	CreadentialPath string
	Baseurl         string
	Url             string
	Dbid       string
	Meta            []Db
	Authurl         string
}

func GetClient()*http.Client{
	return client
}

func GetConfig() *Config{
	return config
}

var config *Config
var client *http.Client
var oauthConfig *oauth2.Config

func init(){
	prepare()
}

func prepare(){
	c,err := ioutil.ReadFile(CONF)
	if err != nil{
		log.Fatalln(err)
	}

	err = yaml.Unmarshal(c,&config)
	if err != nil{
		log.Fatalln(err)
	}

	//google認証
	hostCnf ,err := ioutil.ReadFile(HOST_CONF)
	if err != nil{
		log.Fatalln(err)
	}
	err = yaml.Unmarshal(hostCnf,&oauthConfig)
	if err != nil{
		log.Fatalln(err)
	}
	certification()
}

func certification(){
	tok,err := tokenFromFile(config.CreadentialPath)
	if err != nil{
		tok = tokenFromWeb()
		saveToken(tok)
	}
	client = oauthConfig.Client(context.Background(), tok)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func tokenFromWeb() *oauth2.Token {
	authURL := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	log.Println(authCode)
	tok, err := oauthConfig.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func saveToken(tok *oauth2.Token){
	fmt.Printf("Saving credential file to: %s\n", config.CreadentialPath)
	f, err := os.OpenFile(config.CreadentialPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	json.NewEncoder(f).Encode(tok)
}
