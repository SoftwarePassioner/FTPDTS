package storage

import (
	"errors"
	"github.com/starshiptroopers/uidgenerator"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

//comprehensive FsDataStorage test
func TestFsDataStorage(t *testing.T) {

	type testData struct {
		S string
	}
	var tData = testData{"test"}

	dir, err := ioutil.TempDir("", "testing")
	if err != nil {
		t.Fatalf("can't create temporay directory for testing, %v", err)
		return
	}
	defer func() { _ = os.Remove(dir) }()

	compare := func(data testData, loadedData interface{}) error {
		err := errors.New("wrong data type has been returned")
		d, ok := loadedData.(map[string]interface{})
		if !ok {
			return err
		}
		text, ok := d["S"].(string)
		if !ok {
			return err
		}
		if text != data.S {
			return errors.New("wrong data has been returned")
		}
		return nil
	}

	UIDGenerator := uidgenerator.New(
		&uidgenerator.Cfg{
			Alfa:      "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
			Format:    "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			Validator: "[0-9a-zA-Z]{32}",
		},
	)

	ds := NewFsDataStorage(dir, UIDGenerator)
	uid1 := UIDGenerator.New()
	uid2 := UIDGenerator.New()

	err = ds.Put(uid1, &tData, nil)
	if err != nil {
		t.Fatalf("can't put data into the storage: %v", err)
	}

	di, _, _, err := ds.Get(uid1)

	if err := compare(tData, di); err != nil {
		t.Errorf("Wrong data has been read from the storage: %v", err)
		return
	}

	err = ds.Put(uid2, tData, nil)
	if err != nil {
		t.Fatalf("can't put data into the storage: %v", err)
	}

	//check data loading
	var uid1ok, uid2ok bool
	err = ds.Pass(func(uid string, createdAt time.Time, data interface{}) {
		if uid != uid1 {
			uid1ok = true
		} else if uid != uid2 {
			uid2ok = true
		} else {
			t.Errorf("Unexpected uid has been received")
			return
		}
		if err := compare(tData, data); err != nil {
			t.Errorf("Wrong data has been read from the storage: %v", err)
			return
		}
	})
	if err != nil {
		t.Errorf("can't read data: %v", err)
	}

	if !uid1ok || !uid2ok {
		t.Errorf("not all data received")
	}
}
