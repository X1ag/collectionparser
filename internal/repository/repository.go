package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func NewRepository(DB *pgxpool.Pool) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) InsertAndGetId(ctx context.Context, floor int, collectionName string) (int, error) {
	var id int
	err := r.DB.QueryRow(ctx, `INSERT INTO floor_price (collectionName, floor)
        VALUES ($1, $2)
        ON CONFLICT (collectionName)
        DO UPDATE SET floor = EXCLUDED.floor, created_at = CURRENT_TIMESTAMP
        RETURNING id`, collectionName, floor).Scan(&id)
	if err != nil {
		log.Default().Printf("Error inserting floor: %v", err)
		return -1, err
	}
	log.Default().Println("id", id)

	return id, nil
}

func (r *Repository) GetFloorByCollectionName(ctx context.Context, collectionName string) (int, error) {
	var floor int
	err := r.DB.QueryRow(ctx, "SELECT floor FROM floor_price WHERE collectionName = $1", collectionName).Scan(&floor)
	if err != nil {
		log.Default().Printf("Error getting floor: %v", err)
		return -1, err
	}

	log.Default().Printf("floor of %s is %d", collectionName, floor)

	return floor, nil
}