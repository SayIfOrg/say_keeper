package commenting

import (
	"context"
	"fmt"
	"github.com/SayIfOrg/say_keeper/utils"
	"github.com/redis/go-redis/v9"
	"log"
	"sync"
)

//Subs keeps track of subscribers to comments
type Subs[T interface{}] struct {
	Subs []chan *T
	Mu   sync.Mutex
}

const rChannel string = "commentingChan"

//SubscribeComment subscribes to a redis channel named `rChannel` for new comments
// and callback the subscribers throu `Subs`
func SubscribeComment[T interface{}](ctx context.Context, rdb *redis.Client, subs *Subs[T], unmarshalFn func([]byte) (*T, error)) {
	pubsub := rdb.Subscribe(ctx, rChannel)
	defer func() {
		err := pubsub.Close()
		if err != nil {
			panic(err)
		}
		log.Println("redis pubsub closed")
	}()

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Println(msg.Channel, msg.Payload)

		// Deserialize message from JSON string into object
		receivedObj, err := unmarshalFn([]byte(msg.Payload))
		if err != nil {
			panic(err)
		}

		subs.Mu.Lock()
		var invalidSubsIndex []int
		for i, ch := range subs.Subs {
			// The channel may have gotten closed due to the client disconnecting.
			// To not have our Goroutine block or panic, we do the send in a select block.
			// This will jump to the default case if the channel is closed.
			select {
			case ch <- receivedObj: // This is the actual send.
				// Our message went through, do nothing
			default: // This is run when our send does not work.
				log.Println("a closed graph chan detected")
				invalidSubsIndex = append(invalidSubsIndex, i)
				// You can handle any deregistration of the channel here.
			}
		}
		//remove the closed graphql channel from the list
		subs.Subs = utils.RemoveByIndexes(subs.Subs, invalidSubsIndex)
		subs.Mu.Unlock()
	}
}

//PublishComment publishes comment to redis channel named "rChannel"
func PublishComment(rdb *redis.Client, ctx context.Context, jsComment []byte) error {
	jsonStr := string(jsComment)
	err := rdb.Publish(ctx, rChannel, jsonStr).Err()
	if err != nil {
		return err
	}
	return nil
}
