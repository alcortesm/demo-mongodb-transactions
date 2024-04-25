package e2etest

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

var (
	mongoURI string
)

// TestMain performs some setup/cleanup before/after running the tests in this package.
//
// Setup:
//   - starts a MongoDB Docker container and fills mongoURI with its connection
//     string.
//
// Cleanup:
//   - terminate the MongoDB Docker container
func TestMain(m *testing.M) {
	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	mongodbContainer, err := mongodb.RunContainer(ctx, testcontainers.WithImage("mongo:6.0.15"))
	if err != nil {
		log.Fatalf("starting MongoDB container: %v", err)
	}

	// clean up the mongo container
	defer func() {
		timeout := 2 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := mongodbContainer.Terminate(ctx); err != nil {
			log.Fatalf("terminating container: %v", err)
		}
	}()

	mongoURI, err = mongodbContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("getting MongoDB connection string: %v", err)
	}

	os.Exit(m.Run())
}
