package argus

import (
	"errors"
	"testing"
)

var (
	benchmarkBreakerResult any
	benchmarkBreakerError  error
)

func BenchmarkBreaker_SuccessfulCall(b *testing.B) {
	breaker := NewBreaker("benchmark-success")
	successFunc := func() (any, error) {
		return "success", nil
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, err := breaker.Execute(successFunc)
		benchmarkBreakerResult = result
		benchmarkBreakerError = err
	}
}

func BenchmarkBreaker_SuccessfulCallWithWork(b *testing.B) {
	breaker := NewBreaker("benchmark-work")
	workFunc := func() (any, error) {
		sum := 0
		for i := 0; i < 100; i++ {
			sum += i
		}
		return sum, nil
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, err := breaker.Execute(workFunc)
		benchmarkBreakerResult = result
		benchmarkBreakerError = err
	}
}

func BenchmarkBreaker_OpenCircuit(b *testing.B) {
	breaker := NewBreaker("benchmark-open")

	failFunc := func() (any, error) {
		return nil, errors.New("service unavailable")
	}
	for i := 0; i < 4; i++ {
		breaker.Execute(failFunc)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, err := breaker.Execute(failFunc)
		benchmarkBreakerResult = result
		benchmarkBreakerError = err
	}
}

func BenchmarkBreaker_FailingCall(b *testing.B) {
	failFunc := func() (any, error) {
		return nil, errors.New("temporary failure")
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		breaker := NewBreaker("benchmark-fail")
		result, err := breaker.Execute(failFunc)
		benchmarkBreakerResult = result
		benchmarkBreakerError = err
	}
}

func BenchmarkBreaker_MixedTraffic(b *testing.B) {
	breaker := NewBreaker("benchmark-mixed")

	successFunc := func() (any, error) {
		return "ok", nil
	}

	failFunc := func() (any, error) {
		return nil, errors.New("occasional failure")
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result any
		var err error
		if i%10 == 0 {
			result, err = breaker.Execute(failFunc)
		} else {
			result, err = breaker.Execute(successFunc)
		}
		benchmarkBreakerResult = result
		benchmarkBreakerError = err
	}
}

func BenchmarkBreaker_Parallel(b *testing.B) {
	breaker := NewBreaker("benchmark-parallel")
	successFunc := func() (any, error) {
		return "parallel-success", nil
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			result, err := breaker.Execute(successFunc)
			benchmarkBreakerResult = result
			benchmarkBreakerError = err
		}
	})
}

func BenchmarkBreaker_NoWrapper(b *testing.B) {
	directFunc := func() (any, error) {
		return "direct", nil
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, err := directFunc()
		benchmarkBreakerResult = result
		benchmarkBreakerError = err
	}
}

func BenchmarkBreaker_AllocationsOnly(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		breaker := NewBreaker("benchmark-alloc")
		benchmarkBreakerResult = breaker
	}
}
