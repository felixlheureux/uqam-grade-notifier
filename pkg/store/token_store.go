package store

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

type TokenStore struct {
	path string
	mu   sync.RWMutex
}

type tokenData struct {
	Tokens map[string]string `json:"tokens"` // email -> token
}

func NewTokenStore(path string) (*TokenStore, error) {
	store := &TokenStore{path: path}
	if err := store.load(); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *TokenStore) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := os.Stat(s.path); err != nil {
		if os.IsNotExist(err) {
			return s.save(&tokenData{Tokens: make(map[string]string)})
		}
		return err
	}

	data, err := ioutil.ReadFile(s.path)
	if err != nil {
		return err
	}

	var store tokenData
	if err := json.Unmarshal(data, &store); err != nil {
		return err
	}

	if store.Tokens == nil {
		store.Tokens = make(map[string]string)
	}

	return s.save(&store)
}

func (s *TokenStore) save(data *tokenData) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(s.path, bytes, 0644)
}

func (s *TokenStore) SaveToken(email, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := ioutil.ReadFile(s.path)
	if err != nil {
		return err
	}

	var store tokenData
	if err := json.Unmarshal(data, &store); err != nil {
		return err
	}

	store.Tokens[email] = token
	return s.save(&store)
}

func (s *TokenStore) GetToken(email string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := ioutil.ReadFile(s.path)
	if err != nil {
		return "", err
	}

	var store tokenData
	if err := json.Unmarshal(data, &store); err != nil {
		return "", err
	}

	token, exists := store.Tokens[email]
	if !exists {
		return "", nil
	}

	return token, nil
}

func (s *TokenStore) DeleteToken(email string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := ioutil.ReadFile(s.path)
	if err != nil {
		return err
	}

	var store tokenData
	if err := json.Unmarshal(data, &store); err != nil {
		return err
	}

	delete(store.Tokens, email)
	return s.save(&store)
}
