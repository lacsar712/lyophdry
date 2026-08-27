package store

import (
	"testing"

	"github.com/lacsar712/lyophdry/internal/model"
)

func TestCase(t *testing.T) {
	orig := LyoSnapshot{
		Tower: model.TowerID("tower-a1"),
		Segments: []SegmentSnapshot{
			{Zone: model.ZoneID("tower-a1-zone-00"), TempC: 42.0},
		},
	}
	clone := CloneLyoSnapshot(orig)
	clone.Segments[0].TempC = 99.0
	if orig.Segments[0].TempC == 99.0 {
		t.Fatal("clone mutated original batch Segments backing array")
	}
}
