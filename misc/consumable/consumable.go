package consumable

import (
	"github.com/eiko-team/eiko/misc/structures"
)

// IsBoycotted check if the consumable is boycotted by the user
func IsBoycotted(c structures.Consumable, user structures.User) bool {
	for _, boycotted := range user.SBoycott {
		if boycotted == c.Company {
			return true
		}
	}
	return false
}

// IsBio return true is the consumable is Bio
func IsBio(structures.Consumable) bool {
	// TODO: Create a function `IsBio` to filter consumables
	return true
}

// IsVegan return true is the consumable is Vegan
func IsVegan(structures.Consumable) bool {
	// TODO: Create a function `IsVegan` to filter consumables
	return true
}

// IsHalal return true is the consumable is Halal
func IsHalal(structures.Consumable) bool {
	// TODO: Create a function `IsHalal` to filter consumables
	return true
}

// IsCasher return true is the consumable is Casher
func IsCasher(structures.Consumable) bool {
	// TODO: Create a function `IsCasher` to filter consumables
	return true
}

// ContainSodium return true is the consumable contain salt
func ContainSodium(c structures.Consumable) bool {
	return c.Sodium >= 0
}

// ContainEgg return true is the consumable contain eggs
func ContainEgg(structures.Consumable) bool {
	// TODO: Create a function `ContainEgg` to filter consumables
	return true
}

// ContainPenut return true is the consumable contain penut
func ContainPenut(structures.Consumable) bool {
	// TODO: Create a function `ContainPenut` to filter consumables
	return true
}

// ContainCrustace return true is the consumable contain crustace
func ContainCrustace(structures.Consumable) bool {
	// TODO: Create a function `ContainCrustace` to filter consumables
	return true
}

// IsGlutenFree return true is the consumable is guten free
func IsGlutenFree(structures.Consumable) bool {
	// TODO: Create a function `IsGlutenFree` to filter consumables
	return true
}

// ForDiabetique return true is the consumable is suited for diabetique
func ForDiabetic(c structures.Consumable) bool {
	return c.Glucides <= 0
}
