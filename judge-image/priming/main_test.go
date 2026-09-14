package main

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// このテストは、レッスンの採点で実際に使われるパッケージ
// （testing / httptest / gin のルーター実行経路）のビルドキャッシュを
// イメージに焼き込むためだけに存在する。
// これが無いと、Ginを使うレッスンの「初回提出」だけテストバイナリの
// フルコンパイルが走り、タイムアウトしてしまう。
func TestPriming(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("priming request failed: %d", rec.Code)
	}
}

// echoServiceDesc は protoc生成コードを使わず手作業で定義した最小のgRPC
// サービス記述で、プライミング専用。実際のレッスンではprotoc生成の
// *_grpc.pb.go を使うが、ここではgrpcランタイム・protobufランタイム・
// bufconnによるテスト経路のビルド/リンクキャッシュを焼き込めれば十分なため
// 生成コードは不要。
var echoServiceDesc = grpc.ServiceDesc{
	ServiceName: "priming.Echo",
	HandlerType: (*any)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler: func(_ any, ctx context.Context, dec func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
				req := new(wrapperspb.StringValue)
				if err := dec(req); err != nil {
					return nil, err
				}
				return wrapperspb.String(req.GetValue()), nil
			},
		},
	},
	Metadata: "priming.proto",
}

// TestPrimingGRPC は、gRPCを使うレッスンの採点で実際に使われるパッケージ
// （grpc / bufconn / protobufランタイム）のビルドキャッシュを焼き込むため
// だけに存在する（Ginの場合と同じ理由）。
func TestPrimingGRPC(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	srv := grpc.NewServer()
	srv.RegisterService(&echoServiceDesc, nil)
	go func() { _ = srv.Serve(lis) }()
	defer srv.Stop()

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	req := wrapperspb.String("hello")
	resp := new(wrapperspb.StringValue)
	if err := conn.Invoke(context.Background(), "/priming.Echo/Echo", req, resp); err != nil {
		t.Fatalf("invoke failed: %v", err)
	}
	if !proto.Equal(req, resp) {
		t.Fatalf("unexpected response: %q", resp.GetValue())
	}
}
