package purchase

import (
	coffeeco "coffeeco/internal"
	"coffeeco/internal/payment"
	store "coffeeco/internal/store"
	"errors"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
)

type Purchase struct {
	id                 uuid.UUID
	Store              store.Store
	ProductsToPurchase []coffeeco.Product
	total              money.Money
	PaymentMeans       payment.Means
	timeOfPurchase     time.Time
	CardToken          *string
}

func (p *Purchase) validateAndEnrich() error {
	if len(p.ProductsToPurchase) == 0 {
		return errors.New("no products to purchase")
	}
	p.total = *money.New(0, "USD")

	for _, product := range p.ProductsToPurchase {
		newTotal, _ := p.total.Add(&product.BasePrice)
		p.total = *newTotal
	}

}
