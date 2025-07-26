package adapters

// imports
import (
	"go.mongodb.org/mongo-driver/mongo"
)

// wraps a mongo.SingleResult
type MongoSingleResultAdapter struct {
	Result *mongo.SingleResult
}

// decodes the single result into the provided value
func (m *MongoSingleResultAdapter) Decode(v interface{}) error {
	return m.Result.Decode(v)
}
