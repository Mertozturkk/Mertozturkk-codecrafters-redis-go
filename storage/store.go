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

var KeyValues = make(map[string]string)

func NewStore(initialKey, initialValue string) *Store {
	return &Store{
		KeyValue: map[string]string{initialKey: initialValue},
		Expiry:   make(map[string]time.Time),
	}
}

func (s *Store) Set(key string, value string, duration time.Duration) error {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	duration = 1000 * time.Second
	if len(key) == 0 {
		return errors.New("key cannot be empty")
	}

	expiry := time.Now().Add(duration)
	s.KeyValue[key] = value
	s.Expiry[key] = expiry
	fmt.Printf("Key: %v, Value: %v, Expiry: %v\n", key, value, expiry)

	return nil
}
func (s *Store) Get(key string) (string, bool) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	val, ok := s.KeyValue[key]
	return val, ok
}
