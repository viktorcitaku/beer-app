package beerapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type BeerResult struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Client struct {
	beersApiUrl string
	client      *http.Client
}

func NewClient(beersApiUrl string) *Client {
	return &Client{
		beersApiUrl: beersApiUrl,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetBeers(pairingFood, bitterness, fermentation string) ([]BeerResult, error) {
	var err error
	r, err := c.client.Get(c.beersApiUrl)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("Error closing body: %v", err)
		}
	}(r.Body)

	var beers []beer
	err = json.NewDecoder(r.Body).Decode(&beers)
	if err != nil {
		return nil, err
	}

	var results []BeerResult
	for _, b := range beers {
		if filter(b, pairingFood, bitterness, fermentation) {
			results = append(results, BeerResult{
				ID:   b.ID,
				Name: b.Name,
			})
		}
	}

	return results, nil
}

func filter(beer beer, pf, b, f string) bool {
	pfMatch := -1
	if pf != "" {
		sentences := beer.FoodPairing
		for _, sentence := range sentences {
			for _, word := range strings.Split(pf, "_") {
				if strings.Contains(sentence, word) {
					pfMatch = 1
				}
			}
		}
	} else {
		pfMatch = 0
	}

	bMatch := -1
	if b != "" {
		switch b {
		case "high":
			if highIbuScale(beer.Ibu) {
				bMatch = 1
			}
		case "medium":
			if mediumIbuScale(beer.Ibu) {
				bMatch = 1
			}
		case "low":
			if lowIbuScale(beer.Ibu) {
				bMatch = 1
			}
		}
	} else {
		bMatch = 0
	}

	fMatch := -1
	if f != "" {
		temp := beer.Method.Fermentation.Temp.Value
		switch f {
		case "top_fermented":
			if isAleBeer(temp) {
				fMatch = 1
			}
		case "bottom_fermented":
			if isLagerBeer(temp) {
				fMatch = 1
			}
		}
	} else {
		fMatch = 0
	}

	return pfMatch != -1 && bMatch != -1 && fMatch != -1
}

func highIbuScale(ibu float64) bool {
	return ibu > 80 && ibu <= 120
}

func mediumIbuScale(ibu float64) bool {
	return ibu > 40 && ibu <= 80
}

func lowIbuScale(ibu float64) bool {
	return ibu >= 0 && ibu <= 40
}

func isAleBeer(temp float64) bool {
	return temp >= 10 && temp <= 25
}

func isLagerBeer(temp float64) bool {
	return temp >= 7 && temp <= 15
}
