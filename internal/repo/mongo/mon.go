package mon

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// this defines our instance, we do this so we only have 1 connection to the database.
// if we create a new connection to the database everytime we want to make request
// it'll be slower because the connection has to be made everytime.
// but if we do this we just make it once and then reuse that connection
var clientInstance *mongo.Client

// Sometimes depending on the configuration you setup the InitConfig thing can be called multiple times.
// so for safety we use this thing, it should stop us from reininting db connections if the InitConfig is run
// multuple times.
var mongoOnce sync.Once

const (
	connectTimeout    = 10
	disconnectTimeout = 5
)

func Client(uri string) *mongo.Client {
	fmt.Println(uri)

	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI'...")
	}

	// so here we tell it to do once, if it's already been done
	// if no connection been made already it'll do it here.
	mongoOnce.Do(func() {
		log.Println("📡 Attempting to connect to database.")
		// all the rest below is mongo related, so need to learn from their docs on how to set this part up
		clientOptions := options.Client().ApplyURI(uri)

		ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
		defer cancel()

		// here we connect to the db
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			// if it failed we fatal here
			log.Fatal(err)
		}

		// double check if it's running
		err = client.Ping(context.Background(), nil)
		if err != nil {
			log.Fatal(err)
		}
		// done
		clientInstance = client
		log.Println("📡 Connected to database.")
	})

	// we just return the already inited instance here
	return clientInstance
}

func DisconnectMongoClient() {
	ctx, cancel := context.WithTimeout(context.Background(), disconnectTimeout*time.Second)
	defer cancel()

	if err := clientInstance.Disconnect(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Database connection closed.")
}
