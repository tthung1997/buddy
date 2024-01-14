package bgg

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tthung1997/buddy/core/bgg"
)

func TestClient_GetThing_Success(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Respond with a successful response
		rw.Write([]byte(`<item id="1"></item>`))
	}))
	defer server.Close()

	// Create a Client pointing to the mock server
	client := bgg.NewClient(
		bgg.ClientConfig{
			Root:                server.URL,
			MaxRetries:          1,
			RetryDelayInSeconds: 1,
		},
	)

	// Call GetThing
	thing, err := client.GetThing("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test that the thing is what we expect
	if thing != `<item id="1"></item>` {
		t.Errorf("unexpected thing: %v", thing)
	}
}

func TestClient_GetThing_ServerError(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Respond with a server error
		rw.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Create a Client pointing to the mock server
	client := bgg.NewClient(
		bgg.ClientConfig{
			Root:                server.URL,
			MaxRetries:          1,
			RetryDelayInSeconds: 1,
		},
	)

	// Call GetThing
	_, err := client.GetThing("1")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	// Test that the error is what we expect
	if err.Error() != "unexpected response code: 500" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestClient_GetUser_Success(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Respond with a successful response
		rw.Write([]byte(`<user id="1"></user>`))
	}))
	defer server.Close()

	// Create a Client pointing to the mock server
	client := bgg.NewClient(
		bgg.ClientConfig{
			Root:                server.URL,
			MaxRetries:          1,
			RetryDelayInSeconds: 1,
		},
	)

	// Call GetUser
	user, err := client.GetUser("testuser")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test that the user is what we expect
	if user.Id != "1" {
		t.Errorf("unexpected user: %v", user)
	}
}

func TestClient_GetUser_NotFound(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Respond with a user not found response
		rw.Write([]byte(`<user></user>`))
	}))
	defer server.Close()

	// Create a Client pointing to the mock server
	client := bgg.NewClient(
		bgg.ClientConfig{
			Root:                server.URL,
			MaxRetries:          1,
			RetryDelayInSeconds: 1,
		},
	)

	// Call GetUser
	_, err := client.GetUser("testuser")
	if !errors.Is(err, bgg.ErrUserNotFound) {
		t.Fatalf("expected user not found error, got: %v", err)
	}
}

func TestClient_GetCollection_Success(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Respond with a successful response
		rw.Write([]byte(`<items><item objectid="1"></item></items>`))
	}))
	defer server.Close()

	// Create a Client pointing to the mock server
	client := bgg.NewClient(
		bgg.ClientConfig{
			Root:                server.URL,
			MaxRetries:          1,
			RetryDelayInSeconds: 1,
		},
	)

	// Call GetCollection
	collection, err := client.GetCollection(bgg.CollectionFilter{Username: "testuser"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test that the collection is what we expect
	if len(collection.Items) != 1 || collection.Items[0].Id != "1" {
		t.Errorf("unexpected collection: %v", collection)
	}
}

func TestClient_GetCollection_NotFound(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Respond with a collection not found response
		rw.Write([]byte(`<items></items>`))
	}))
	defer server.Close()

	// Create a Client pointing to the mock server
	client := bgg.NewClient(
		bgg.ClientConfig{
			Root:                server.URL,
			MaxRetries:          1,
			RetryDelayInSeconds: 1,
		},
	)

	// Call GetCollection
	collection, err := client.GetCollection(bgg.CollectionFilter{Username: "testuser"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test that the collection is empty
	if len(collection.Items) != 0 {
		t.Errorf("expected empty collection, got: %v", collection)
	}
}

func TestClient_GetCollection_Retry(t *testing.T) {
	attempts := 0

	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		attempts++

		if attempts == 1 {
			// First attempt, respond with StatusAccepted
			rw.WriteHeader(http.StatusAccepted)
		} else {
			// Second attempt, respond with StatusOK and a valid collection
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(`<items><item objectid="1"></item></items>`))
		}
	}))
	defer server.Close()

	// Create a Client pointing to the mock server
	client := bgg.NewClient(
		bgg.ClientConfig{
			Root:                server.URL,
			MaxRetries:          2,
			RetryDelayInSeconds: 1,
		},
	)

	// Call GetCollection
	collection, err := client.GetCollection(bgg.CollectionFilter{Username: "testuser"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test that the collection is what we expect
	if len(collection.Items) != 1 || collection.Items[0].Id != "1" {
		t.Errorf("unexpected collection: %v", collection)
	}
}

func TestClient_GetCollection_MaxRetriesExceeded(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Always respond with StatusAccepted
		rw.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	// Create a Client pointing to the mock server
	client := bgg.NewClient(
		bgg.ClientConfig{
			Root:                server.URL,
			MaxRetries:          2,
			RetryDelayInSeconds: 1,
		},
	)

	// Call GetCollection
	_, err := client.GetCollection(bgg.CollectionFilter{Username: "testuser"})
	if err == nil || err.Error() != "max retries exceeded: 2" {
		t.Fatalf("expected max retries exceeded error, got: %v", err)
	}
}

func TestClient_GetCollection_UnexpectedResponseCode(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Always respond with StatusBadRequest
		rw.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	// Create a Client pointing to the mock server
	client := bgg.NewClient(
		bgg.ClientConfig{
			Root:                server.URL,
			MaxRetries:          2,
			RetryDelayInSeconds: 1,
		},
	)

	// Call GetCollection
	_, err := client.GetCollection(bgg.CollectionFilter{Username: "testuser"})
	if err == nil || err.Error() != "unexpected response code: 400" {
		t.Fatalf("expected unexpected response code error, got: %v", err)
	}
}
