package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"

	"acon3d.com/function"
	"github.com/pelletier/go-toml/v2"
)

const (
	UT_NAME   = "acon3d assessment unit test"
	UT_SCHEMA = "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
	UT_URL    = "https://api.getpostman.com/collections/10804095-e761a4b8-9920-4252-a7f1-52b7a46734ab"

	UAT_NAME   = "acon3d assessment acceptance test"
	UAT_SCHEMA = "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
	UAT_URL    = "https://api.getpostman.com/collections/10804095-80f41f9e-2a8c-4083-8dcc-c5f85e5eb83b"
)

func main() {
	// Read test cases
	testCases := readTestCases()

	// create unittest
	createUnittest(testCases)

	// create uat
	createUAT(testCases)

}

func createUnittest(testCases *TestCases) {
	var cfg PostmanConfig

	cfg.Collection.Info.Name = UT_NAME
	cfg.Collection.Info.Schema = UT_SCHEMA
	cfg.Collection.Item = make([]Folder, 0)

	// add test cases
	for _, f := range testCases.Folder {
		var folder Folder

		folder.Name = f.Name

		for _, tc := range f.TestCases {
			var request Request

			request.Name = tc.Name
			request.ProtocolProfileBehavior.DisableBodyPruning = true
			request.Request.Body.Mode = "raw"
			request.Request.Body.Options.Raw.Language = "json"
			request.Request.Body.Raw = tc.Payload
			request.Request.URL = fmt.Sprintf("%s%s", testCases.Endpoint, tc.URL)
			request.Request.Method = tc.Method

			// making header
			if tc.UserId != "" {
				request.Request.Header = make([]Header, 0)

				request.Request.Header = append(request.Request.Header, Header{
					Key:   "user_id",
					Type:  "text",
					Value: tc.UserId,
				})
			} else {
				request.Request.Header = make([]Header, 0)
			}

			// if script is Empty
			if tc.Script == "" {
				tc.Script = fmt.Sprintf(`
				pm.test("'%s'Status code is 200", function () {
					pm.response.to.have.status(200);
				});				
				`, request.Name)
			}

			request.Event = append(request.Event, Event{
				Listen: "test",
				Script: Script{
					Type: "text/javascript",
					Exec: []string{tc.Script},
				},
			})

			folder.Requests = append(folder.Requests, request)
		}

		cfg.Collection.Item = append(cfg.Collection.Item, folder)
	}

	// add variable
	for _, v := range testCases.Variable {
		cfg.Collection.Variable = append(cfg.Collection.Variable, Variable{Key: v.Key})
	}

	// creat new config file
	jsonPostmanConfig, err := json.Marshal(cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonPostmanConfig))

	//////////////////////////////////////////////////////////////////////////////////
	// Create Unit Test
	url := UT_URL // send to postman

	payload := function.MakePayload(cfg) // create payload

	// create header
	var headerMap = map[string]string{
		"Content-Type": "application/json",
		"X-Api-Key":    "PMAK-63a2c8dc270b5562cff5e9e1-aa818cd2eb8cbbbff83e66354a99270a6a",
	}

	// create Unit test
	resBodyUnitTest := function.HttpCall("PUT", url, &headerMap, payload)

	// http call
	fmt.Printf("%+v\n", resBodyUnitTest)
}

// newman run https://api.getpostman.com/collections/10804095-80f41f9e-2a8c-4083-8dcc-c5f85e5eb83b?apikey=PMAK-63a2c8dc270b5562cff5e9e1-aa818cd2eb8cbbbff83e66354a99270a6a --env-var "endpoint=localhost:8080"
func createUAT(testCases *TestCases) {
	// extract test_cases for UAT
	req := make(map[int]StructTestCase)

	for _, f := range testCases.Folder {
		for _, tc := range f.TestCases {
			if tc.Uat != 0 {
				req[tc.Uat] = tc
			}
		}
	}

	// sort
	keys := make([]int, 0, len(req))
	for k := range req {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	// set general
	var cfg PostmanUATConfig

	cfg.Collection.Info.Name = UAT_NAME
	cfg.Collection.Info.Schema = UAT_SCHEMA
	cfg.Collection.Item = make([]Request, 0)

	// iterate
	for _, k := range keys {
		var tmp Request

		tmp.Name = req[k].Name
		tmp.ProtocolProfileBehavior.DisableBodyPruning = true
		tmp.Request.Body.Mode = "raw"
		tmp.Request.Body.Options.Raw.Language = "json"
		tmp.Request.Body.Raw = req[k].Payload
		tmp.Request.URL = fmt.Sprintf("%s%s", testCases.Endpoint, req[k].URL)
		tmp.Request.Method = req[k].Method

		// making header
		if req[k].UserId != "" {
			tmp.Request.Header = make([]Header, 0)

			tmp.Request.Header = append(tmp.Request.Header, Header{
				Key:   "user_id",
				Type:  "text",
				Value: req[k].UserId,
			})
		} else {
			tmp.Request.Header = make([]Header, 0)
		}

		// if script is Empty
		if req[k].Script == "" {
			script := fmt.Sprintf(`
			pm.test("'%s'Status code is 200", function () {
				pm.response.to.have.status(200);
			});
			`, tmp.Name)
			tmp.Event = append(tmp.Event, Event{
				Listen: "test",
				Script: Script{
					Type: "text/javascript",
					Exec: []string{script},
				},
			})
		} else {
			tmp.Event = append(tmp.Event, Event{
				Listen: "test",
				Script: Script{
					Type: "text/javascript",
					Exec: []string{req[k].Script},
				},
			})
		}

		cfg.Collection.Item = append(cfg.Collection.Item, tmp)
	}

	// add variable
	for i := 0; i < len(testCases.Variable); i++ {
		cfg.Collection.Variable = append(cfg.Collection.Variable, Variable{Key: testCases.Variable[i].Key})
	}

	jsonPostmanConfig, err := json.Marshal(cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonPostmanConfig))

	// url :=  // send to postman

	payload := function.MakePayload(cfg) // create payload

	// create header
	var headerMap = map[string]string{
		"Content-Type": "application/json",
		"X-Api-Key":    "PMAK-63a2c8dc270b5562cff5e9e1-aa818cd2eb8cbbbff83e66354a99270a6a",
	}

	// create Unit test
	resBodyUnitTest := function.HttpCall("PUT", UAT_URL, &headerMap, payload)

	// http call
	fmt.Printf("%+v\n", resBodyUnitTest)
}

////////////////////////////////////////////////////////////////////////////
// postman config
////////////////////////////////////////////////////////////////////////////
type Header struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Script struct {
	Exec []string `json:"exec"`
	Type string   `json:"type"`
}

type Event struct {
	Listen string `json:"listen"`
	Script Script `json:"script"`
}

type Request struct {
	Name                    string `json:"name"`
	ProtocolProfileBehavior struct {
		DisableBodyPruning bool `json:"disableBodyPruning"`
	} `json:"protocolProfileBehavior"`
	Event   []Event `json:"event,omitempty"`
	Request struct {
		Body struct {
			Mode    string `json:"mode"`
			Options struct {
				Raw struct {
					Language string `json:"language"`
				} `json:"raw"`
			} `json:"options"`
			Raw string `json:"raw"`
		} `json:"body"`
		Header []Header `json:"header"`
		Method string   `json:"method"`
		URL    string   `json:"url"`
	} `json:"request"`
	// Response []interface{} `json:"response"`
}

type Folder struct {
	Name     string    `json:"name"`
	Requests []Request `json:"item"`
}

type Variable struct {
	Key string `json:"key"`
}

type PostmanConfig struct {
	Collection struct {
		Info struct {
			Name   string `json:"name"`
			Schema string `json:"schema"`
			// UpdatedAt time.Time `json:"updatedAt"`
		} `json:"info"`
		Item     []Folder   `json:"item"`
		Variable []Variable `json:"variable"`
	} `json:"collection"`
}

type PostmanUATConfig struct {
	Collection struct {
		Info struct {
			Name   string `json:"name"`
			Schema string `json:"schema"`
			// UpdatedAt time.Time `json:"updatedAt"`
		} `json:"info"`
		Item     []Request  `json:"item"`
		Variable []Variable `json:"variable"`
	} `json:"collection"`
}

////////////////////////////////////////////////////////////////////////////
// testcases
////////////////////////////////////////////////////////////////////////////
type StructTestCase struct {
	Name    string `toml:"name"`
	Method  string `toml:"method"`
	URL     string `toml:"url"`
	UserId  string `toml:"user_id"`
	Payload string `toml:"payload"`
	Script  string `toml:"script,omitempty"`
	Uat     int    `toml:"uat"`
}

type TestCases struct {
	Endpoint string `toml:"endpoint"`
	Folder   []struct {
		Name      string           `toml:"name"`
		TestCases []StructTestCase `toml:"test_cases"`
	} `toml:"folder"`
	Variable []struct {
		Key string `toml:"key"`
	} `toml:"variable"`
}

func readTestCases() *TestCases {
	f, err := ioutil.ReadFile("./test_case.toml")
	if err != nil {
		panic("can not read test case file")
	}

	var ret TestCases
	err = toml.Unmarshal(f, &ret)
	if err != nil {
		panic(err)
	}

	return &ret
}

// newman run https://api.getpostman.com/collections/7ebf3ed7-5ed3-4aa4-8741-12e56f8ed090/folder/17304433-97d8ba9e-857b-4ce5-a7e8-443c754f9c82?apikey=PMAK-62ccbe7b33de391acae3b31e-ae0c974ed898bcc6b4846a2a0ccfd16421
