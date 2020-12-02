package main

//go:generate go-enum -f=$GOFILE --names

// Platform represents
// ENUM(
// _ // invalid default int value
// apple
// android
//)
type Platform int

// Status represents
// ENUM(
// _ // invalid default int value
// active
// inactive
//)
type Status int
