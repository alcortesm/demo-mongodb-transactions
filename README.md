# Intro

The purpose of this project is to explore how MongoDB transactions can be used
to ensure consistency of Domain Drive Design (DDD) aggregates.

# Run the demo

The `internal/e2etest.Test_Concurrency_AddLotsOfUsersConcurrentlyToGroup` test
verifies Mongo transactions allows to keep DDD aggregates consistent
when facing concurrent calls in our application layer.

The easiest why to run the test is running all the test in the repo:

```
; go test ./...
ok      github.com/alcortesm/demo-mongodb-transactions/internal/application     (cached)
?       github.com/alcortesm/demo-mongodb-transactions/internal/testhelp        [no test files]
ok      github.com/alcortesm/demo-mongodb-transactions/internal/domain  (cached)
ok      github.com/alcortesm/demo-mongodb-transactions/internal/e2etest (cached)
ok      github.com/alcortesm/demo-mongodb-transactions/internal/infra/mongo     (cached)
```
