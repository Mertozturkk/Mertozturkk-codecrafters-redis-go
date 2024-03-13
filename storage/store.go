package storage

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Store struct {
	KeyValue map[string]string
	Expiry   map[string]time.Time
	Lock     sync.Mutex
}

func NewStore(initialKey, initialValue string) *Store {
	return &Store{
		KeyValue: map[string]string{initialKey: initialValue},
		Expiry:   make(map[string]time.Time),
	}
}

func (s *Store) Set(key string, value string, duration time.Duration) error {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	if len(key) == 0 {
		return errors.New("key cannot be empty")
	}

	expiry := time.Now().Add(duration)
	s.KeyValue[key] = value
	s.Expiry[key] = expiry
	fmt.Println("Key:", key, "Value:", value, "Expiry:", expiry)

	return nil
}

func (s *Store) Get(key string) (string, bool) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	expiry, ok := s.Expiry[key]
	if !ok || expiry.Before(time.Now()) {
		delete(s.KeyValue, key)
		delete(s.Expiry, key)
		return "", false
	}

	val := s.KeyValue[key]
	return val, true
}
