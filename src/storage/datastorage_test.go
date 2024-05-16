package storage

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

type fakeStorage struct {
	data map[string]interface{}
}

type testData struct {
	s string
}

var tData = testData{"test"}

func NewFakeStorage() *fakeStorage {
	return &fakeStorage{make(map[string]interface{})}
}

func (s *fakeStorage) Get(uid string) (payload interface{}, createdAt time.Time, ttl time.Duration, err error) {
	payload, ok := s.data[uid]
	if !ok {
		err = errors.New("uid not found")
	}
	return
}

func (s *fakeStorage) Put(uid string, payload interface{}, ttl *time.Duration) error {
	_, _, _, err := s.Get(uid)
	if err == nil {
		return errors.New("data with uid already exists")
	}

	s.data[uid] = payload
	return nil
}

func CheckMemoryStorageTest(storage *DataStorage, storageM *fakeStorage, storageP *fakeStorage) error {
	const UID = "CHECK_MEMORY_TEST"

	if err := storage.Put(UID, &tData, nil); err != nil {
		return fmt.Errorf("error on Put: %v", err)
	}

	dMi, _, _, err := storageM.Get(UID)
	if err != nil {
		return fmt.Errorf("error on Get from memory storage: %v", err)
	}

	_, _, _, err = storageP.Get(UID)
	if err == nil {
		return fmt.Errorf("data has been found in the persistent storage: %v", err)
	}

	return CheckReturnedData(dMi)
}

func CheckPersistentStorageTest(storage *DataStorage, storageM *fakeStorage, storageP *fakeStorage) error {
	const UID = "CHECK_PERSISTENT_TEST"

	if err := storage.Put(UID, &tData, &ttlForever); err != nil {
		return fmt.Errorf("error on Put: %v", err)
	}

	dMi, _, _, err := storageM.Get(UID)
	if err != nil {
		return fmt.Errorf("error on getting from the memory storage: %v", err)
	}

	dPi, _, _, err := storageP.Get(UID)
	if err != nil {
		return fmt.Errorf("error on getting from the persistent storage: %v", err)
	}

	if err = CheckReturnedData(dMi); err != nil {
		return fmt.Errorf("memory storage error: %v", err)
	}

	if err = CheckReturnedData(dPi); err != nil {
		return fmt.Errorf("persistent storage error: %v", err)
	}
	return nil
}

func CheckUniq(storage *DataStorage) error {
	const UID = "CHECK_UNIQ_TEST"

	if err := storage.Put(UID, &tData, nil); err != nil {
		return fmt.Errorf("error on Put: %v", err)
	}

	if err := storage.Put(UID, &tData, nil); err == nil {
		return errors.New("put should return an error on putting the same UID twice")
	}
	return nil
}

func CheckReturnedData(payload interface{}) error {
	d, ok := payload.(*testData)
	if !ok {
		return fmt.Errorf("wrong data type has been returned")
	}

	if d.s != tData.s {
		return fmt.Errorf("wrong data has been returned")
	}
	return nil
}

func TestDataStorage(t *testing.T) {

	storageM, storageP := NewFakeStorage(), NewFakeStorage()
	ds := NewDataStorage(storageM, storageP)

	if err := CheckUniq(ds); err != nil {
		t.Errorf("Uniq test: %v", err)
	}

	if err := CheckMemoryStorageTest(ds, storageM, storageP); err != nil {
		t.Errorf("Memory storage test: %v", err)
	}

	if err := CheckPersistentStorageTest(ds, storageM, storageP); err != nil {
		t.Errorf("Persistent storage test: %v", err)
	}
}
