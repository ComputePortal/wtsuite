package main

type Node interface {
  Name() string // names be unique
  Write(indent string) string
}

