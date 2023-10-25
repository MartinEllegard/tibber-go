package tibberclient

import (
	"net/http"
	"time"

	"github.com/hasura/go-graphql-client"
)

type SubscriptionHandler func([]byte, error) error

func (ctx *TibberClient) StartSubscription(houseId string, handler SubscriptionHandler) error {
	// get the demo token from the graphiql playground
	if houseId == "" {
		return graphql.Error{Message: "missing argument houseId"}
	}

	var sub struct {
		LiveMeasurement struct {
			Timestamp              time.Time `graphql:"timestamp"`
			Power                  int       `graphql:"power"`
			AccumulatedConsumption float64   `graphql:"accumulatedConsumption"`
			AccumulatedCost        float64   `graphql:"accumulatedCost"`
			Currency               string    `graphql:"currency"`
			MinPower               int       `graphql:"minPower"`
			AveragePower           float64   `graphql:"averagePower"`
			MaxPower               float64   `graphql:"maxPower"`
		} `graphql:"liveMeasurement(homeId: $homeId)"`
	}

	variables := map[string]interface{}{
		"homeId": graphql.ID(houseId),
	}
	_, err := ctx.wsClient.Subscribe(sub, variables, handler)

	if err != nil {
		panic(err)
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
