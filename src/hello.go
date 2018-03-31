package main

import(
	"net/http"
	"io/ioutil"
)

func main(){
	url := "http://www.dmm.co.jp/digital/videoc/-/detail/ajax-movie/=/cid=gsiro032/"
	res,err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	byteList,err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	println(string(byteList))

}
