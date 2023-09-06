package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrNoDiscount = errors.New("no discount available")

type Repository interface {
	GetStoreDiscount(ctx context.Context, storeID uuid.UUID) (int64, error)
}

type MongoRepository struct {
	storeDiscounts *mongo.Collection
}

func NewMongoRepo(ctx context.Context, connectionString string) (*MongoRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo: %w", err)
	}
	discounts := client.Database("coffeeco").Collection("store_discounts")
	return &MongoRepository{
		storeDiscounts: discounts,
	}, nil
}

func (mr MongoRepository) GetStoreDiscount(ctx context.Context, storeID uuid.UUID) (int64, error) {
	var discount int64
	if err := mr.storeDiscounts.FindOne(ctx, bson.D{{"store_id", storeID.String()}}).Decode(&discount); err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, ErrNoDiscount
		}
		return 0, fmt.Errorf("failed to get store discount: %w", err)
	}
	return discount, nil
}

func (m MongoRepository) Ping(ctx context.Context) error {
	if _, err := m.storeDiscounts.EstimatedDocumentCount(ctx); err != nil {
		return fmt.Errorf("failed to ping DB: %w", err)
	}
	return nil
}
