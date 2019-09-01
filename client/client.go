package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func Call(method string, url string) {
	req, _ := http.NewRequest(method, url, nil)

	client := new(http.Client)
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byteArray))

}
