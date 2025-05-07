// Package firestore provides a key-value abstraction on top of Google Cloud Firestore.
//
// It defines a KV interface representing simple key-value operations (Get, Set, Close).
// The FirestoreKV type is a concrete implementation of this KV interface, using a
// Firestore collection as the backend. If a key does not exist during a Get operation,
// FirestoreKV returns an empty string without an error, adhering to a common pattern
// for key-value stores where absence is not necessarily an error state.
//
// The Set operation writes or overwrites values, and Close releases underlying
// Firestore client resources.
//
// Relationship with store.KV:
// The `store` package in this library defines its own `store.KV` interface for broader
// storage abstraction purposes. Due to Go's structural typing, `firestore.FirestoreKV`
// can implicitly satisfy `store.KV` if their method signatures match, allowing
// `FirestoreKV` to be used where a `store.KV` is expected.
//
// This package (firestore) focuses on providing the specific Firestore-backed
// key-value store implementation.
//
// For testing code that depends on FirestoreKV without connecting to a real Firestore
// instance, consider using the MockFirestoreKV from the internal testutil package
// (github.com/duizendstra/dui-go/internal/testutil).
package firestore
