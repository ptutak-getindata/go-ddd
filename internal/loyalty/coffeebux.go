package loyalty

import (
	coffeeco "coffeeco/internal"
	"coffeeco/internal/store"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type CoffeeBux struct {
	ID                                    uuid.UUID
	store                                 store.Store
	coffeeLover                           coffeeco.CoffeeLover
	FreeDrinksAvailable                   int
	RemainingDrinkPurchasesUntilFreeDrink int
}

func (cb *CoffeeBux) AddStamp() {
	if cb.RemainingDrinkPurchasesUntilFreeDrink == 1 {
		cb.RemainingDrinkPurchasesUntilFreeDrink = 10
		cb.FreeDrinksAvailable++
	} else {
		cb.RemainingDrinkPurchasesUntilFreeDrink--
	}
}

func (cb *CoffeeBux) Pay(ctx context.Context, purchases []coffeeco.Product) error {
	allPurchases := len(purchases)
	if allPurchases == 0 {
		return errors.New("no purchases to pay")
	}

	if cb.FreeDrinksAvailable < allPurchases {
		return fmt.Errorf("not enough coffeeBux to cover entire purchase. Have %d, need %d", len(purchases), cb.FreeDrinksAvailable)
	}

	cb.FreeDrinksAvailable -= allPurchases
	return nil
}
