// Command priming はジャッジ用Dockerイメージのビルド時にだけ使う。
// Gin/GORM/SQLite/JWT/bcryptをここでimportしてビルドすることで、
// これらのモジュールをイメージのGOMODCACHE/GOCACHEに焼き込む
// （実行時は --network none でオフラインになるため、事前に用意しておく必要がある）。
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

var _ = gin.New
var _ = gorm.Open
var _ = sqlite.Open
var _ = jwt.New
var _ = bcrypt.GenerateFromPassword
var _ = grpc.NewServer
var _ = codes.OK
var _ = insecure.NewCredentials
var _ = metadata.New
var _ = status.New
var _ = bufconn.Listen
var _ = proto.Marshal

func main() {}
