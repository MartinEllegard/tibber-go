package tibber

import (
	"fmt"
	"log"
	"net/http"

	graphql "github.com/hasura/go-graphql-client"
)

// https://developer.tibber.com/explorer

type TibberClient struct {
	token     string
	userAgent string
	apiClient *graphql.Client
	wsClient  *graphql.SubscriptionClient
}

const (
	apiEndpoint          = "https://api.tibber.com/v1-beta/gql"
	subscriptionEndpoint = "wss://websocket-api.tibber.com/v1-beta/gql/subscriptions"
	version              = "0.1.0"
)

func CreateTibberClient(token string, agent string) *TibberClient {
	userAgent := agent + " tibber-go/" + version

	apiClient := graphql.NewClient(apiEndpoint, &http.Client{
		Transport: headerRoundTripper{
			setHeaders: func(req *http.Request) {
				req.Header.Set("User-Agent", userAgent)
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			},
			rt: http.DefaultTransport,
		},
	})

	wsclient := graphql.NewSubscriptionClient(subscriptionEndpoint).
		WithProtocol(graphql.GraphQLWS).
		WithWebSocketOptions(graphql.WebsocketOptions{
			HTTPClient: &http.Client{
				Transport: headerRoundTripper{
					setHeaders: func(req *http.Request) {
						req.Header.Set("User-Agent", userAgent)
					},
					rt: http.DefaultTransport,
				},
			},
		}).
		WithConnectionParams(map[string]interface{}{
			"token": token,
		}).WithLog(log.Println).
		OnError(func(sc *graphql.SubscriptionClient, err error) error {
			panic(err)
		})

	return &TibberClient{
		token,
		userAgent,
		apiClient,
		wsclient,
	}
}

func (ctx *TibberClient) Close() {
	defer ctx.wsClient.Close()
}
