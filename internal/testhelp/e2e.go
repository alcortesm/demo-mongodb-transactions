package testhelp

import (
	"context"
	"math/rand"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// newTestDatabase is a test helper that creates a mongo database optimized for running
// parallel integration tests:
//
//   - it connects to a Mongo instance at localhost:27017
//   - it uses a unique database name, based on the test name, thus allowing for safe concurrent tests
//   - it drops the database content during test cleanup.
func NewTestDatabase(t *testing.T, uri string) *mongo.Database {
	t.Helper()

	dbName := databaseName(t)
	t.Logf("using MongoDB database name %s", dbName)

	client := newTestClient(t, uri)
	db := client.Database(dbName)
	t.Cleanup(func() {
		t.Helper()

		const timeout = 2 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		t.Cleanup(cancel)

		if err := db.Drop(ctx); err != nil {
			t.Fatalf("cleaning up database after test: %+v", err)
		}
	})

	return db
}

// newTestClient is a test helper that returns a MongoDB client connected to a
// server at uri.
func newTestClient(t *testing.T, uri string) *mongo.Client {
	timeout := 2 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	opts := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		t.Fatalf("connecting to MongoDB at %q: %s", uri, err)
	}

	t.Cleanup(func() {
		timeout := time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err = client.Disconnect(ctx); err != nil {
			t.Fatalf("disconnecting from MongoDB: %v", err)
		}
	})

	// Check to make sure connection is usable
	if err := client.Ping(ctx, nil); err != nil {
		t.Fatalf("pinging: %v", err)
	}

	return client
}

// databaseName is a test helper that generates a MongoDB database name for the
// given test using the following criteria:
//   - trim the test name to 50 bytes (mongo db names cannot be longer than 63 bytes)
//   - while trimming the test name, avoid splitting unicode symbols at the end.
//   - add 13 random bytes to database name to avoid having collisions on similar test names from other packages
//   - replace all forbidden characters in database names by "_"
//
// For example, given the following tests names:
//   - TestFoo
//   - TestSuperVeryLongTestName/scenarioOnALeapYear/scenarioWhenItRains/subTestA
//   - TestSuperVeryLongTestName/scenarioOnALeapYear/scenarioWhenItRains/subTestB
//
// We will return database names like this:
//   - TestFoo_x3la8aXkI2fi
//   - TestSuperVeryLongTestName_scenarioOnALeapYear_scen_t0p8UoMDIafi
//   - TestSuperVeryLongTestName_scenarioOnALeapYear_scen_L1ljD7cT7Wmn
func databaseName(t *testing.T) string {
	t.Helper()

	const (
		// maxDatabaseNameLen is the max number of bytes a database name can have,
		// this is a MongoDB constraint.
		maxDatabaseNameLen = 63
		// randomSuffixLen is the random suffix length added to test database names, including the "_" separator
		randomSuffixLen = 13
	)

	result := t.Name()

	removeForbiddenRunes := func(r rune) rune {
		// forbidenDatabaseNameRunes is the set of runes that are not allowed in a database name,
		// see https://www.mongodb.com/docs/manual/reference/limits/#naming-restrictions
		const forbidenDatabaseNameRunes = `/\. "$*<>:|?`

		if strings.ContainsRune(forbidenDatabaseNameRunes, r) {
			return '_'
		}
		return r
	}
	result = strings.Map(removeForbiddenRunes, result)

	// the (random) suffix to be added to the database name
	randomSuffix := func() string {
		const validChars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

		b := make([]byte, randomSuffixLen-1)
		for i := range b {
			//nolint:gosec // weak random generation is ok here
			b[i] = validChars[rand.Intn(len(validChars))]
		}

		return "_" + string(b)
	}()

	// the max len of the original input we can keep so after adding the suffix
	// we are still under the max database name len for mongo.
	inputLenThreshold := maxDatabaseNameLen - randomSuffixLen
	if len(result) <= inputLenThreshold {
		return result + randomSuffix
	}

	end := 0 // index of the first rune that doesn't fit in 50 bytes
	for i, r := range result {
		nextRuneIndex := i + utf8.RuneLen(r)
		if nextRuneIndex > inputLenThreshold {
			break
		}
		end = nextRuneIndex
	}

	result = result[:end] + randomSuffix

	return result
}
