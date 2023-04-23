package commenting

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SayIfOrg/say_keeper/graph/gmodel"
	"github.com/redis/go-redis/v9"
	"log"
)

const rChannel string = "commentingChan"

//SubscribeComment subscribes to a redis channel named "rChannel" for new comments
func SubscribeComment(ctx context.Context, rdb *redis.Client, ch chan *gmodel.Comment) {
	pubsub := rdb.Subscribe(ctx, rChannel)
	defer func() {
		err := pubsub.Close()
		if err != nil {
			panic(err)
		}
		log.Println("sub closed")
	}()

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Println(msg.Channel, msg.Payload)

		// Deserialize message from JSON string into object
		receivedObj := &gmodel.Comment{}
		err = json.Unmarshal([]byte(msg.Payload), receivedObj)
		if err != nil {
			panic(err)
		}

		// The channel may have gotten closed due to the client disconnecting.
		// To not have our Goroutine block or panic, we do the send in a select block.
		// This will jump to the default case if the channel is closed.
		select {
		case ch <- receivedObj: // This is the actual send.
			// Our message went through, do nothing
		default: // This is run when our send does not work.
			log.Println("Channel closed.")
			// You can handle any deregistration of the channel here.
			return // We'll just return ending the routine.
		}
	}
}

//PublishComment publishes comment to redis channel named "rChannel"
func PublishComment(rdb *redis.Client, ctx context.Context, comment *gmodel.Comment) error {
	jsonBytes, err := json.Marshal(comment)
	if err != nil {
		panic(err)
	}
	jsonStr := string(jsonBytes)
	err = rdb.Publish(ctx, rChannel, jsonStr).Err()
	if err != nil {
		return err
	}
	return nil
}
