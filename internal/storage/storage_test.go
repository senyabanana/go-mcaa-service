package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMemStorage(t *testing.T) {
	ms := NewMemStorage()
	assert.NotNil(t, ms, "expected MemStorage instance, got nil")
	assert.NotNil(t, ms.gauges, "expected gauges map to be initialized, got nil")
	assert.NotNil(t, ms.counters, "expected counters map to be initialized, got nil")
}

func TestUpdateGauge(t *testing.T) {
	tests := []struct {
		name      string
		gaugeName string
		value     float64
		want      float64
	}{
		{
			name:      "UpdateGauge with positive value",
			gaugeName: "testGauge1",
			value:     123.45,
			want:      123.45,
		},
		{
			name:      "UpdateGauge with negative value",
			gaugeName: "testGauge2",
			value:     -2.4,
			want:      -2.4,
		},
		{
			name:      "UpdateGauge with zero value",
			gaugeName: "testGauge3",
			value:     0,
			want:      0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := NewMemStorage()
			ms.UpdateGauge(tt.gaugeName, tt.value)
			val, ok := ms.gauges[tt.gaugeName]
			assert.True(t, ok, "expected gauge '%s' to be present", tt.gaugeName)
			assert.Equal(t, tt.want, val, "expected gauge '%s' to be %f, got %f", tt.gaugeName, tt.want, val)
		})
	}
}

func TestUpdateCounter(t *testing.T) {
	tests := []struct {
		name        string
		counterName string
		value       int64
		want        int64
	}{
		{
			name:        "UpdateCounter with positive value",
			counterName: "testCounter1",
			value:       33,
			want:        33,
		},
		{
			name:        "UpdateCounter with negative value",
			counterName: "testCounter2",
			value:       -22,
			want:        -22,
		},
		{
			name:        "UpdateCounter with zero value",
			counterName: "testCounter3",
			value:       0,
			want:        0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := NewMemStorage()
			ms.UpdateCounter(tt.counterName, tt.value)
			val, ok := ms.counters[tt.counterName]
			assert.True(t, ok, "expected counter '%s' to be present", tt.counterName)
			assert.Equal(t, tt.want, val, "expected counter '%s' to be %d, got %d", tt.counterName, tt.want, val)
		})
	}
}
