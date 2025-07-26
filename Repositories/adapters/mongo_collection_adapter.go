package adapters

// imports
import (
	"context"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// an adapter for the mongo.Collection type
type MongoCollectionAdapter struct {
	Collection *mongo.Collection
}

// this inserts a single document into the collection
func (m *MongoCollectionAdapter) InsertOne(ctx context.Context, doc interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return m.Collection.InsertOne(ctx, doc, opts...)
}

// this returns a cursor for the documents that match the filter
func (m *MongoCollectionAdapter) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return m.Collection.Find(ctx, filter, opts...)
}

// this retrieves a single document from the collection that matches the filter
func (m *MongoCollectionAdapter) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) domain.SingleResult {
	result := m.Collection.FindOne(ctx, filter, opts...)
	return &MongoSingleResultAdapter{Result: result}
}

// this updates a single document in the collection that matches the filter
func (m *MongoCollectionAdapter) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) domain.SingleResult {
	result := m.Collection.FindOneAndUpdate(ctx, filter, update, opts...)
	return &MongoSingleResultAdapter{Result: result}
}

// this deletes a single document from the collection that matches the filter
func (m *MongoCollectionAdapter) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return m.Collection.DeleteOne(ctx, filter, opts...)
}

// this returns the count of documents in the collection that match the filter
func (a *MongoCollectionAdapter) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return a.Collection.CountDocuments(ctx, filter, opts...)
}




