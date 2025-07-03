package main

import (
	"context"
	"floor_parser/internal/parser"
	"floor_parser/internal/repository"
	"floor_parser/internal/usecases"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	DB_URL := os.Getenv("DATABASE_URL")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()


	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	pool, err := pgxpool.New(ctx, DB_URL)
	if err != nil {
		log.Default().Println("Error while connect to the pool")
	}
	defer pool.Close()

	r := gin.Default()

	repository := repository.NewRepository(pool)

	r.GET("/getFloorByCollectionName", func(c *gin.Context) {
		collName := c.Query("collectionName")
		floorPrice, err := usecases.GetFloorPrice(c, repository, collName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"floor": floorPrice}) 
	})

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func()  {
		defer wg.Done()
		for {
			select {
			case <- ctx.Done():
				log.Default().Println("Goroutine stopping")
				return
			default: 
				collections, err := parser.ParseCollections()
				if err != nil {
					log.Default().Printf("Failed while parsing collections, err: %s", err)
					select {
					case <- ctx.Done():
						log.Println("Parser goroutine stopped (during sleep after error)")
						return
					case <- time.After(20 * time.Second):
						continue
					}
				} 
				for _, collection := range collections {
					id, err := usecases.InsertPrice(ctx, repository, collection.Name, collection.Floor)
					if err != nil {
						log.Printf("Error while insert price in the DB, err: %s", err)
						continue
					}
					log.Printf("Insert success, id:%d", id)
				}
				select {
				case <- ctx.Done():
					log.Println("Parser goroutine stopped (during sleep after error)")
					return
				case <- time.After(20 * time.Second):
					continue
				}
			}	
		}
	}()

	server := &http.Server{
		Addr: ":8080",
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error while running server, err:%s", err)
		}
	}()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	<-sigCh

	cancel()
	wg.Wait()
	log.Default().Print("Connection closed")
}