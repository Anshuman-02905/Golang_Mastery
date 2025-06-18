package lru

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {

	test_cases := []struct {
		name               string
		maxEntries         int
		expectedMaxEntries int
	}{{"base cases", 5, 5},
		{"zero entries", 0, 0},
	}

	for _, tc := range test_cases {

		t.Run(tc.name, func(t *testing.T) {
			lru := New(tc.maxEntries)
			if lru.MaxEntries != tc.expectedMaxEntries {
				t.Fatalf("expected MaxEntries = %d, got = %d", tc.expectedMaxEntries, lru.MaxEntries)
			}
			if lru.ll == nil {
				t.Fatal("Expected ll (list) to be initialized")
			}
			if lru.cache == nil {
				t.Fatal("expected cache map to be initialized")
			}

		})

	}
}

func TestAdd(t *testing.T) {
	test_cases := []struct {
		name       string
		MaxEntries int
		keys       []Key
		values     []interface{}
	}{
		{
			name:       "simple insert",
			MaxEntries: 2,
			keys:       []Key{"Node1"},
			values:     []interface{}{"value1"}},

		{
			name:       "multiple insert ",
			MaxEntries: 2,
			keys:       []Key{"Node1", "Node2"},
			values:     []interface{}{"value1", "value2"}},
	}
	for _, tc := range test_cases {
		t.Run(tc.name, func(t *testing.T) {
			ch := New(tc.MaxEntries)

			for i, k := range tc.keys {
				ch.Add(k, tc.values[i])
			}
			//
			start := len(tc.keys) - min(len(tc.keys), tc.MaxEntries)
			for i := start; i < len(tc.keys); i++ {
				k := tc.keys[i]
				if ele, ok := ch.cache[k]; !ok {
					t.Errorf("expected key %v to be present", k)
				} else {
					v := ele.Value.(*entry).value
					if v != tc.values[i] {
						t.Errorf("expected value %v to be present", v)
					}
				}
			}

			for i := 0; i < start; i++ {
				t.Logf("these keys should be evicted %d", i)

			}

		})
	}
}

func TestGet(t *testing.T) {
	test_cases := []struct {
		name       string
		maxEntries int
		key        []Key
		value      []interface{}
	}{
		{
			name:       "one Insert",
			maxEntries: 1,
			key:        []Key{"Node1"},
			value:      []interface{}{"Value1"},
		},
		{
			name:       "Multiple Insert",
			maxEntries: 3,
			key:        []Key{"Node1", "Node2", "Node3"},
			value:      []interface{}{"Value1", "Value2", "Value3"},
		},
		{
			name:       "No Insert",
			maxEntries: 1,
			key:        []Key{},
			value:      []interface{}{},
		},
	}

	for _, tc := range test_cases {
		t.Log(tc.name)
		ch := New(tc.maxEntries)
		for i, k := range tc.key {
			ch.Add(k, tc.value[i])
		}
		for i, k := range tc.key {
			get_value, ok := ch.Get(k)
			if !ok {
				t.Fatal("Expected data to be be present in the cache")
			}
			if get_value != tc.value[i] {
				t.Fatalf("Expected value %v in the cache but got %v ", get_value, tc.value[i])
			}
		}

	}
}

func TestRemove(t *testing.T) {
	test_cases := []struct {
		name                 string
		MaxEntries           int
		Keys                 []Key
		Keys_to_Remove       []Key
		Keys_to_be_Present   []Key
		Values               []interface{}
		Values_to_be_Present []interface{}
	}{
		{
			name:                 "single Entry Single remove",
			MaxEntries:           1,
			Keys:                 []Key{"Node1"},
			Keys_to_Remove:       []Key{"Node1"},
			Values:               []interface{}{"val1"},
			Keys_to_be_Present:   []Key{},
			Values_to_be_Present: []interface{}{},
		},
		{
			name:                 "Double Entry Single remove",
			MaxEntries:           2,
			Keys:                 []Key{"Node1", "Node2"},
			Keys_to_Remove:       []Key{"Node1"},
			Values:               []interface{}{"val1", 5},
			Keys_to_be_Present:   []Key{"Node2"},
			Values_to_be_Present: []interface{}{5},
		},
		{
			name:                 "Double Entry Double remove",
			MaxEntries:           2,
			Keys:                 []Key{"Node1", "Node2"},
			Keys_to_Remove:       []Key{"Node1", "Node2"},
			Values:               []interface{}{"val1", 5},
			Keys_to_be_Present:   []Key{},
			Values_to_be_Present: []interface{}{},
		},
	}

	for _, tc := range test_cases {
		t.Log(tc.name)
		ch := New(tc.MaxEntries)
		for i, key := range tc.Keys {
			ch.Add(key, tc.Values[i])
		}

		for _, key := range tc.Keys_to_Remove {
			ch.Remove(key)
		}
		for j, key := range tc.Keys_to_be_Present {
			if value, ok := ch.Get(key); !ok {
				if tc.Values_to_be_Present[j] != nil {
					t.Fatalf("Value Expeceted nil but got %v", tc.Values_to_be_Present[j])

				}

			} else {
				if tc.Values_to_be_Present[j] != value {
					t.Fatalf("Value is matched between fetched %v and test %v ", value, tc.Values_to_be_Present[j])
					return
				}
			}
		}

	}
}

func TestEvict(t *testing.T){

	evictedKeys :=make([]Key,0)
	onEvictedFun :=func(key Key, value interface{}){
		evictedKeys = append(evictedKeys, key)
	}

	lru:=New(20)
	lru.OnEvicted=onEvictedFun
	for i:=0;i<22;i++{
		lru.Add(fmt.Sprintf("myKey%d",i),1234)
	}
	if len(evictedKeys)!=2{
		t.Fatalf("got %d evicted Keys :want 2",len(evictedKeys))
	}
	if evictedKeys[0]!=Key("myKey0"){
		t.Fatalf("got %v in the first evicted key: want %s",evictedKeys[0],"myKey0")
	}
		if evictedKeys[1]!=Key("myKey1"){
		t.Fatalf("got %v in the first evicted key: want %s",evictedKeys[1],"myKey1")
	}

}