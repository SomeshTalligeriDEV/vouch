package service

import "testing"

// BenchmarkProblemService_Create measures the overhead of problem creation
// (slug generation, validation) with an in-memory repo.
func BenchmarkProblemService_Create(b *testing.B) {
	svc := NewProblemService(newFakeProblemRepo(), &recordingEnqueuer{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.Create(ctx(), "poster1", ProblemInput{
			Title:       "Benchmark problem title here",
			Description: "Some description that is long enough to be valid for the benchmark run.",
			BudgetMin:   100,
			BudgetMax:   500,
		})
	}
}
