package mock_repositories

// imports
import (
    "context"
    "github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Domain"
    "github.com/stretchr/testify/mock"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// mock collection for testing
type MockCollection struct {
    mock.Mock
}

// mocks InsertOne method of the collection
func (m *MockCollection) InsertOne(contx context.Context, doc interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
    args := m.Called(contx, doc)
    res := args.Get(0)
    if res == nil {
        return nil, args.Error(1)
    }
    return res.(*mongo.InsertOneResult), args.Error(1)
}

// mocks Find method of the collection
func (m *MockCollection) Find(contx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
    args := m.Called(contx, filter)
    return args.Get(0).(*mongo.Cursor), args.Error(1)
}

// mocks FindOne method of the collection
func (m *MockCollection) FindOne(contx context.Context, filter interface{}, opts ...*options.FindOneOptions) domain.SingleResult {
    args := m.Called(contx, filter)
    return args.Get(0).(domain.SingleResult)
}

// mocks FindOneAndUpdate method of the collection
func (m *MockCollection) FindOneAndUpdate(contx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) domain.SingleResult {
    args := m.Called(contx, filter, update)
    return args.Get(0).(domain.SingleResult)
}

// mocks DeleteOne method of the collection
func (m *MockCollection) DeleteOne(contx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
    args := m.Called(contx, filter)
    return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

// mocks CountDocuments method of the collection
func (m *MockCollection) CountDocuments(contx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
    args := m.Called(contx, filter)
    return args.Get(0).(int64), args.Error(1)
}