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
	"gorm.io/gorm"
)

var _ = gin.New
var _ = gorm.Open
var _ = sqlite.Open
var _ = jwt.New
var _ = bcrypt.GenerateFromPassword

func main() {}
