package stdserver

import (
	"crypto/ecdsa"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ext = ".pem"
)

type KeyEntry struct {
	Id  string
	Key interface{}
}

type KeyTable struct {
	entries []KeyEntry
	sign    map[string]any
	verify  map[string]jwtware.SigningKey
}

func NewKeyTable() *KeyTable {
	return &KeyTable{
		entries: make([]KeyEntry, 0),
		sign:    make(map[string]any),
		verify:  make(map[string]jwtware.SigningKey),
	}
}

func (k *KeyTable) PutECKey(id string, key *ecdsa.PrivateKey) {
	k.entries = append(k.entries, KeyEntry{
		Id:  id,
		Key: key,
	})
	k.sign[id] = key
	alg := "ES" + strconv.Itoa(key.Curve.Params().N.BitLen())
	k.verify[id] = jwtware.SigningKey{
		JWTAlg: alg,
		Key:    &key.PublicKey,
	}
}

func (k *KeyTable) GetPrivateKeys() map[string]any {
	return k.sign
}

func (k *KeyTable) GetPublicKeys() map[string]jwtware.SigningKey {
	return k.verify
}

func LoadKeyTableFromDir(root string) (*KeyTable, error) {
	k := NewKeyTable()
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ext {
			dat, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			key, err := jwt.ParseECPrivateKeyFromPEM(dat)
			if err != nil {
				return err
			}
			name := info.Name()
			id := name[:strings.LastIndex(name, ext)]
			k.PutECKey(id, key)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return k, nil
}
