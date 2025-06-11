package auth

import (
    "os"
    "time"
    "log"
    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
)

type Claims struct {
    UserID string `json:"user_id"`
    jwt.RegisteredClaims
}

func getJWTSecret() []byte {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        log.Fatal("JWT_SECRET is not set in .env")
    }
    return []byte(secret)
}

func getJWTExpiryDuration() time.Duration {
    durationStr := os.Getenv("JWT_EXPIRES_IN")
    if durationStr == "" {
        durationStr = "24h"
    }
    duration, err := time.ParseDuration(durationStr)
    if err != nil {
        log.Fatalf("Invalid JWT_EXPIRES_IN format: %v", err)
    }
    return duration
}

func GenerateJWT(userID uuid.UUID) (string, error) {
    claims := Claims{
        UserID: userID.String(),
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(getJWTSecret())
}

func ParseJWT(tokenStr string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
        return getJWTSecret(), nil
    })

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, err
}
