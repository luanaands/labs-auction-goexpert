package auction

import (
	"context"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection
	Updater    AuctionUpdater
}

var auctionTimer = func(d time.Duration) *time.Timer {
	return time.NewTimer(d)
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	ar := &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
	ar.Updater = ar
	return ar
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	ar.StartAuctionCloser(context.Background(), auctionEntity.Id)

	return nil
}

func (ar *AuctionRepository) StartAuctionCloser(ctx context.Context, auctionId string) {
	go func() {
		duration := getAuctionDuration()
		timer := auctionTimer(duration)
		defer timer.Stop()

		select {
		case <-timer.C:
			logger.Info(fmt.Sprintf("O leilão %s foi encerrado automaticamente", auctionId))
			if err := ar.Updater.UpdateAuctionStatus(
				ctx,
				auctionId,
				auction_entity.Completed); err != nil {
				logger.Error(fmt.Sprintf("Erro ao fechar o leilão %s", auctionId), err)
				return
			}

		case <-ctx.Done():
			logger.Info(fmt.Sprintf("Contexto foi cancelado para o leilão %s", auctionId))
			return
		}
	}()
}

func getAuctionDuration() time.Duration {
	auctionDuration := os.Getenv("AUCTION_DURATION")
	duration, err := time.ParseDuration(auctionDuration)
	if err != nil {
		return time.Minute * 8
	}

	return duration
}
