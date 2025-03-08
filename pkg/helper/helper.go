package helper

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func MustLoadConfig(c interface{}) {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(data, c); err != nil {
		log.Fatal(err)
	}
}
