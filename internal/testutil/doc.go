// Package testutil provides testing utilities, including mock implementations
// of various interfaces (e.g., cache.Cache, firestore.KV) used by other packages.
//
// These mocks allow developers to test code that depends on these interfaces
// without requiring real services or complex setups. Keep this package focused
// on test-related functionality, ensuring production code remains free of
// testing scaffolding.
package testutil
