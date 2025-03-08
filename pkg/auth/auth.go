package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const ENDPOINT = "https://monportail.uqam.ca/authentification"

type response struct {
	Token string `json:"token"`
}

func MustGenerateToken(user, pass string) string {
	payload, err := json.Marshal(map[string]string{
		"identifiant": user,
		"motDePasse":  pass,
	})
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", ENDPOINT, bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
	var response response
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatal(err)
	}

	if response.Token == "" {
		log.Fatal("Token is empty")
	}

	return response.Token
}
