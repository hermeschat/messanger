package auth

import (
	"encoding/base64"
	"strings"
)

// HasAccountAccess check if current account has access to given account_id
func HasAccountAccess(accountID string, roles ...string) bool {
	return false
}

// AccountHasRole Check if current account has role
func AccountHasRole(roles ...string) bool {
	return false
}

// base64Decode Decode specific base64url encoding with padding stripped
func base64Decode(seg string) ([]byte, error) {
	seg = base64ScapePadding(seg)
	return base64.URLEncoding.DecodeString(seg)
}

func base64ScapePadding(seg string) string {
	seg = strings.TrimRight(seg, "=")
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}
	return seg
}

func base64ScapeCharacter(seg string) string {
	for i := 0; i < len(seg); i++ {
		if seg[i] == '_' {
			seg = seg[:i] + "/" + seg[i+1:]
		} else if seg[i] == '-' {
			seg = seg[:i] + "+" + seg[i+1:]
		}
	}
	return seg
}
