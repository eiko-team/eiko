package consumable_test

import (
	"testing"

	"github.com/eiko-team/eiko/misc/consumable"
	"github.com/eiko-team/eiko/misc/structures"
)

func TestIsBoycotted(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name   string
		c      structures.Consumable
		user   structures.User
		result bool
	}{
		{"sanity", structures.Consumable{}, structures.User{}, false},
		{"true", structures.Consumable{Company: "Test"},
			structures.User{SBoycott: []string{"Test"}}, true},
		{"false", structures.Consumable{Company: "Test"},
			structures.User{SBoycott: []string{"not Test"}}, false},
		{"empty boycott array", structures.Consumable{Company: "Test"},
			structures.User{}, false},
		{"empty company", structures.Consumable{},
			structures.User{SBoycott: []string{"Test"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.result != consumable.IsBoycotted(tt.c, tt.user) {
				t.Errorf("IsBoycotted(%+v, %+v) got %t want %t",
					tt.c, tt.user, consumable.IsBoycotted(tt.c, tt.user), tt.result)
			}
		})
	}
}
