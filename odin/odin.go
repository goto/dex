package odin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type odinStream struct {
	Urn           string `json:"urn"`
	Name          string `json:"name"`
	Entity        string `json:"entity"`
	Organization  string `json:"organization"`
	Landscape     string `json:"landscape"`
	Environment   string `json:"environment"`
	Type          string `json:"type"`
	AdvertiseMode struct {
		Host    string `json:"host"`
		Address string `json:"address"`
	} `json:"advertise_mode"`
	Brokers   []broker `json:"brokers"`
	Created   string   `json:"created"`
	Updated   string   `json:"updated"`
	ProjectID string   `json:"projectID"`
	URL       string   `json:"url"`
	ID        string   `json:"id"`
}

type broker struct {
	Name    string `json:"name"`
	Host    string `json:"host"`
	Address string `json:"address"`
}

func GetOdinStream(odinAddr, urn string) (string, error) {
	url := fmt.Sprintf("%s/streams/%s", odinAddr, urn)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var odinResp odinStream
	err = json.Unmarshal(body, &odinResp)
	if err != nil {
		return "", err
	}

	return odinResp.URL, nil
}
