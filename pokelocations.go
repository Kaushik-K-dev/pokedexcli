package main

import ("net/http"; "encoding/json")

var LocationURL = "https://pokeapi.co/api/v2/location-area/"

type LocationData struct {
	Count    int    `json:"count"`
	Next     *string `json:"next"`
	Previous *string    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (cfg *config) LocationReq(url *string) (*LocationData, error) {
	if url == nil {url=&LocationURL}
	dat, ok := cfg.cache.Get(url)
	if ok {
		cfg.prevURL = dat.Previous
		cfg.nextURL = dat.Next
		return dat, nil}

	req, err := http.NewRequest("GET", *url, nil)
	if err!=nil {return nil, err}
	resp, err := http.DefaultClient.Do(req)
	if err!=nil {return nil, err}
	defer resp.Body.Close()

	var LocationResp LocationData
	err = json.NewDecoder(resp.Body).Decode(&LocationResp)
	if err!=nil {return nil, err}

	cfg.cache.Add(url, LocationResp)

	cfg.prevURL = LocationResp.Previous
	cfg.nextURL = LocationResp.Next

	return &LocationResp, nil
}