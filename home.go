package tibber

import (
	"context"
	"encoding/json"
	"time"

	graphql "github.com/hasura/go-graphql-client"
)

type QueryResponse struct {
	Viewer Viewer
}
type Viewer struct {
	Name   string `json:"name"`
	UserId string `json:"userId"`
	Homes  []Home `json:"homes"`
	Home   Home   `json:"home"`
}
type PreviousMeterData struct {
	Power           float64 `json:"power"`
	PowerProduction float64 `json:"powerProduction"`
}

// Home structure
type Home struct {
	ID                   string              `json:"id"`
	AppNickname          string              `json:"appNickname"`
	MeteringPointData    MeteringPointData   `json:"meteringPointData"`
	Features             Features            `json:"features"`
	Address              Address             `json:"address"`
	Size                 int                 `json:"size"`
	MainFuseSize         int                 `json:"mainFuseSize"`
	NumberOfResidents    int                 `json:"numberOfResidents"`
	PrimaryHeatingSource string              `json:"primaryHeatingSource"`
	HasVentilationSystem bool                `json:"hasVentilationSystem"`
	CurrentSubscription  CurrentSubscription `json:"currentSubscription"`
	PreviousMeterData    PreviousMeterData   `json:"previousMeterData"`
}

type Address struct {
	Address    string `json:"address1"`
	PostalCode string `json:"postalCode"`
	City       string `json:"city"`
	Country    string `json:"country"`
	Latitude   string `json:"latitude"`
	Longitude  string `json:"longitude"`
}

// MeteringPointData - meter number
type MeteringPointData struct {
	ConsumptionEan string `json:"consumptionEan"`
}

// Features - tibber pulse connected
type Features struct {
	RealTimeConsumptionEnabled bool `json:"realTimeConsumptionEnabled"`
}

type CurrentSubscription struct {
	Id        string    `json:"id"`
	ValidFrom time.Time `json:"validFrom"`
	ValidTo   time.Time `json:"validTo"`
	PriceInfo PriceInfo `json:"priceInfo"`
}

type PriceInfo struct {
	Current  Price   `json:"current"`
	Today    []Price `json:"today"`
	Tomorrow []Price `json:"tomorrow"`
}

type Price struct {
	Level    string    `json:"level"`
	Total    float64   `json:"total"`
	Energy   float64   `json:"energy"`
	Tax      float64   `json:"tax"`
	Currency string    `json:"currency"`
	StartsAt time.Time `json:"startsAt"`
}

func (ctx *TibberClient) GetHomes() (QueryResponse, error) {
	var query struct {
		Viewer struct {
			Name   string `graphql:"name"`
			UserId string `graphql:"userId"`
			Homes  struct {
				Id                   string `graphql:"id"`
				TimeZone             string `graphql:"id"`
				AppNickname          string `graphql:"appNickname"`
				Size                 int    `graphql:"size"`
				Type                 string `graphql:"type"`
				NumberOfResidents    int    `graphql:"numberOfResidents"`
				PrimaryHeatingSource string `graphql:"primaryHeatingSource"`
				HasVentilationSystem string `graphql:"hasVentilationSystem "`
				MainFuseSize         string `graphql:"mainFuseSize"`
				Address              struct {
					Address    string `graphql:"address1"`
					City       string `graphql:"city"`
					PostalCode string `graphql:"postalCode"`
					Country    string `graphql:"country"`
					Latitude   string `graphql:"latitude"`
					Longitude  string `graphql:"longitude"`
				} `graphql:"address"`
				CurrentSubscription struct {
					Id        string `graphql:"id"`
					PriceInfo struct {
						Current struct {
							Level    string    `graphql:"level"`
							Total    float64   `graphql:"total"`
							Energy   float64   `graphql:"energy"`
							Tax      float64   `graphql:"tax"`
							Currency string    `graphql:"currency"`
							StartsAt time.Time `graphql:"startsAt"`
						} `graphql:"current"`
						Today []struct {
							Level    string    `graphql:"level"`
							Total    float64   `graphql:"total"`
							Energy   float64   `graphql:"energy"`
							Tax      float64   `graphql:"tax"`
							Currency string    `graphql:"currency"`
							StartsAt time.Time `graphql:"startsAt"`
						} `graphql:"today"`
						Tomorrow []struct {
							Level    string    `graphql:"level"`
							Total    float64   `graphql:"total"`
							Energy   float64   `graphql:"energy"`
							Tax      float64   `graphql:"tax"`
							Currency string    `graphql:"currency"`
							StartsAt time.Time `graphql:"startsAt"`
						} `graphql:"tomorrow"`
					} `graphql:"priceInfo"`
					ValidFrom string `graphql:"validFrom"`
					ValidTo   string `graphql:"validTo"`
				} `graphql:"currentSubscription "`
			} `graphql:"homes"`
		} `graphql:"viewer"`
	}

	data, err := ctx.apiClient.QueryRaw(context.Background(), query, nil)

	var queryResponse QueryResponse

	if err != nil {
		println(err.Error())
		return queryResponse, err
	}

	jsonErr := json.Unmarshal(data, &queryResponse)

	return queryResponse, jsonErr
}

func (ctx *TibberClient) GetPriceInfo(homeId string) (CurrentSubscription, error) {
	var query struct {
		Viewer struct {
			Home struct {
				Id                  string `graphql:"id"`
				CurrentSubscription struct {
					Id        string `graphql:"id"`
					PriceInfo struct {
						Current struct {
							Level    string    `graphql:"level"`
							Total    float64   `graphql:"total"`
							Energy   float64   `graphql:"energy"`
							Tax      float64   `graphql:"tax"`
							Currency string    `graphql:"currency"`
							StartsAt time.Time `graphql:"startsAt"`
						} `graphql:"current"`
						Today []struct {
							Level    string    `graphql:"level"`
							Total    float64   `graphql:"total"`
							Energy   float64   `graphql:"energy"`
							Tax      float64   `graphql:"tax"`
							Currency string    `graphql:"currency"`
							StartsAt time.Time `graphql:"startsAt"`
						} `graphql:"today"`
						Tomorrow []struct {
							Level    string    `graphql:"level"`
							Total    float64   `graphql:"total"`
							Energy   float64   `graphql:"energy"`
							Tax      float64   `graphql:"tax"`
							Currency string    `graphql:"currency"`
							StartsAt time.Time `graphql:"startsAt"`
						} `graphql:"tomorrow"`
					} `graphql:"priceInfo"`
					ValidFrom string `graphql:"validFrom"`
					ValidTo   string `graphql:"validTo"`
				} `graphql:"currentSubscription"`
			} `graphql:"home(id: $id)"`
		} `graphql:"viewer"`
	}

	variables := map[string]interface{}{
		"id": graphql.ID(homeId),
	}

	dataRaw, err := ctx.apiClient.QueryRaw(context.Background(), query, variables)
	if err != nil {
		return CurrentSubscription{}, err
	}

	var result QueryResponse
	jsonErr := json.Unmarshal(dataRaw, &result)

	if jsonErr != nil {
		return CurrentSubscription{}, jsonErr
	}

	return result.Viewer.Home.CurrentSubscription, jsonErr
}
