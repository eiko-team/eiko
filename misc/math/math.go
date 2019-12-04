package math

// ScoreStore calculate the final score of a store
func ScoreStore(global, update, n int) int {
	return (global*n + update) / (n + 1)
}
