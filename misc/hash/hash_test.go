package hash

import (
	"reflect"
	"testing"
)

var (
	pass    = "pass"
	hash, _ = Hash(pass)
)

func TestHashExample(t *testing.T) {
	pass := "pass"
	hash, err := Hash(pass)
	if err != nil {
		t.Errorf("Hash() error = %v", err)
	}
	if !CompareHash(hash, pass) {
		t.Error("CompareHash() = false")
	}
}

func Test_saltPassword(t *testing.T) {
	type args struct {
		pass string
	}
	tests := []struct {
		name string
		args args
	}{
		{"test1", args{"aze"}},
		{"test2", args{"qsd"}},
		{"test3", args{"aze"}},
		{"test4", args{"wxc"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1 := saltPassword(tt.args.pass)
			got2 := saltPassword(tt.args.pass)
			if !reflect.DeepEqual(got1, got2) {
				t.Errorf("saltPassword() = %v, want %v", got1, got2)
			}
		})
	}
}

func TestHash(t *testing.T) {
	type args struct {
		pass string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{"aze"}, false},
		{"test2", args{"qsd"}, false},
		{"test3", args{"aze"}, false},
		{"test4", args{"wxc"}, false},
		{"empty pass", args{""}, false},
		{"empty hash", args{"abc"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Hash(tt.args.pass)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !CompareHash(got, tt.args.pass) {
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

func TestCompareHash(t *testing.T) {
	type args struct {
		hash string
		pass string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"sanity", args{hash, pass}, true},
		{"wrong password back", args{hash, pass + "wrong"}, false},
		{"wrong password front", args{hash, "wrong" + pass}, false},
		{"wrong password both", args{hash, "wrong" + pass + "wrong"}, false},
		// Weird
		// {"wrong hash back", args{hash + "wronghash", pass}, false},
		{"wrong hash front", args{"test" + hash, pass}, false},
		{"wrong hash both", args{"test" + hash + "wrong", pass}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("%+v", tt.args)
			got := CompareHash(tt.args.hash, tt.args.pass)
			if got != tt.want {
				t.Errorf("CompareHash() = %v, want = %v", got, tt.want)
			}
		})
	}
}
