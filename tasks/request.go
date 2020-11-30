package tasks

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func apiRequest(method string, u string, body interface{}) []byte {

	b, err := json.Marshal(body)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	req, err := http.NewRequest(method, u, bytes.NewBuffer(b))
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return respData
}
