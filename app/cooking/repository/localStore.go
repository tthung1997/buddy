package repository

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

var errEmptyKey = errors.New("key must not be empty")

// localStore persists a keyed collection as a single JSON object on disk. It
// follows the same self seeding approach as app/random/choice so a fresh clone,
// where the gitignored .db directory does not exist yet, still works.
type localStore[T any] struct {
	filePath string
	mutex    sync.RWMutex
}

func newLocalStore[T any](filePath string) (*localStore[T], error) {
	if filePath == "" {
		return nil, errors.New("file path must not be empty")
	}
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return nil, err
	}

	store := &localStore[T]{filePath: filePath}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := store.save(map[string]T{}); err != nil {
			return nil, err
		}
	}

	return store, nil
}

func (s *localStore[T]) load() (map[string]T, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]T{}, nil
		}
		return nil, err
	}

	items := map[string]T{}
	if len(data) == 0 {
		return items, nil
	}
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (s *localStore[T]) save(items map[string]T) error {
	data, err := json.MarshalIndent(items, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0644)
}

func (s *localStore[T]) all() (map[string]T, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.load()
}

func (s *localStore[T]) get(key string) (T, bool, error) {
	var zero T
	if key == "" {
		return zero, false, errEmptyKey
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	items, err := s.load()
	if err != nil {
		return zero, false, err
	}
	item, exists := items[key]
	return item, exists, nil
}

func (s *localStore[T]) put(key string, value T) error {
	if key == "" {
		return errEmptyKey
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	items, err := s.load()
	if err != nil {
		return err
	}
	items[key] = value
	return s.save(items)
}

func (s *localStore[T]) remove(key string) error {
	if key == "" {
		return errEmptyKey
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	items, err := s.load()
	if err != nil {
		return err
	}
	delete(items, key)
	return s.save(items)
}
