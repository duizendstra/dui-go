// Package store provides a generic key-value storage abstraction through the Store
// interface. Implementations of this Store interface can utilize various backends.
//
// The package also defines a KV interface (Get, Set, Close). This KV interface
// represents the contract for basic key-value operations that a backend (like one
// based on Firestore) might provide. The FirestoreStore implementation in this
// package, for example, uses an object that satisfies this KV interface.
//
// Hypothetical Store Interface Usage:
//
//	import "github.com/duizendstra/dui-go/store"
//
//	// var myAppStore store.Store = // ... obtain a store.Store implementation (e.g., NewFirestoreStore)
//	//
//	// func processItem(ctx context.Context, itemID string, itemValue string, s store.Store) error {
//	//     if err := s.Set(ctx, itemID, itemValue); err != nil {
//	//         return err
//	//     }
//	//     retrievedValue, err := s.Get(ctx, itemID)
//	//     // ...
//	//     return nil
//	// }
//
// Relationship with firestore.KV:
// The `firestore` package in this library provides `firestore.FirestoreKV`, a concrete
// key-value implementation using Google Cloud Firestore. `firestore.FirestoreKV` can
// serve as an implementation for the `store.KV` interface defined herein, thanks to
// Go's structural typing, if their method signatures are compatible.
// The `store.KV` interface is defined in this consumer package (`store`) to specify
// the dependencies of `FirestoreStore` more explicitly from the `store` package's perspective.
//
// By depending on these interfaces (Store, KV) rather than concrete types, applications
// can more easily swap storage implementations or mock the store and its underlying
// key-value mechanism in tests.
//
// This package focuses on production code. For testing without relying on external
// services, use mocks from the internal testutil package. For example,
// testutil.NewMockFirestoreKV() can simulate Firestore KV behavior in-memory and
// can be used to back a FirestoreStore instance during tests.
package store
