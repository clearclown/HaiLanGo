package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// RSAKeyBits はRSA鍵のビット数
	RSAKeyBits = 2048
	// RefreshTokenLength はリフレッシュトークンの長さ
	RefreshTokenLength = 32
)

var (
	// AccessTokenExpiry はアクセストークンの有効期限（15分）
	AccessTokenExpiry = 15 * time.Minute
	// RefreshTokenExpiry はリフレッシュトークンの有効期限（7日）
	RefreshTokenExpiry = 7 * 24 * time.Hour

	// privateKey はRSA秘密鍵
	privateKey *rsa.PrivateKey
	// publicKey はRSA公開鍵
	publicKey *rsa.PublicKey

	// エラー
	ErrEmptyUserID    = errors.New("ユーザーIDが空です")
	ErrEmptyEmail     = errors.New("メールアドレスが空です")
	ErrInvalidToken   = errors.New("無効なトークンです")
	ErrExpiredToken   = errors.New("トークンの有効期限が切れています")
	ErrTokenNotParsed = errors.New("トークンを解析できません")
)

// Claims はJWTのクレーム情報
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateRSAKeys はRSA鍵ペアを生成する（テスト用または初期化時に使用）
func GenerateRSAKeys() error {
	key, err := rsa.GenerateKey(rand.Reader, RSAKeyBits)
	if err != nil {
		return err
	}

	privateKey = key
	publicKey = &key.PublicKey
	return nil
}

// SetRSAKeys は既存のRSA鍵ペアを設定する
func SetRSAKeys(privKey *rsa.PrivateKey, pubKey *rsa.PublicKey) {
	privateKey = privKey
	publicKey = pubKey
}

// GenerateToken はJWTアクセストークンを生成する
func GenerateToken(userID, email string) (string, error) {
	// バリデーション
	if userID == "" {
		return "", ErrEmptyUserID
	}
	if email == "" {
		return "", ErrEmptyEmail
	}

	// 秘密鍵が設定されていない場合はエラー
	if privateKey == nil {
		return "", errors.New("RSA秘密鍵が設定されていません")
	}

	// クレームの作成
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenExpiry)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	// トークンの生成
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// VerifyToken はJWTトークンを検証してクレームを返す
func VerifyToken(tokenString string) (*Claims, error) {
	// 空のトークンチェック
	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	// 公開鍵が設定されていない場合はエラー
	if publicKey == nil {
		return nil, errors.New("RSA公開鍵が設定されていません")
	}

	// トークンのパース
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 署名方式の確認
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("予期しない署名方式です")
		}
		return publicKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// クレームの取得
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrTokenNotParsed
	}

	return claims, nil
}

// GenerateRefreshToken はリフレッシュトークンを生成する
func GenerateRefreshToken(userID string) (string, time.Time, error) {
	// バリデーション
	if userID == "" {
		return "", time.Time{}, ErrEmptyUserID
	}

	// ランダムなバイト列を生成
	bytes := make([]byte, RefreshTokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", time.Time{}, err
	}

	// Base64エンコード
	token := base64.URLEncoding.EncodeToString(bytes)

	// 有効期限
	expiresAt := time.Now().Add(RefreshTokenExpiry)

	return token, expiresAt, nil
}

// IsTokenExpired はトークンが期限切れかどうかを確認する
func IsTokenExpired(claims *Claims) bool {
	return claims.ExpiresAt.Before(time.Now())
}

// LoadRSAKeysFromFiles はPEM形式のRSA鍵をファイルから読み込んで設定する
// 本番環境では鍵を永続化し、この関数で読み込むことを推奨する
func LoadRSAKeysFromFiles(privateKeyPath, publicKeyPath string) error {
	if strings.TrimSpace(privateKeyPath) == "" || strings.TrimSpace(publicKeyPath) == "" {
		return errors.New("JWT鍵ファイルパスが不正です（JWT_PRIVATE_KEY_PATH と JWT_PUBLIC_KEY_PATH を両方設定してください）")
	}

	privBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("秘密鍵ファイル読み込み失敗: %w", err)
	}

	pubBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return fmt.Errorf("公開鍵ファイル読み込み失敗: %w", err)
	}

	return LoadRSAKeysFromPEM(privBytes, pubBytes)
}

// LoadRSAKeysFromPEM はPEM形式のRSA鍵を読み込んで設定する
func LoadRSAKeysFromPEM(privateKeyPEM, publicKeyPEM []byte) error {
	priv, err := parseRSAPrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return fmt.Errorf("秘密鍵の解析に失敗: %w", err)
	}

	pub, err := parseRSAPublicKeyFromPEM(publicKeyPEM)
	if err != nil {
		return fmt.Errorf("公開鍵の解析に失敗: %w", err)
	}

	SetRSAKeys(priv, pub)
	return nil
}

func parseRSAPrivateKeyFromPEM(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("PEMデータが見つかりません")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("RSA秘密鍵ではありません（PKCS#8）")
		}
		return rsaKey, nil
	default:
		return nil, fmt.Errorf("未対応の秘密鍵タイプ: %s", block.Type)
	}
}

func parseRSAPublicKeyFromPEM(pemBytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("PEMデータが見つかりません")
	}

	switch block.Type {
	case "PUBLIC KEY":
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaPub, ok := pub.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("RSA公開鍵ではありません（PKIX）")
		}
		return rsaPub, nil
	case "RSA PUBLIC KEY":
		return x509.ParsePKCS1PublicKey(block.Bytes)
	case "CERTIFICATE":
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaPub, ok := cert.PublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("証明書の公開鍵がRSAではありません")
		}
		return rsaPub, nil
	default:
		return nil, fmt.Errorf("未対応の公開鍵タイプ: %s", block.Type)
	}
}
