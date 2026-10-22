package service

import (
	"crypto/sha256"
	"encoding/hex"
)

// DeleteUser TODO: проверка: выполнить только если id == callerID || callerRole == admin
// password validate: if err == bcrypt.ErrPasswordTooLong (>72byte),

func hashToken(token string) string {
    h := sha256.Sum256([]byte(token))
    return hex.EncodeToString(h[:])
}