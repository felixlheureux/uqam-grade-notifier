package helper

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func MustLoadConfig(path string, c interface{}) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(data, c); err != nil {
		log.Fatal(err)
	}
}
