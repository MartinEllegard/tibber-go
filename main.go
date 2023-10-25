package main

import (
	"log"
	"os"

	tibberclient "github.com/MartinEllegard/tibber-go/tibber_client"
	"github.com/joho/godotenv"
)

// https://developer.tibber.com/explorer

func main() {
	godotenv.Load()

	token := os.Getenv("TIBBER_TOKEN")
	houseId := os.Getenv("TEST_HOUSE")

	client := tibberclient.CreateTibberClient(token, "test/0.1.0")

	response, err := client.GetHomes()
	if err != nil {
		panic("failed query")
	}

	println(response.Viewer.Name)

	client.StartSubscription(houseId, test_handler)
	client.Close()
}

func test_handler(data []byte, err error) error {

	if err != nil {
		log.Println("ERROR: ", err)
		return nil
	}

	if data == nil {
		return nil
	}
	log.Println(string(data))
	return nil
}
