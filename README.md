# cadence-api

A simple REST api using mongo and built with a DDD style controller-service-repo pattern. 

Testing utilizes a more "idiomatic" Ginkgo approach as a personal PoC towards the viability of the package in terms of balancing readability and DRY principles. Not entirely sure if it was worth it as I think I lean towards keeping tests DAMP.

### Run

requires mongodb to be running on localhost:27017

```bash
make run
```

### Test

```bash
make test
```



