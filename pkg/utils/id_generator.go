package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"github.com/google/uuid"
)

func GenerateID() string {
	return uuid.New().String() // universal, semua versi kompatibel
}

// ==========================
// CONFIG
// ==========================
const secretKey = "super-secret-key-CHANGE-ME"

// ==========================
// SIGN ID
// ==========================
func SignID(id int64, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(strconv.FormatInt(id, 10)))
	signature := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("%d.%s", id, signature)
}

// ==========================
// VERIFY ID
// ==========================
func VerifyID(signed string, secret string) (int64, bool) {
	parts := strings.Split(signed, ".")
	if len(parts) != 2 {
		return 0, false
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, false
	}

	expected := SignID(id, secret)

	// prevent timing attack
	if !hmac.Equal([]byte(expected), []byte(signed)) {
		return 0, false
	}

	return id, true
}

// ==========================
// MAIN
// ==========================
// func main() {
	// contoh BIGINT ID dari database
// 	var originalID int64 = 123456789

// 	fmt.Println("Original ID :", originalID)

// 	// SIGN
// 	signedID := SignID(originalID, secretKey)
// 	fmt.Println("Signed ID   :", signedID)

// 	// VERIFY (VALID)
// 	id, ok := VerifyID(signedID, secretKey)
// 	fmt.Println("\nVerify valid ID:")
// 	fmt.Println("Valid :", ok)
// 	fmt.Println("ID    :", id)

// 	// VERIFY (TAMPERED)
// 	tampered := strings.Replace(signedID, "123456789", "123456788", 1)

// 	fmt.Println("\nVerify tampered ID:")
// 	id2, ok2 := VerifyID(tampered, secretKey)
// 	fmt.Println("Valid :", ok2)
// 	fmt.Println("ID    :", id2)
// }
