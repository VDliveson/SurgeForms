package dependencies

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/VDliveson/SurgeForms/go-producer/internal/dependencies/cache"
	"github.com/VDliveson/SurgeForms/go-producer/internal/dependencies/database"
	"github.com/VDliveson/SurgeForms/go-producer/internal/dependencies/queue"
)

type AppDependencies struct {
	Ctx   context.Context
	Cache cache.CacheInterface
	DB    *database.DB
	Queue *queue.Queue
}

// New initializes dependencies and injects them into the App struct
func New(ctx context.Context) (*AppDependencies, error) {
	app := &AppDependencies{Ctx: ctx}

	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Connecting to RabbitMQ...")
		q, err := queue.ConnectQueue(ctx)
		if err != nil {
			errChan <- fmt.Errorf("error connecting to RabbitMQ: %v", err)
			return
		}
		app.Queue = q
	}()

	// Initialize MongoDB
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Connecting to MongoDB...")
		db, err := database.ConnectDB(ctx)
		if err != nil {
			errChan <- fmt.Errorf("error connecting to MongoDB: %v", err)
			return
		}
		app.DB = db
	}()

	// Initialize Redis
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Connecting to Redis...")
		cache, err := cache.ConnectCache(ctx, "redis")
		if err != nil {
			errChan <- fmt.Errorf("error connecting to Redis: %v", err)
			return
		}
		app.Cache = cache
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	log.Println("All services connected successfully.")
	return app, nil
}

// Shutdown cleans up resources
func (a *AppDependencies) Shutdown() {
	log.Println("Shutting down services...")

	if a.Queue != nil {
		a.Queue.Channel.Close()
		a.Queue.Connection.Close()
		log.Println("RabbitMQ connection closed.")
	}

	if a.DB != nil {
		a.DB.Client.Disconnect(a.Ctx)
		log.Println("MongoDB connection closed.")
	}

	if a.Cache != nil {
		a.Cache.Close()
		log.Println("Redis connection closed.")
	}
}
