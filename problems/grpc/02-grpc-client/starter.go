package main

import "context"

// CallSayHello は GreeterClient を使って SayHello を呼び出し、
// レスポンスのメッセージ文字列を返します。
func CallSayHello(ctx context.Context, client GreeterClient, name string) (string, error) {
	// TODO:
	// 1. client.SayHello(ctx, &HelloRequest{Name: name}) を呼び出す
	// 2. エラーが返ってきたら、"call SayHello: " を前置してラップして返す
	//    （fmt.Errorf の %w を使う）
	// 3. 成功したら、レスポンスのメッセージ文字列と nil を返す
	return "", nil
}

func main() {}
