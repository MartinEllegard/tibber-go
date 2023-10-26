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
	Timestamp                      time.Time `json:"timestamp"`
	Power                          int       `json:"power"`
	MinPower                       int       `json:"minPower"`
	AveragePower                   float32   `json:"averagePower"`
	MaxPower                       float32   `json:"maxPower"`
	LastMeterConsumption           float32   `json:"lastMeterConsumption"`
	LastMeterProduction            float32   `json:"lastMeterProduction"`
	AccumulatedConsumption         float32   `json:"accumulatedConsumption"`
	AccumulatedProduction          float32   `json:"accumulatedProduction"`
	AccumulatedCost                float32   `json:"accumulatedCost"`
	AccumulatedConsumptionLastHour float32   `json:"accumulatedConsumptionLastHour"`
	AccumulatedProductionLastHour  float32   `json:"accumulatedProductionLastHour"`
	Currency                       string    `json:"currency"`
}

type SubscriptionHandler func(LiveMeasurement, error) error

func (ctx *TibberClient) StartSubscription(houseId string, handler SubscriptionHandler) error {
	// get the demo token from the graphiql playground

	if houseId == "" {
		return graphql.Error{Message: "missing argument houseId"}
	}
	log.Println("Trying")
	var sub struct {
		LiveMeasurement struct {
			Timestamp                      time.Time `graphql:"timestamp"`
			Power                          int       `graphql:"power"`
			MinPower                       int       `graphql:"minPower"`
			AveragePower                   float32   `graphql:"averagePower"`
			MaxPower                       float32   `graphql:"maxPower"`
			LastMeterConsumption           float32   `graphql:"lastMeterConsumption"`
			LastMeterProduction            float32   `graphql:"lastMeterProduction"`
			AccumulatedConsumption         float32   `graphql:"accumulatedConsumption"`
			AccumulatedProduction          float32   `graphql:"accumulatedProduction"`
			AccumulatedCost                float32   `graphql:"accumulatedCost"`
			AccumulatedConsumptionLastHour float32   `graphql:"accumulatedConsumptionLastHour"`
			AccumulatedProductionLastHour  float32   `graphql:"accumulatedProductionLastHour"`
			Currency                       string    `graphql:"currency"`
		} `graphql:"liveMeasurement(homeId: $homeId)"`
	}

	variables := map[string]interface{}{
		"homeId": graphql.ID(houseId),
	}
	_, err := ctx.wsClient.Subscribe(sub, variables, func(message []byte, err error) error {
		data := subscriptionResponse{}

		jsonutil.UnmarshalGraphQL(message, &data)
		return handler(data.LiveMeasurement, err)
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
