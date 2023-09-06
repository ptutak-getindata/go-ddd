package purchase

import (
	coffeeco "coffeeco/internal"
	"coffeeco/internal/loyalty"
	"coffeeco/internal/payment"
	"coffeeco/internal/store"
	"context"
	"errors"
	"fmt"
	"log"
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
	if p.total.IsZero() {
		return errors.New("total purchase cannot be zero")
	}
	p.id = uuid.New()
	p.timeOfPurchase = time.Now()
	return nil
}

type CardChargeService interface {
	ChargeCard(ctx context.Context, amount money.Money, cardToken string) error
}

type Service struct {
	cardService  CardChargeService
	purchaseRepo Repository
}

func (s Service) CompletePurchase(ctx context.Context, purchase *Purchase, coffeeBuxCard *loyalty.CoffeeBux) error {
	if err := purchase.validateAndEnrich(); err != nil {
		return err
	}
	switch purchase.PaymentMeans {
	case payment.MEANS_CARD:
		if err := s.cardService.ChargeCard(ctx, purchase.total, *purchase.CardToken); err != nil {
			return errors.New("error charging card")
		}
	case payment.MEANS_CASH:
		log.Println("cash payment")
	case payment.MEANS_COFFEEBUX:
		if err := coffeeBuxCard.Pay(ctx, purchase.ProductsToPurchase); err != nil {
			return fmt.Errorf("error paying with coffeeBux: %w", err)
		}
	default:
		return errors.New("invalid payment means")
	}
	if err := s.purchaseRepo.Store(ctx, purchase); err != nil {
		return errors.New("error storing purchase")
	}
	if coffeeBuxCard != nil {
		coffeeBuxCard.AddStamp()
	}
	return nil
}
