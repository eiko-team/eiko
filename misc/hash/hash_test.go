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
		{"empty pass", args{"", "qsd"}, false},
		{"empty hash", args{"abc", ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Hash(tt.args.pass, tt.args.salt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !CompareHash(got, tt.args.pass, tt.args.salt) {
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

var (
	salt    = "misc.GenerateKey(999)"
	pass    = "pass"
	hash, _ = Hash(pass, salt)
)

func TestCompareHash(t *testing.T) {
	type args struct {
		hash string
		pass string
		salt string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"sanity", args{hash, pass, salt}, true},
		{"wrong password", args{hash, pass + "wrong", salt}, false},
		// Wierd
		// {"wrong hash back", args{hash + "wronghash", pass, salt}, false},
		{"wrong hash front", args{"test" + hash, pass, salt}, false},
		{"wrong hash both", args{"test" + hash + "wrong", pass, salt}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("%+v", tt.args)
			got := CompareHash(tt.args.hash, tt.args.pass, tt.args.salt)
			if got != tt.want {
				t.Errorf("CompareHash() = %v, want = %v", got, tt.want)
			}
		})
	}
}
