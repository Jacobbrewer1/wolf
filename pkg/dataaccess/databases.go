package dataaccess

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDB is the Mongo client. This is a connection pool.
var MongoDB *mongo.Client

const mongoDatabase = "wolf"
