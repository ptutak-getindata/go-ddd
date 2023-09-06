package purchase

import (
	coffeeco "coffeeco/internal"
	"coffeeco/internal/payment"
	"coffeeco/internal/store"
	"context"
	"fmt"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Store(ctx context.Context, purchase *Purchase) error
}

type MongoRepository struct {
	purchases *mongo.Collection
}

func NewMongoRepo(ctx context.Context, connectionString string) (*MongoRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}
	return &MongoRepository{
		purchases: client.Database("coffeeco").Collection("purchases"),
	}, nil
}

func (mongoRepo *MongoRepository) Store(ctx context.Context, purchase Purchase) error {
	mongoDocument := toMongoPurchase(purchase)
	_, err := mongoRepo.purchases.InsertOne(ctx, mongoDocument)
	if err != nil {
		return fmt.Errorf("failed to persist purchase: %w", err)
	}
	return nil
}

type mongoPurchase struct {
	id                 uuid.UUID
	store              store.Store
	productsToPurchase []coffeeco.Product
	total              money.Money
	paymentMeans       payment.Means
	timeOfPurchase     time.Time
	cardToken          *string
}

func toMongoPurchase(purchase Purchase) mongoPurchase {
	return mongoPurchase{
		id:                 purchase.id,
		store:              purchase.Store,
		productsToPurchase: purchase.ProductsToPurchase,
		total:              purchase.total,
		paymentMeans:       purchase.PaymentMeans,
		timeOfPurchase:     purchase.timeOfPurchase,
		cardToken:          purchase.CardToken,
	}
}
