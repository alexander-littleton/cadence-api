package cadence_errors

import (
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrNotFound = mongo.ErrNoDocuments
var ValidationErr = errors.New("validation failed")
