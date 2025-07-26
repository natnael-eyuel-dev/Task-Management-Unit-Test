package mock_repositories

// imports
import (
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Domain"
	"github.com/stretchr/testify/mock"
)

// mocks mongo.SingleResult behavior for Decode
type MockSingleResult struct {
	mock.Mock
	Result interface{}
	Err    error
}

// simulates decoding a result
func (m *MockSingleResult) Decode(v interface{}) error {
	
	// check if there is an error
	if m.Err != nil {
		return m.Err
	}

	// attempt type assertion and copy values
	switch out := v.(type) {
	case *interface{}:
		*out = m.Result
	case *map[string]interface{}:
		*out = *(m.Result.(*map[string]interface{}))
	default:
		// for domain types use reflection or direct assertion
		switch typed := m.Result.(type) {
		case *domain.User:
			*out.(*domain.User) = *typed
		case *domain.Task:
			*out.(*domain.Task) = *typed
		}
	}

	return nil
}
