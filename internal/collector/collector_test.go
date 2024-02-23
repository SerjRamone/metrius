package collector

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/SerjRamone/metrius/internal/metrics"
)

func TestCollector_Collect(t *testing.T) {
	c := New()
	c.Collect()

	if len(c.collections) < 1 {
		t.Error("empty collections slice")
	}

	if c.collections[0] == nil {
		t.Error("first collection is nil")
	}
}

func TestCollector_Export(t *testing.T) {
	c := New()
	mockData := metrics.NewCollection(runtime.MemStats{})
	c.collections = append(c.collections, mockData)
	collectionsCopy := make([]metrics.Collection, len(c.collections))
	copy(collectionsCopy, c.collections)
	exported := c.Export()

	if !reflect.DeepEqual(exported, collectionsCopy) {
		t.Error("exported collections is not equal to imported")
	}

	if len(c.collections) != 0 {
		t.Error("collections is not cleared after export")
	}
}
