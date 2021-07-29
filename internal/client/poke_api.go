package client

import (
	"pokedex/internal/domain"
	"encoding/json"
	"net/http"
	"time"
)

type PokeApiClient struct {
	host string
	client *http.Client
}

type searchPokemonResults struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Stats []struct {
		BaseStat string `json:"base_stat"`
		Effort int `json:"effort"`
		Stat struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
}

type encounter struct {
	LocationArea struct {
		Name string `json:"name"`
	} `json:"location_area"`
	VersionDetails []struct {
		EncounterDetails []struct {
			Method struct{
				Name string `json:"name"`
			} `json:"method"`
		} `json:"encounter_details"`
		MaxChance int `json:"max_chance"`
		Version struct {
			Name string `json:"name"`
		} `json:"version"`
	} `json:"version_details"`
}

func NewPokeApiClient() *PokeApiClient {
	return &PokeApiClient{
		"https://pokeapi.co/api/v2",
		&http.Client{
			Timeout: time.Duration(time.Second * 5),
		},
	}
}

func (p *PokeApiClient) GetBasicInfo(keyword string) (*domain.PokeData, error) {
	url := p.host + "/pokemon/" + keyword
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	var res searchPokemonResults
	json.NewDecoder(resp.Body).Decode(&res)

	typesArr := make([]string, 0)
	for _, t := range res.Types {
		typesArr = append(typesArr, t.Type.Name)
	}
	statsArr := make([]domain.StatItem, len(res.Stats))
	for i, statItem := range res.Stats {
		statsArr[i] = domain.StatItem{
			BaseStat: statItem.BaseStat,
			Effort: statItem.Effort,
			StatName: statItem.Stat.Name,
		}
	}
	return &domain.PokeData{
		Id: res.Id,
		Name: res.Name,
		Types: typesArr,
		Stats: statsArr,
	}, nil
}

func (p *PokeApiClient) GetEncounters(keyword string) ([]domain.LocationAreaEncounter, error) {
	url := p.host + "/pokemon/" + keyword + "/encounters"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	var res []encounter
	json.NewDecoder(resp.Body).Decode(&res)

	encounters := make([]domain.LocationAreaEncounter, len(res))
	for i, item := range res {
		details := make([]domain.VersionDetail, len(item.VersionDetails))
		for n, detail := range item.VersionDetails {
			details[n] = domain.VersionDetail{
				MaxChance: detail.MaxChance,
			}
		}
		encounters[i] = domain.LocationAreaEncounter{
			Name: item.LocationArea.Name,
			VersionDetails: details,
		}
	}
	return encounters, nil
}
