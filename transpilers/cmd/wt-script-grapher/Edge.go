package main

type Edge interface {
  Write(indent string) string
}
