package structures

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func ExampleMergeUser() {
	UserWithMissingInformations := User{
		Email: "test@test.test",
		ID:    42,
	}

	SameUserWithMoreInformations := User{
		Pass:      "Pass",
		Created:   time.Now(),
		Validated: true,
		ID:        21,
	}
	NewUser := MergeUser(UserWithMissingInformations,
		SameUserWithMoreInformations)

	fmt.Println(NewUser.Email)
	fmt.Println(NewUser.Pass)
	fmt.Println(NewUser.Validated)
	fmt.Println(NewUser.ID)

	// Output:
	// test@test.test
	// Pass
	// true
	// 42
}

func TestMergeUser(t *testing.T) {
	fullUser := User{
		Email:     "test",
		Pass:      "Pass",
		Created:   time.Now(),
		Validated: true,
		ID:        42,
	}
	tests := []struct {
		name  string
		user1 User
		user2 User
		want  User
	}{
		{"sanity", User{}, User{}, User{}},
		{"no field", User{Email: "test"}, User{}, User{Email: "test"}},
		{"same field", User{Email: "test"}, User{Email: "not test"}, User{Email: "test"}},
		{"same field not set", User{}, User{Email: "test"}, User{Email: "test"}},
		{"full fields", fullUser, User{}, fullUser},
		{"none fields set", User{}, fullUser, fullUser},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeUser(tt.user1, tt.user2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User.Merge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleMergeStore() {
	StoreWithMissingInformations := Store{
		Name: "Store Test.",
		ID:   42,
	}

	SameStoreWithMoreInformations := Store{
		Name:    "The Best Test Store",
		ID:      21,
		Country: "In",
	}
	NewStore := MergeStore(StoreWithMissingInformations,
		SameStoreWithMoreInformations)

	fmt.Println(NewStore.Name)
	fmt.Println(NewStore.ID)
	fmt.Println(NewStore.Country)

	// Output:
	// The Best Test Store
	// 42
	// In
}

func TestMergeStore(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name string
		i1   Store
		i2   Store
		want Store
	}{
		{"sanity", Store{}, Store{}, Store{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeStore(tt.i1, tt.i2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeStore() = %v, want %v", got, tt.want)
			}
		})
	}
}
