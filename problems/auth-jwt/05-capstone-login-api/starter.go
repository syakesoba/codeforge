package main

import (
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// jwtSecret はトークンの署名に使う秘密鍵です。
const jwtSecret = "codeforge-capstone-secret"

// User は登録済みユーザーを表します。
type User struct {
	Username     string
	PasswordHash string
}

var (
	mu    sync.Mutex
	users = map[string]User{}
)

// --- ここから下のヘルパーはすべて実装済みです ---

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CheckPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func IssueToken(secret, username string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(ttl).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func VerifyToken(secret, tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	return claims.GetSubject()
}

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		tokenString, ok := strings.CutPrefix(authHeader, "Bearer ")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		username, err := VerifyToken(secret, tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Set("username", username)
		c.Next()
	}
}

// --- ここから下があなたの実装範囲です ---

// registerHandler はユーザーを登録します。
// 成功: 201 {"username": "..."} / バリデーションエラー: 400 / ユーザー名重複: 409
func registerHandler(c *gin.Context) {
	// TODO:
	// 1. credentialsRequest 型を定義する
	//    Username: binding:"required"
	//    Password: binding:"required,min=8"
	// 2. c.ShouldBindJSON でバインドし、失敗したら400を返す
	// 3. HashPassword でハッシュ化する（失敗したら500）
	// 4. mu.Lock() / defer mu.Unlock() で保護しつつ、
	//    既に同じusernameが存在すれば409を返す
	// 5. users に登録し、201で {"username": ...} を返す
}

// loginHandler は認証してJWTを発行します。
// 成功: 200 {"token": "..."} / バリデーションエラー: 400 / 認証失敗: 401
func loginHandler(c *gin.Context) {
	// TODO:
	// 1. リクエストボディをバインドする（失敗したら400）
	// 2. users からユーザーを探す
	// 3. ユーザーが存在しない、またはパスワードが一致しない場合は
	//    どちらも同じ 401 {"error": "invalid credentials"} を返す
	// 4. IssueToken(jwtSecret, username, time.Hour) でトークンを発行し、
	//    200 で {"token": ...} を返す
}

// profileHandler は認証済みユーザーの情報を返します（200）。
func profileHandler(c *gin.Context) {
	// TODO: c.GetString("username") で取得し、200で {"username": ...} を返す
}

// newRouter はこのレッスンで使うルーティングを登録したルーターを返します。
func newRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/register", registerHandler)
	r.POST("/login", loginHandler)
	r.GET("/profile", AuthMiddleware(jwtSecret), profileHandler)

	return r
}

func main() {
	newRouter().Run(":8080")
}
