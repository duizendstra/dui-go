// Package firestore provides a key-value abstraction on top of Google Cloud Firestore.
//
// It defines a KV interface representing simple key-value operations. The FirestoreKV
// type is a concrete implementation of KV that uses a Firestore collection. If a key does
// not exist, Get returns an empty string without error. Set writes values and Close
// releases Firestore resources.
//
// This package does not include mocks for testing. To test code that depends on
// FirestoreKV without connecting to a real Firestore instance, use the MockFirestoreKV
// from the testutil package.
package firestore
