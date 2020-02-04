package user

import (
	"github.com/eiko-team/eiko/misc/consumable"
	"github.com/eiko-team/eiko/misc/structures"
)

// IsGood checks if the consumable is sutable for the user
func IsGood(user structures.User, c structures.Consumable) bool {
	return (user.SBio && !consumable.IsBio(c)) ||
		(user.SVegan && !consumable.IsVegan(c)) ||
		(user.SHalal && !consumable.IsHalal(c)) ||
		(user.SCasher && !consumable.IsCasher(c)) ||
		(user.SSel && !consumable.ContainSel(c)) ||
		(user.SOeuf && !consumable.ContainOeuf(c)) ||
		(user.SArachide && !consumable.ContainArachide(c)) ||
		(user.SCrustace && !consumable.ContainCrustace(c)) ||
		(user.SGluten && !consumable.IsGlutenFree(c)) ||
		(user.SDiabetique && !consumable.ForDiabetique(c))
}
