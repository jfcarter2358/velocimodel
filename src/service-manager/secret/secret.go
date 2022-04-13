package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"service-manager/config"
	"service-manager/utils"

	"github.com/jfcarter2358/ceresdb-go/connection"
)

func LoadSecrets() {
	jsonFile, err := os.Open(config.Config.SecretsPath)
	if err != nil {
		panic(err)
	}

	var secrets map[string]interface{}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &secrets)
	if err != nil {
		panic(err)
	}
	err = UpdateSecrets(secrets)
	if err != nil {
		panic(err)
	}
}

func DeleteSecrets(paramNames []string) error {
	queryString := fmt.Sprintf("get record %v.secrets", config.Config.DBName)
	currentData, err := connection.Query(queryString)
	if err != nil {
		return err
	}
	ids := make([]string, 0)
	for _, datum := range currentData {
		if utils.Contains(paramNames, datum["name"].(string)) {
			ids = append(ids, datum[".id"].(string))
		}
	}
	queryData, _ := json.Marshal(&ids)
	queryString = fmt.Sprintf("delete record %v.secrets %v", config.Config.DBName, queryData)
	_, err = connection.Query(queryString)
	return err
}

func GetSecrets() (map[string]interface{}, error) {
	queryString := fmt.Sprintf("get record %v.secrets", config.Config.DBName)
	data, err := connection.Query(queryString)
	if err != nil {
		return nil, err
	}
	output := make(map[string]interface{})
	for _, datum := range data {
		val, err := decrypt(datum["value"].(string), config.Config.EncryptionKey)
		if err != nil {
			return nil, err
		}
		output[datum["name"].(string)] = val
	}
	return output, nil
}

func SetSecrets(input map[string]interface{}) error {
	queryList := make([]map[string]interface{}, 0)
	for key, val := range input {
		encryptedVal, err := encrypt(val.(string), config.Config.EncryptionKey)
		if err != nil {
			return err
		}
		tmpObject := map[string]interface{}{"name": key, "value": encryptedVal}
		queryList = append(queryList, tmpObject)
	}

	queryData, _ := json.Marshal(&queryList)
	queryString := fmt.Sprintf("post record %v.secrets %v", config.Config.DBName, string(queryData))
	_, err := connection.Query(queryString)
	return err
}

func UpdateSecrets(input map[string]interface{}) error {
	queryString := fmt.Sprintf("get record %v.secrets", config.Config.DBName)
	currentData, err := connection.Query(queryString)
	if err != nil {
		return err
	}
	newData := make([]map[string]interface{}, 0)
	for key, val := range input {
		exists := false
		for idx, data := range currentData {
			if data["name"] == key {
				encryptedVal, err := encrypt(val.(string), config.Config.EncryptionKey)
				if err != nil {
					return err
				}
				data["value"] = encryptedVal
				currentData[idx] = data
				exists = true
				break
			}
		}
		if !exists {
			encryptedVal, err := encrypt(val.(string), config.Config.EncryptionKey)
			if err != nil {
				return err
			}
			tmpObject := map[string]interface{}{"name": key, "value": encryptedVal}
			newData = append(newData, tmpObject)
		}
	}

	queryData, _ := json.Marshal(&currentData)
	queryString = fmt.Sprintf("put record %v.secrets %v", config.Config.DBName, string(queryData))
	_, err = connection.Query(queryString)
	if err != nil {
		return err
	}
	if len(newData) > 0 {
		queryData, _ := json.Marshal(&newData)
		queryString = fmt.Sprintf("post record %v.secrets %v", config.Config.DBName, string(queryData))
		_, err = connection.Query(queryString)
		if err != nil {
			return err
		}
	}
	return nil
}

// Encrypt encrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Output takes the
// form nonce|ciphertext|tag where '|' indicates concatenation.
func encrypt(textToEncrypt string, key [32]byte) (string, error) {
	plaintext := []byte(textToEncrypt)
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext), nil
}

// Decrypt decrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Expects input
// form nonce|ciphertext|tag where '|' indicates concatenation.
func decrypt(textToDecrypt string, key [32]byte) (string, error) {
	ciphertext, _ := hex.DecodeString(textToDecrypt)
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.New("malformed ciphertext")
	}

	plaintext, err := gcm.Open(nil, ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():], nil)
	if err != nil {
		return "", nil
	}
	return string(plaintext), nil
}
