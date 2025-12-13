package argus

import (
	"errors"
	"testing"
)

func TestBreaker_FailsOpen(t *testing.T) {
	breaker := NewBreaker("test-breaker")

	failFunc := func() (any, error) {
		return nil, errors.New("ai service unavailable")
	}

	for i := range 4 {
		_, err := breaker.Execute(failFunc)
		if err == nil {
			t.Fatalf("iteration %d: expected error from failFunc, got nil", i)
		}
	}

	panicFunc := func() (any, error) {
		panic("breaker should not call this function when open!")
	}

	_, err := breaker.Execute(panicFunc)

	if err == nil {
		t.Fatal("Expected breaker error (open state), got nil")
	}

	if err.Error() != "circuit breaker is open" {
		t.Errorf("Expected 'circuit breaker is open' error, got: %v", err)
	}
}

func TestBreaker_SuccessReset(t *testing.T) {
	breaker := NewBreaker("success-breaker")

	successFunc := func() (any, error) {
		return "success", nil
	}

	res, err := breaker.Execute(successFunc)
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}

	if res.(string) != "success" {
		t.Errorf("Expected 'success', got %v", res)
	}
}
