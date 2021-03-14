package pcache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expected := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expected) {
		t.Errorf("callback failed")
	}
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	pcache := NewGroup("scores", 2<<10, GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			if _, ok := loadCounts[key]; !ok {
				loadCounts[key] = 0
			}
			loadCounts[key] += 1
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))

	for k, v := range db {
		if view, err := pcache.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value of Tom")
		}
		if _, err := pcache.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}

	if view, err := pcache.Get("unknow"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}
