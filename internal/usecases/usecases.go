package usecases

import (
	"context"
	"floor_parser/internal/repository"

	"github.com/gin-gonic/gin"
)

func GetFloorPrice(ctx *gin.Context, repository *repository.Repository, collectionName string) (int, error) {
	return repository.GetFloorByCollectionName(ctx, collectionName)
}

func InsertPrice(ctx context.Context, repository *repository.Repository, collectionName string, floorPrice int) (int, error) {
	id, err := repository.InsertAndGetId(ctx, floorPrice, collectionName)

	return id, err
}