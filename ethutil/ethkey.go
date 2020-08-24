package ethutil

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"fmt"
	"io/ioutil"
)

func DecryptKeyFile(keyFile, password string) (*keystore.Key, error) {
	data, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read KeyFile %s [%v]", keyFile, err)
	}
	key, err := keystore.DecryptKey(data, password)
	if err != nil {
		return nil, fmt.Errorf("unable to decrypt %s [%v]", keyFile, err)
	}
	return key, nil
}