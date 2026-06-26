package batch_test

import (
	"errors"
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/batch"
)

func TestOf_EvenSplit(t *testing.T) {
	got := batch.Of([]int{1, 2, 3, 4}, 2)
	if len(got) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(got))
	}
	if len(got[0]) != 2 || len(got[1]) != 2 {
		t.Error("expected each chunk to have 2 items")
	}
}

func TestOf_UnevenSplit(t *testing.T) {
	got := batch.Of([]int{1, 2, 3, 4, 5}, 2)
	if len(got) != 3 {
		t.Fatalf("expected 3 chunks, got %d", len(got))
	}
	if len(got[2]) != 1 {
		t.Errorf("expected last chunk to have 1 item, got %d", len(got[2]))
	}
}

func TestOf_EmptyInput(t *testing.T) {
	got := batch.Of([]string{}, 5)
	if got != nil {
		t.Errorf("expected nil for empty input, got %v", got)
	}
}

func TestOf_SizeLargerThanItems(t *testing.T) {
	got := batch.Of([]int{1, 2}, 100)
	if len(got) != 1 || len(got[0]) != 2 {
		t.Errorf("expected single chunk of 2, got %v", got)
	}
}

func TestDo_ProcessesAllChunks(t *testing.T) {
	var processed []int
	err := batch.Do([]int{1, 2, 3, 4, 5}, 2, func(chunk []int) error {
		processed = append(processed, chunk...)
		return nil
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(processed) != 5 {
		t.Errorf("expected 5 items processed, got %d", len(processed))
	}
}

func TestDo_StopsOnError(t *testing.T) {
	calls := 0
	errStop := errors.New("stop")
	err := batch.Do([]int{1, 2, 3, 4, 5, 6}, 2, func(chunk []int) error {
		calls++
		if calls == 2 {
			return errStop
		}
		return nil
	})
	if !errors.Is(err, errStop) {
		t.Errorf("expected errStop, got %v", err)
	}
	if calls != 2 {
		t.Errorf("expected fn to be called 2 times, got %d", calls)
	}
}

func TestMap_CollectsResults(t *testing.T) {
	items := []int{1, 2, 3, 4}
	result, err := batch.Map(items, 2, func(chunk []int) ([]int, error) {
		out := make([]int, len(chunk))
		for i, v := range chunk {
			out[i] = v * 2
		}
		return out, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 4 {
		t.Fatalf("expected 4 results, got %d", len(result))
	}
	for i, v := range result {
		if v != (i+1)*2 {
			t.Errorf("result[%d] = %d, want %d", i, v, (i+1)*2)
		}
	}
}
