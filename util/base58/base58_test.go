package base58

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestDecodeCheck(t *testing.T) {
	decodeBytes := DecodeCheck("27ZESitosJfKouTBrGg6Nk5yEjnJHXMbkZp")

	decode := hex.EncodeToString(decodeBytes)

	if strings.EqualFold(decode, "a06f61d05912402335744c288d4b72a735ede35604") {
		t.Log("success")
	} else {
		t.Fatalf("failure: %s", decode)
	}
}
