package hash

import (
	"reflect"
	"testing"
)

func Test_saltPassword(t *testing.T) {
	type args struct {
		pass string
		salt string
	}
	tests := []struct {
		name string
		args args
	}{
		{"test1", args{"aze", "aze"}},
		{"test2", args{"qsd", "qsd"}},
		{"test3", args{"aze", "qsd"}},
		{"test4", args{"wxc", "qsd"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1 := saltPassword(tt.args.pass, tt.args.salt)
			got2 := saltPassword(tt.args.pass, tt.args.salt)
			if !reflect.DeepEqual(got1, got2) {
				t.Errorf("saltPassword() = %v, want %v", got1, got2)
			}
		})
	}
}

func TestHash(t *testing.T) {
	type args struct {
		pass string
		salt string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{"aze", "aze"}, false},
		{"test2", args{"qsd", "qsd"}, false},
		{"test3", args{"aze", "qsd"}, false},
		{"test4", args{"wxc", "qsd"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Hash(tt.args.pass, tt.args.salt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := CompareHash(got, tt.args.pass, tt.args.salt); err != nil {
				t.Errorf("Hash() = %v, err: %v", got, err)
			}
		})
	}
}

func TestGenerateKey(t *testing.T) {
	tests := []struct {
		name string
		n    int
	}{
		// TODO: Add test cases.
		{"simple", 42},
		{"medium", 21},
		{"short", 1},
		{"empty", 0},
		{"long", 9999},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(GenerateKey(tt.n)); got != tt.n {
				t.Errorf("len(GenerateKey()) = %v, want %v", got, tt.n)
			}
		})
	}
}
