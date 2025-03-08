package grade

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const ENDPOINT = "https://monportail.uqam.ca/apis/resultatActivite/identifiant"

type GradeResponse struct {
	Data struct {
		Resultats []struct {
			Programmes []struct {
				Activites []struct {
					Total string `json:"total"`
				} `json:"activites"`
			} `json:"programmes"`
		} `json:"resultats"`
	} `json:"data"`
}

func FetchGrade(token, semester, course string) (string, error) {
	url := fmt.Sprintf("%s/%s/%s", ENDPOINT, semester, course)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response GradeResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	return response.Data.Resultats[0].Programmes[0].Activites[0].Total, nil
}
