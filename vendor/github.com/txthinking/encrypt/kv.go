package encrypt

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"
)

// KV can be used to crypt value by key and give it a lifecycle (s).
// Can be used in cookies
type KV struct {
	AESKey []byte
}

// Encrypt key, value
func (kv *KV) Encrypt(k string, v string) (string, error) {
	m := map[string]interface{}{
		"k": k,
		"v": v,
		"t": time.Now().Unix(),
	}
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	b, err = AESCFBEncrypt(b, kv.AESKey)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Decrypt key, value, lifecycle is second, 0 is no expired time.
func (kv *KV) Decrypt(c string, k string, lifecycle int64) (string, error) {
	b, err := hex.DecodeString(c)
	if err != nil {
		return "", err
	}
	m := make(map[string]interface{})
	d, err := AESCFBDecrypt(b, kv.AESKey)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(d, &m); err != nil {
		return "", err
	}
	if lifecycle != 0 {
		if int64(m["t"].(float64))+lifecycle < time.Now().Unix() {
			return "", errors.New("Expired")
		}
	}
	if m["k"].(string) != k {
		return "", errors.New("Unmarch key")
	}
	return m["v"].(string), nil
}

// JSON it first
func (kv *KV) EncryptStruct(k string, v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	m := map[string]interface{}{
		"k": k,
		"v": hex.EncodeToString(b),
		"t": time.Now().Unix(),
	}
	b, err = json.Marshal(m)
	if err != nil {
		return "", err
	}
	b, err = AESCFBEncrypt(b, kv.AESKey)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Decrypt key, value, lifecycle is second, 0 is no expired time.
func (kv *KV) DecryptStruct(c string, v interface{}, k string, lifecycle int64) error {
	b, err := hex.DecodeString(c)
	if err != nil {
		return err
	}
	m := make(map[string]interface{})
	d, err := AESCFBDecrypt(b, kv.AESKey)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(d, &m); err != nil {
		return err
	}
	if lifecycle != 0 {
		if int64(m["t"].(float64))+lifecycle < time.Now().Unix() {
			return errors.New("Expired")
		}
	}
	if m["k"].(string) != k {
		return errors.New("Unmarch key")
	}
	b, err = hex.DecodeString(m["v"].(string))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, v); err != nil {
		return err
	}
	return nil
}
