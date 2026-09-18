package hash

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

func MD5(v any) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:]), nil
}
