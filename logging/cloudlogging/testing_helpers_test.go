package cloudlogging

import (
	"context"
	"sync"
)

// mockProjectIDFetcher allows simulating responses from the metadata service
// for deterministic testing of project ID determination.
type mockProjectIDFetcher struct {
	id  string // Project ID to return
	err error  // Error to return, if any
}

// ProjectID implements the projectIDFetcher interface for the mock.
func (m *mockProjectIDFetcher) ProjectID(ctx context.Context) (string, error) {
	return m.id, m.err
}

// resetDetermineProjectID resets the sync.Once and cached project ID,
// ensuring a clean state for each test run that manipulates this global state.
func resetDetermineProjectID() {
	projectIDOnce = sync.Once{}
	determinedProjectID = ""
}
