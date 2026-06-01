package auction

import (
	"context"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"

	"go.mongodb.org/mongo-driver/bson"
)

type AuctionUpdater interface {
	UpdateAuctionStatus(ctx context.Context, auctionId string, status auction_entity.AuctionStatus) *internal_error.InternalError
}

func (ar *AuctionRepository) UpdateAuctionStatus(
	ctx context.Context,
	auctionId string,
	status auction_entity.AuctionStatus) *internal_error.InternalError {

	filter := bson.M{"_id": auctionId}
	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}

	result, err := ar.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error(fmt.Sprintf("Erro ao atualizar o status do leilão %s", auctionId), err)
		return internal_error.NewInternalServerError("Erro ao atualizar o status do leilão")
	}

	if result.MatchedCount == 0 {
		logger.Error(fmt.Sprintf("Leilão não encontrado para o id = %s", auctionId), nil)
		return internal_error.NewNotFoundError("leilão não encontrado")
	}

	return nil
}
