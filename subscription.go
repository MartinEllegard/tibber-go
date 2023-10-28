package tibber

import (
	"log"
	"net/http"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/hasura/go-graphql-client/pkg/jsonutil"
)

type subscriptionResponse struct {
	LiveMeasurement LiveMeasurement
}

type LiveMeasurement struct {
	HomeId                         string    `json:"homeId"`
	Timestamp                      time.Time `json:"timestamp"`
	Power                          float64   `json:"power"`
	MinPower                       float64   `json:"minPower"`
	AveragePower                   float64   `json:"averagePower"`
	MaxPower                       float64   `json:"maxPower"`
	LastMeterConsumption           float64   `json:"lastMeterConsumption"`
	LastMeterProduction            float64   `json:"lastMeterProduction"`
	AccumulatedConsumption         float64   `json:"accumulatedConsumption"`
	AccumulatedProduction          float64   `json:"accumulatedProduction"`
	AccumulatedCost                float64   `json:"accumulatedCost"`
	AccumulatedConsumptionLastHour float64   `json:"accumulatedConsumptionLastHour"`
	AccumulatedProductionLastHour  float64   `json:"accumulatedProductionLastHour"`
	Currency                       string    `json:"currency"`
}

func (ctx *TibberClient) StartSubscription(homeId string, outputChannel chan<- LiveMeasurement) error {
	if homeId == "" {
		return graphql.Error{Message: "missing argument homeId"}
	}

	var sub struct {
		LiveMeasurement struct {
			Timestamp                      time.Time `graphql:"timestamp"`
			Power                          float64   `graphql:"power"`
			MinPower                       float64   `graphql:"minPower"`
			AveragePower                   float64   `graphql:"averagePower"`
			MaxPower                       float64   `graphql:"maxPower"`
			LastMeterConsumption           float64   `graphql:"lastMeterConsumption"`
			LastMeterProduction            float64   `graphql:"lastMeterProduction"`
			AccumulatedConsumption         float64   `graphql:"accumulatedConsumption"`
			AccumulatedProduction          float64   `graphql:"accumulatedProduction"`
			AccumulatedCost                float64   `graphql:"accumulatedCost"`
			AccumulatedConsumptionLastHour float64   `graphql:"accumulatedConsumptionLastHour"`
			AccumulatedProductionLastHour  float64   `graphql:"accumulatedProductionLastHour"`
			Currency                       string    `graphql:"currency"`
		} `graphql:"liveMeasurement(homeId: $homeId)"`
	}

	variables := map[string]interface{}{
		"homeId": graphql.ID(homeId),
	}
	_, err := ctx.wsClient.Subscribe(sub, variables, func(message []byte, err error) error {
		data := subscriptionResponse{}

		if err != nil {
			return err
		}

		jsonerror := jsonutil.UnmarshalGraphQL(message, &data)

		if (jsonerror == nil && data.LiveMeasurement != LiveMeasurement{}) {
			data.LiveMeasurement.HomeId = homeId
			outputChannel <- data.LiveMeasurement
			return nil
		}

		return nil
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	return ctx.wsClient.Run()
}

type headerRoundTripper struct {
	setHeaders func(req *http.Request)
	rt         http.RoundTripper
}

func (h headerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	h.setHeaders(req)
	return h.rt.RoundTrip(req)
}
