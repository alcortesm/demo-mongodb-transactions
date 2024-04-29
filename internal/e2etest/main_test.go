package e2etest

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

var (
	// the URI of the MongoDB Docker container
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

	mongodbContainer, err := mongodb.RunContainer(
		ctx,
		testcontainers.WithImage("mongo:6.0.15"),
		withReplicaSet(),
	)
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

// withReplicaSet configures a MongoDB testcontainer to start with a replica set named "rs".
func withReplicaSet() testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) {
		req.Cmd = append(req.Cmd, "--replSet", "rs")

		hook := testcontainers.ContainerLifecycleHooks{
			PostReadies: []testcontainers.ContainerHook{
				func(ctx context.Context, c testcontainers.Container) error {
					cIP, err := c.ContainerIP(ctx)
					if err != nil {
						return err
					}

					cmd := eval("rs.initiate({ _id: 'rs', members: [ { _id: 0, host: '%s:27017' } ] })", cIP)

					if exitCode, _, err := c.Exec(ctx, cmd); err != nil || exitCode != 0 {
						return fmt.Errorf("failed to initiate the replica set with status %d: %s", exitCode, err)
					}

					return nil
				},
			},
		}
		req.LifecycleHooks = append(req.LifecycleHooks, hook)
	}
}

// eval builds an mongosh|mongo eval command.
func eval(command string, args ...any) []string {
	command = "\"" + fmt.Sprintf(command, args...) + "\""

	return []string{
		"sh",
		"-c",
		// In previous versions, the binary "mongosh" was named "mongo".
		"mongosh --quiet --eval " + command + " || mongo --quiet --eval " + command,
	}
}
