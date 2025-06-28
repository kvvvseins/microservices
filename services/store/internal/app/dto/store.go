package dto

import "github.com/kvvvseins/mictoservices/services/store/internal/app/model"

type ViewStore struct {
	Guid     string `json:"guid"`
	Name     string `json:"name"`
	Price    uint   `json:"price"`
	Quantity int    `json:"quantity"`
}

type CreateStore struct {
	Price    uint   `json:"price"`
	Name     string `json:"name"`
	Quantity uint   `json:"quantity"`
}

func BuildViewsFromModels(stores []model.Store) []ViewStore {
	storesDto := make([]ViewStore, 0, len(stores))
	for _, store := range stores {
		storesDto = append(storesDto, ViewStore{
			Guid:     store.Guid.String(),
			Name:     store.Name,
			Price:    store.Price,
			Quantity: store.Quantity,
		})
	}

	return storesDto
}
