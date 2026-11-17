package safeorderedmap_test

import (
	"encoding/json"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/thalesfsp/go-common-types/safeorderedmap"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func finish(t *testing.T, call func()) {
	t.Helper()

	done := make(chan struct{})
	go func() {
		defer close(done)

		call()
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("collection operation deadlocked")
	}
}

func TestCallbackReentrySnapshot(t *testing.T) {
	calls := map[string]func(*safeorderedmap.SafeOrderedMap[int], func(string, int)){
		"all": func(m *safeorderedmap.SafeOrderedMap[int], visit func(string, int)) {
			m.All(func(k string, v int) bool { visit(k, v); return true })
		},
		"map": func(m *safeorderedmap.SafeOrderedMap[int], visit func(string, int)) {
			m.Map(func(k string, v int) int { visit(k, v); return v })
		},
		"filter": func(m *safeorderedmap.SafeOrderedMap[int], visit func(string, int)) {
			m.Filter(func(k string, v int) bool { visit(k, v); return true })
		},
		"each": func(m *safeorderedmap.SafeOrderedMap[int], visit func(string, int)) { m.Each(visit) },
		"reduce": func(m *safeorderedmap.SafeOrderedMap[int], visit func(string, int)) {
			m.Reduce(func(a int, k string, v int) int { visit(k, v); return a + v }, 0)
		},
		"find": func(m *safeorderedmap.SafeOrderedMap[int], visit func(string, int)) {
			m.Find(func(k string, v int) bool { visit(k, v); return true })
		},
		"any": func(m *safeorderedmap.SafeOrderedMap[int], visit func(string, int)) {
			m.Any(func(k string, v int) bool { visit(k, v); return true })
		},
		"take": func(m *safeorderedmap.SafeOrderedMap[int], visit func(string, int)) {
			m.TakeWhile(func(k string, v int) bool { visit(k, v); return true })
		},
		"drop": func(m *safeorderedmap.SafeOrderedMap[int], visit func(string, int)) {
			m.DropWhile(func(k string, v int) bool { visit(k, v); return false })
		},
	}
	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			m := safeorderedmap.New[int]().Add("original", 1)
			visits := 0

			finish(t, func() {
				call(m, func(string, int) { visits++; m.Add("added", 2) })
			})

			if visits != 1 || m.Size() != 2 {
				t.Fatalf("expected one snapshot visit and a reentrant insertion; visits=%d size=%d", visits, m.Size())
			}
		})
	}
}

func TestBinaryOperationsConcurrent(t *testing.T) {
	a := safeorderedmap.New[int]().Add("stable", 1)
	b := safeorderedmap.New[int]().Add("stable", 1)

	finish(t, func() {
		var workers sync.WaitGroup

		workers.Add(3)
		go func() {
			defer workers.Done()

			for range 2000 {
				b.Add("changing", 2).Delete("changing")
			}
		}()

		for range 2 {
			go func() {
				defer workers.Done()

				for range 2000 {
					a.Union(b)
					a.Difference(b)
					a.Intersection(b)
					a.Subset(b)
					a.Superset(b)
					b.Subset(a)
				}
			}()
		}

		workers.Wait()
	})

	if !a.Union(b).Contains("stable") || !a.Intersection(b).Contains("stable") {
		t.Fatal("binary operations lost the stable entry")
	}

	if !a.Subset(a) || !a.Superset(a) || !a.Difference(a).Empty() {
		t.Fatal("self operations violated set identities")
	}
}

func TestBSONWireRoundTrip(t *testing.T) {
	for _, document := range []bson.D{{}, {{Key: "key", Value: int32(1)}}} {
		wire, err := bson.Marshal(document)
		if err != nil {
			t.Fatal(err)
		}

		err = bson.Raw(wire).Validate()
		if err != nil {
			t.Fatal(err)
		}

		m := safeorderedmap.New[int32]()

		finish(t, func() {
			err = bson.Unmarshal(wire, m)
		})

		if err != nil || m.Size() != len(document) {
			t.Fatalf("BSON decode: size=%d err=%v", m.Size(), err)
		}

		encoded, err := bson.Marshal(m)
		if err != nil {
			t.Fatal(err)
		}

		err = bson.Raw(encoded).Validate()
		if err != nil {
			t.Fatalf("MarshalBSON returned invalid wire bytes: %v", err)
		}

		var decoded bson.D

		err = bson.Unmarshal(encoded, &decoded)
		if err != nil || !reflect.DeepEqual(document, decoded) {
			t.Fatalf("typed BSON contents changed: %v -> %v, err=%v", document, decoded, err)
		}

		before := m.Clone()

		err = m.UnmarshalBSON(wire[:len(wire)-1])
		if err == nil || !reflect.DeepEqual(before.Values(), m.Values()) {
			t.Fatal("truncated BSON must fail without changing the map")
		}

		err = m.UnmarshalBSON(append(append([]byte(nil), wire...), 0))
		if err == nil || !reflect.DeepEqual(before.Values(), m.Values()) {
			t.Fatal("BSON with trailing bytes must fail without changing the map")
		}
	}
}

func TestJSONReplacementAtomicity(t *testing.T) {
	m := safeorderedmap.New[int]().Add("stale", 1)

	err := json.Unmarshal([]byte(`{"fresh":2}`), m)
	if err != nil {
		t.Fatal(err)
	}

	m.Add("stale", 3)

	encoded, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}

	var values map[string]int

	err = json.Unmarshal(encoded, &values)
	if err != nil || !reflect.DeepEqual(values, map[string]int{"fresh": 2, "stale": 3}) || m.Size() != 2 {
		t.Fatalf("reused map lost an inserted value: %s, size=%d, err=%v", encoded, m.Size(), err)
	}

	err = m.UnmarshalJSON([]byte(`{"bad":`))
	if err == nil || !reflect.DeepEqual(m.Keys(), []string{"fresh", "stale"}) {
		t.Fatal("malformed JSON changed the map")
	}

	err = m.UnmarshalJSON([]byte(`null`))
	if err != nil {
		t.Fatal(err)
	}

	m.Add("after-null", 4)

	if !m.Contains("after-null") || m.Size() != 1 {
		t.Fatal("map cannot be reused after JSON null")
	}
}

type jsonReentry struct{ mutate func() }

func (v *jsonReentry) MarshalJSON() ([]byte, error) {
	v.mutate()
	return []byte(`1`), nil
}

func TestJSONMarshalerReentry(t *testing.T) {
	m := safeorderedmap.New[*jsonReentry]()
	m.Add("original", &jsonReentry{mutate: func() { m.Add("new", nil) }})

	var (
		encoded []byte
		err     error
	)

	finish(t, func() { encoded, err = json.Marshal(m) })

	if err != nil || string(encoded) != `{"original":1}` || m.Size() != 2 {
		t.Fatalf("JSON callback must use a snapshot: %s, err=%v", encoded, err)
	}
}

func TestBSONTypedContentsAndErrors(t *testing.T) {
	type payload struct {
		Number int64  `bson:"number"`
		Blob   []byte `bson:"blob"`
		Nested bson.D `bson:"nested"`
	}

	value := payload{Number: 1 << 40, Blob: []byte{0, 1}, Nested: bson.D{{Key: "inner", Value: int32(2)}}}
	m := safeorderedmap.New[payload]().Add("second", value).Add("first", value)

	wire, err := bson.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}

	decoded := safeorderedmap.New[payload]().Add("stale", value)

	err = bson.Unmarshal(wire, decoded)
	if err != nil || !reflect.DeepEqual(decoded.Values(), m.Values()) || !reflect.DeepEqual(decoded.Keys(), []string{"second", "first"}) {
		t.Fatalf("BSON lost typed values/order or retained stale data: keys=%v err=%v", decoded.Keys(), err)
	}

	wire, err = bson.Marshal(bson.D{{Key: "duplicate", Value: int32(1)}, {Key: "duplicate", Value: int32(2)}})
	if err != nil {
		t.Fatal(err)
	}

	ints := safeorderedmap.New[int32]()

	err = bson.Unmarshal(wire, ints)
	if err != nil || ints.Size() != 1 {
		t.Fatalf("duplicate BSON keys broke the data/order invariant: %v", err)
	}

	got, ok := ints.Get("duplicate")
	if !ok || got != 2 {
		t.Fatal("duplicate BSON key must keep its last value")
	}

	wire, err = bson.Marshal(bson.D{{Key: "bad", Value: "not an integer"}})
	if err != nil {
		t.Fatal(err)
	}

	err = bson.Unmarshal(wire, ints)
	if err == nil || ints.Size() != 1 || !ints.Contains("duplicate") {
		t.Fatal("BSON value decode error changed existing state")
	}

	_, err = safeorderedmap.New[any]().Add("unsupported", func() {}).MarshalBSON()
	if err == nil {
		t.Fatal("BSON must reject an unsupported value")
	}

	if safeorderedmap.New[any]().Add("unsupported", func() {}).String() != "" {
		t.Fatal("String must preserve its JSON encoding-error convention")
	}
}

func TestBSONWriterProgress(t *testing.T) {
	m := safeorderedmap.New[int32]().Add("stable", 1)

	finish(t, func() {
		var workers sync.WaitGroup
		workers.Go(func() {
			for range 1000 {
				m.Add("changing", 2).Delete("changing")
			}
		})

		for range 1000 {
			wire, err := m.MarshalBSON()
			if err != nil {
				t.Error(err)
				break
			}

			err = bson.Raw(wire).Validate()
			if err != nil {
				t.Error(err)
				break
			}
		}

		workers.Wait()
	})
}
