// Package store provides a key-value storage abstraction through the Store interface.
// It may use various backends to store data, such as Firestore, in-memory implementations,
// or other systems. By depending on interfaces rather than concrete types, you can easily
// swap implementations or mock the store in tests.
//
// This package focuses on production code. For testing without relying on external services,
// use mocks from the testutil package. For example, testutil.NewMockKV() can simulate Firestore KV behavior in-memory.
package store
