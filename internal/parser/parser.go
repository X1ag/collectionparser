package parser

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Collection struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
	Stats    struct {
		Count int    `json:"count"`
		Floor json.Number `json:"floor"`
	} `json:"stats"`
}

type CollectionData struct {
	Name  string 
	Floor int 
}

const url = "API URL HERE"

func ParseCollections() ([]CollectionData, error) {

	log.Println("Calling ParseCollections")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Default().Printf("Error creating request: %v\n", err)
		return nil, err	
	}
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Default().Printf("Request failed: %v\n", err)
		return nil, err	
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Default().Printf("API error: status %d\n", res.StatusCode)
		return nil, err	
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Default().Printf("Error reading response: %v\n", err)
		return nil, err	
	}

	var collections []Collection
	if err := json.Unmarshal(body, &collections); err != nil {
		log.Default().Printf("JSON parse error: %v\n", err)
		return nil, err
	}

	var results []CollectionData
	targetCollections := map[string]bool{
		"Plush Pepe":     true,
		"Durov's Cap":    true,
		"Precious Peach": true,
		"Astral Shard":   true,
		"Loot Bag":       true,
		"Ion Gem":        true,
		"Scared Cat":     true,
	}

	for _, c := range collections {
		if targetCollections[c.Name] {
			floorInt, err := c.Stats.Floor.Int64()
			if err != nil {
				log.Default().Printf("Error converting floor to int: %v\n", err)
				continue
			}
			results = append(results, CollectionData{
				Name:  c.Name,
				Floor: int(floorInt),
			})
		}
	}
	return results, nil
}
