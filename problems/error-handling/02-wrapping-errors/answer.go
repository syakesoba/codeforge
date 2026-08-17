//go:build ignore

package main

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

var users = map[int]string{
	1: "Alice",
	2: "Bob",
}

func FindUser(id int) (string, error) {
	name, ok := users[id]
	if !ok {
		return "", fmt.Errorf("user %d: %w", id, ErrNotFound)
	}
	return name, nil
}

func main() {}
