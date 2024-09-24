package main

import ("net/http"; "encoding/json"; "fmt")

var LocationURL = "https://pokeapi.co/api/v2/location-area/"

type LocationsList struct {
	Count    int    `json:"count"`
	Next     *string `json:"next"`
	Previous *string    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (cfg *config) LocationListReq(url *string) (*LocationsList, error) {
	if url == nil {url=&LocationURL}
	dat, ok := cfg.cache.Get(*url)
	if ok {
		locationList, ok := dat.(*LocationsList)
        if !ok {return nil, fmt.Errorf("cache entry is not of type *LocationsList")}
		cfg.prevURL = locationList.Previous
		cfg.nextURL = locationList.Next
		return locationList, nil}

	req, err := http.NewRequest("GET", *url, nil)
	if err!=nil {return nil, err}
	resp, err := http.DefaultClient.Do(req)
	if err!=nil {return nil, err}
	defer resp.Body.Close()

	var LocationResp LocationsList
	err = json.NewDecoder(resp.Body).Decode(&LocationResp)
	if err!=nil {return nil, err}

	cfg.cache.Add(*url, &LocationResp)

	cfg.prevURL = LocationResp.Previous
	cfg.nextURL = LocationResp.Next

	return &LocationResp, nil
}

type LocationData struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func (cfg *config) LocationDataReq(AreaName string) (*LocationData, error) {
	url := LocationURL + AreaName
	
	dat, ok := cfg.cache.Get(url)
	if ok {
		locationData, ok := dat.(*LocationData)
        if !ok {return nil, fmt.Errorf("cache entry is not of type *LocationData")}
		return locationData, nil}

	req, err := http.NewRequest("GET", url, nil)
	if err!=nil {return nil, err}
	resp, err := http.DefaultClient.Do(req)
	if err!=nil {return nil, err}
	defer resp.Body.Close()

	var LocationResp LocationData
	err = json.NewDecoder(resp.Body).Decode(&LocationResp)
	if err!=nil {return nil, err}

	cfg.cache.Add(url, &LocationResp)

	return &LocationResp, nil
}