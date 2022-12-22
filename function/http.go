package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func MakePayload(obj interface{}) *bytes.Buffer {
	reqJson, _ := json.Marshal(obj)
	payload := bytes.NewBuffer(reqJson)

	return payload
}

func HttpCall(method string, url string, headerMap *map[string]string, payload *bytes.Buffer) string {
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	httpReq, err := http.NewRequest(method, url, payload)
	if err != nil {
		panic(fmt.Sprintf("Got error %s", err.Error()))
	}

	for key, value := range *headerMap {
		httpReq.Header.Set(key, value)
	}

	response, err := client.Do(httpReq)
	if err != nil {
		panic(fmt.Sprintf("Got error %s", err.Error()))
	}
	defer response.Body.Close()

	resBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(fmt.Sprintf("Got error %s", err.Error()))
	}

	return string(resBody)
}

func HttpRequest(method string, url string, headerMap *map[string]string, payload *bytes.Buffer) (string, string) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	httpReq, err := http.NewRequest(method, url, payload)
	if err != nil {
		panic(fmt.Sprintf("Got error %s", err.Error()))
	}

	for key, value := range *headerMap {
		httpReq.Header.Set(key, value)
	}

	response, err := client.Do(httpReq)
	if err != nil {
		panic(fmt.Sprintf("Got error %s", err.Error()))
	}
	defer response.Body.Close()

	resBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(fmt.Sprintf("Got error %s", err.Error()))
	}

	return string(resBody), response.Status
}
