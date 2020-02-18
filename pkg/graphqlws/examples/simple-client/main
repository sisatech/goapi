package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"code.vorteil.io/vorteil/libs/graphqlws"
)

func main() {
	client, err := graphqlws.NewClient(context.TODO(), &graphqlws.ClientConfig{
		Address: "localhost:8000",
		Path:    "/websocket",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// subscribe to data
	subscription, err := client.Subscription(&graphqlws.SubscriptionConfig{
		Query: `subscription {
	data
}`,
		DataCallback: func(payload *graphqlws.GQLDataPayload) {
			fmt.Printf("DATA: %v\n", payload.Data)
			for _, gqlError := range payload.Errors {
				fmt.Printf("ERROR: %s\n", gqlError.Error())
			}
		},
		ErrorCallback: func(err error) {
			panic(err)
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// mutate data
	data, err := client.Mutation(context.TODO(), &graphqlws.MutationConfig{
		Query: `mutation {
	replaceData(data: "new data")
}`})
	if err != nil {
		log.Fatal(err)
	}

	// Check for errors, otherwise discard data (the subscription should be
	// updated with the new information anyway).
	if len(data.Errors) > 1 {
		for _, gqlError := range data.Errors {
			fmt.Printf("ERROR: %s\n", gqlError.Error())
		}
		log.Fatal(errors.New("server returned errors to the mutation"))
	}

	// stop the subscription
	// NOTE: the client's Shutdown method will do this anyway, it's only
	// 	done here for illustrative purposes.
	subscription.Stop()
	err = subscription.WaitUntilFinished(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// shutdown the client
	err = client.Shutdown(context.TODO())
	if err != nil {
		fmt.Printf("Failed to shutdown the GraphQL web socket client gracefully: %v\n", err)
		return
	}

}
