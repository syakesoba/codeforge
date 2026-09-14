//go:build ignore

package main

import (
	"context"
	"fmt"
)

func CallSayHello(ctx context.Context, client GreeterClient, name string) (string, error) {
	resp, err := client.SayHello(ctx, &HelloRequest{Name: name})
	if err != nil {
		return "", fmt.Errorf("call SayHello: %w", err)
	}
	return resp.GetMessage(), nil
}

func main() {}
