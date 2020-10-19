package manager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/viert/kstore/sec"
	"github.com/viert/kstore/store"
)

var (
	clientID     = ""
	clientSecret = ""
)

type authData struct {
	AccessToken string `json:"access_token"`
}

// Authenticate makes all the necessary steps to authenticate user in Yandex.Disk
func (m *Manager) Authenticate() (err error) {
	// get master password from user
	err = m.acquireMasterKey()
	if err != nil {
		return
	}

	// load auth file
	authdata, err := loadAuthData()

	if err != nil {
		if os.IsNotExist(err) {
			// no file found, getting token from web (and saving to auth file)
			err = m.acquireTokenFromWeb()
			if err != nil {
				return err
			}
		} else {
			// something different has happend to file, returning the error
			return
		}
	} else {
		// getting token from file
		err = m.acquireTokenFromAuthData(authdata)
		if err != nil {
			return err
		}
	}

	fmt.Println("Checking credentials...")
	m.store, err = store.New(m.accessToken)
	if err != nil {
		return
	}

	return
}

func (m *Manager) acquireTokenFromAuthData(data []byte) error {
	dec, err := m.enc.Decrypt(data)
	if err != nil {
		return err
	}

	var ad authData
	err = json.Unmarshal(dec, &ad)
	if err != nil {
		return err
	}
	m.accessToken = ad.AccessToken
	return nil
}

func (m *Manager) acquireTokenFromWeb() error {
	uri := fmt.Sprintf("https://oauth.yandex.ru/authorize?response_type=code&client_id=%s", clientID)
	fmt.Printf("Open the following link: %s\n", uri)
	code, err := getString("Enter confirmation code (will be hidden): ")
	if err != nil {
		return err
	}

	values := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	}

	resp, err := http.PostForm("https://oauth.yandex.ru/token", values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var ad authData
	err = json.Unmarshal(data, &ad)
	if err != nil {
		return err
	}
	m.accessToken = ad.AccessToken

	authDataContent, err := json.Marshal(ad)
	if err != nil {
		return err
	}

	authFile, err := getAuthFilename()
	if err != nil {
		return err
	}

	encrypted, err := m.enc.Encrypt(authDataContent)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(authFile, encrypted, 0600)
	return err
}

func (m *Manager) acquireMasterKey() error {
	pass, err := m.rl.ReadPassword("Enter Master Password: ")
	if err != nil {
		return err
	}
	key, err := sec.CreatePassKey(pass)
	if err != nil {
		return err
	}
	m.enc, err = sec.NewAES(key)
	return err
}

func getAuthFilename() (string, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	authFile := path.Join(userHome, ".kstore.auth")
	return authFile, nil
}

func loadAuthData() ([]byte, error) {
	authFile, err := getAuthFilename()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(authFile)
}
