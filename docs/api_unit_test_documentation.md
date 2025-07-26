# Task Management API Documentation (Unit Test)

## Overview
This project implements a comprehensive unit test suite for a Task Management API backend. The suite covers all major layers: **Repositories**, **Usecases**, **Controllers**, **Routers**, and **Infrastructure**. The tests are written using [Testify](https://github.com/stretchr/testify) for assertions, mocking, and test suites, and are organized by feature and layer.

---

## Folder Structure

- **Delivery/controllers/**  
  Tests for HTTP controllers (user and task endpoints).
- **Delivery/routers/**  
  Tests for HTTP routers (user and task routing endpoints).
- **Infrastructure/**  
  Tests for services like JWT and password hashing.
- **Repositories/**  
  Tests for MongoDB-based repositories (task and user).
- **Usecases/**  
  Tests for business logic (usecases) for tasks and users.

---

## Running the Tests

### Prerequisites

- Go 1.18+
- All dependencies installed (`go mod tidy`)
- [Testify](https://github.com/stretchr/testify) installed

### Command

To run all tests in the project, open a terminal in the project root and execute:

```sh
go test ./... -v
```

- `-v` enables verbose output.
- `./...` recursively runs tests in all subfolders.

To run tests in a specific folder (e.g., user repository):

```sh
go test ./Repositories -v
```

Or for a specific file:

```sh
go test ./Repositories/user_repository_test.go -v
```

---

## Test Coverage

### How to Measure

To generate a coverage report:

```sh
go test ./... -cover
```

For an HTML coverage report:

```sh
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Coverage Summary

- **Repositories:**  
  - CRUD operations for tasks and users are tested, including edge cases (not found, duplicate, invalid input).
  - Mock collections simulate MongoDB behavior.
  - Coverage: ~73% (all main methods, error branches, and edge cases).

- **Usecases:**  
  - Business logic for registration, login, promotion, task creation, update, and deletion.
  - Coverage: ~90% (all branches, including error handling and validation).

- **Controllers:**  
  - HTTP endpoints for user and task management.
  - Coverage: ~85% (all endpoints, status codes, and error responses).

- **Infrastructure:**  
  - JWT and password services.
  - Coverage: ~93% (token generation, validation, password hashing, and checking).

- **Routers:**  
  - Routing logic for user and task endpoints is tested, including authentication, authorization, and correct mapping to controllers.
  - Tests cover public, authenticated, and admin-only routes, verifying access control and integration with usecase and JWT services.
  - Coverage: ~100% (all main routes, middleware, and error branches).

---

## Test Suite Details

### Repositories

- **TaskRepositoryTestSuite**
  - Tests for creating, updating, deleting, and fetching tasks.
  - Edge cases: invalid IDs, not found, no fields provided for update.
  - Uses `MockCollection` for simulating MongoDB operations.

- **UserRepositoryTestSuite**
  - Tests for user creation, duplicate users, fetching by username/ID, counting users, and role updates.
  - Edge cases: user not found, duplicate, error propagation.
  - Uses `MockCollection` and `MockSingleResult`.

### Usecases

- **TaskUseCaseTestSuite**
  - Tests for business logic around tasks: creation, update, deletion, validation.
  - Edge cases: invalid due dates, invalid status, not found.

- **UserUseCaseTestSuite**
  - Tests for registration (first user becomes admin), login (success, invalid password, user not found), promotion to admin.
  - Edge cases: invalid user ID, user not found, password validation.

### Controllers

- **TaskControllerTestSuite**
  - Tests for HTTP endpoints: create, get all, get by ID, update, delete.
  - Verifies status codes, response bodies, and error handling.

- **UserControllerTestSuite**
  - Tests for registration, login, promotion endpoints.
  - Verifies status codes, response bodies, and error handling.

### Routers 

- **RouterTestSuite**
  - Tests for HTTP routers, focusing on correct routing, middleware application, and integration with controllers.
  - Verifies that public routes (e.g., registration, login) are accessible without authentication.
  - Tests authenticated routes (e.g., getting tasks) for proper token validation and access control.
  - Covers admin-only routes (e.g., creating, updating, deleting tasks, promoting users) to ensure only users with admin role can access them.
  - Includes edge cases such as unauthorized access, forbidden actions for non-admin users, and invalid input handling.
  - Uses mocks for usecase and JWT services to isolate routing logic from business and infrastructure layers.

### Infrastructure

- **PasswordServiceTestSuite**
  - Tests for password hashing, checking, and length limits.
  - Edge cases: empty password, wrong password.

- **JWTServiceTestSuite**
  - Tests for token generation, validation, expiration, and error cases (empty secret, invalid claims).

---

## Mocking Strategy

- **Repositories/mocks/**  
  Custom mocks for MongoDB collections and single results.
- **Usecases/mocks/**  
  Mocks for usecase interfaces.
- **Infrastructure/mocks/**  
  Mocks for JWT and password services.

Mocks are used to isolate tests from external dependencies and focus on logic and error handling.

---

## Issues Encountered

### 1. **Error Type Mismatches**
- Some tests failed due to incorrect error types being returned (e.g., returning `ErrTaskNotFound` instead of `ErrUserNotFound`).  
  **Resolution:** Standardized error returns in repository implementations.

### 2. **Mocking Complex MongoDB Behavior**
- Simulating MongoDB's `SingleResult` and collection methods required custom mock types.
  **Resolution:** Implemented `MockSingleResult` and `MockCollection` with flexible type assertions.

### 3. **Test Data Consistency**
- Ensuring consistent test data (e.g., ObjectIDs, timestamps) across tests.
  **Resolution:** Used fixed values or generated data within setup methods.

### 4. **Coverage Gaps**
- Some controller error branches (e.g., malformed JSON) were not covered.
  **Resolution:** Added additional tests for invalid input cases.

---

## Best Practices Followed

- **Arrange-Act-Assert Pattern:**  
  Each test is structured for clarity and maintainability.
- **Isolation:**  
  Mocks ensure tests do not depend on external systems.
- **Comprehensive Edge Case Testing:**  
  All error branches and edge cases are covered.
- **Consistent Naming:**  
  Test functions are named for easy identification and reporting.
- **Suite-based Organization:**  
  Test suites group related tests for easier setup and teardown.

---

## Conclusion

The unit test suite provides robust coverage for the Task Management API backend. It ensures reliability, correctness, and maintainability across all layers. By following best practices and addressing encountered issues, the suite supports rapid development and confident refactoring.

For any new features or bug fixes, always add corresponding tests and run the suite to verify correctness.

---

**For troubleshooting, check the output panel in VS Code and ensure all dependencies are installed.**

---

**Happy Testing!**
