package stdserver

import (
	"crypto/ecdsa"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-jwt/jwt/v4"
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
	sign    map[string]interface{}
	verify  map[string]interface{}
}

func NewKeyTable() *KeyTable {
	return &KeyTable{
		entries: make([]KeyEntry, 0),
		sign:    make(map[string]interface{}),
		verify:  make(map[string]interface{}),
	}
}

func (k *KeyTable) PutECKey(id string, key *ecdsa.PrivateKey) {
	k.entries = append(k.entries, KeyEntry{
		Id:  id,
		Key: key,
	})
	k.sign[id] = key
	k.verify[id] = &key.PublicKey
}

func (k *KeyTable) GetPrivateKeys() map[string]interface{} {
	return k.sign
}

func (k *KeyTable) GetPublicKeys() map[string]interface{} {
	return k.verify
}

func LoadKeyTableFromDir(root string) (*KeyTable, error) {
	k := NewKeyTable()
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ext {
			dat, err := ioutil.ReadFile(path)
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
