package math

import "testing"

func TestScoreStore(t *testing.T) {
	type args struct {
		global int
		update int
		n      int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"sanity", args{1, 1, 1}, 1},
		{"simple notation", args{50, 100, 1}, 75},
		{"simple notation", args{95, 10, 1100}, 94},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ScoreStore(tt.args.global, tt.args.update, tt.args.n)
			if got != tt.want {
				t.Errorf("ScoreStore() = %v, want %v", got, tt.want)
			}
		})
	}
}
