package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Store represents an API interface to yandex disk
type Store struct {
	accessToken string
	cli         *http.Client
}

type diskResponse struct {
	Href      string `json:"href"`
	Method    string `json:"method"`
	Templated bool   `json:"templated"`
}

const (
	dbFilename = "/db.bin"
	maxBackups = 5
)

// New creates a new authdata store
func New(accessToken string) (*Store, error) {
	s := &Store{
		accessToken: accessToken,
		cli:         &http.Client{},
	}

	err := s.check()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) newRequest(method string, uri string, body io.Reader) (*http.Request, error) {
	uri = fmt.Sprintf("https://cloud-api.yandex.net%s", uri)
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("OAuth %s", s.accessToken))
	req.Header.Add("Accept", "application/json")
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	return req, nil
}

func (s *Store) check() error {
	req, err := s.newRequest("GET", "/v1/disk/resources/?path=app:%2F", nil)
	if err != nil {
		return err
	}
	resp, err := s.cli.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return responseError(resp)
	}
	return nil
}

func appFilename(path string) string {
	return fmt.Sprintf(
		"app:%s",
		url.QueryEscape(path),
	)
}

func (s *Store) fileExists(path string) (bool, error) {
	fullPath := appFilename(path)

	req, err := s.newRequest("GET",
		fmt.Sprintf("/v1/disk/resources?path=%s", fullPath), nil)

	if err != nil {
		return false, err
	}

	resp, err := s.cli.Do(req)
	if err != nil {
		return false, err
	}

	if resp.StatusCode == 404 {
		return false, nil
	} else if resp.StatusCode == 200 {
		return true, nil
	}

	return false, responseError(resp)
}

func (s *Store) moveFile(src string, dst string) error {
	found, err := s.fileExists(src)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}

	src = appFilename(src)
	dst = appFilename(dst)

	path := fmt.Sprintf(
		"/v1/disk/resources/move?from=%s&path=%s&overwrite=true",
		src,
		dst,
	)

	req, err := s.newRequest("POST", path, nil)
	if err != nil {
		return err
	}

	resp, err := s.cli.Do(req)
	if err != nil {
		return err
	}
	if !isSuccess(resp) {
		return responseError(resp)
	}
	return nil
}

func (s *Store) backup() error {
	var src string
	var dst string
	var err error

	fmt.Print("Creating Backups")

	for i := maxBackups - 1; i > 0; i-- {
		src = fmt.Sprintf("%s.%d", dbFilename, i)
		dst = fmt.Sprintf("%s.%d", dbFilename, i+1)

		err := s.moveFile(src, dst)
		if err != nil {
			fmt.Print("!")
		} else {
			fmt.Print(".")
		}
	}
	dst = src
	src = dbFilename
	err = s.moveFile(src, dst)
	if err != nil {
		fmt.Println("!")
	} else {
		fmt.Println(".")
	}

	return nil
}

// Load loads the actual version of db.bin and returns its contents
func (s *Store) Load() ([]byte, error) {
	filename := appFilename(dbFilename)
	path := fmt.Sprintf("/v1/disk/resources/download?path=%s", filename)
	req, err := s.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !isSuccess(resp) {
		return nil, responseError(resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var dr diskResponse
	err = json.Unmarshal(body, &dr)
	if err != nil {
		return nil, err
	}

	req, err = http.NewRequest(dr.Method, dr.Href, nil)
	if err != nil {
		return nil, err
	}

	resp, err = s.cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !isSuccess(resp) {
		return nil, responseError(resp)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}

// Save saves data into a db.bin file backuping prev versions beforehand
func (s *Store) Save(data []byte) error {
	s.backup()

	fmt.Println("Saving data...")
	filename := appFilename(dbFilename)
	path := fmt.Sprintf("/v1/disk/resources/upload?path=%s&overwrite=true", filename)
	req, err := s.newRequest("GET", path, nil)
	if err != nil {
		return err
	}

	resp, err := s.cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if !isSuccess(resp) {
		return responseError(resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var ur diskResponse
	err = json.Unmarshal(body, &ur)
	if err != nil {
		return err
	}

	req, err = http.NewRequest(ur.Method, ur.Href, bytes.NewReader(data))
	if err != nil {
		return err
	}

	resp, err = s.cli.Do(req)
	if err != nil {
		return err
	}

	if !isSuccess(resp) {
		return responseError(resp)
	}
	return nil

}
