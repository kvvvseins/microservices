package repository

import (
	"github.com/google/uuid"
	"github.com/kvvvseins/mictoservices/services/store/internal/app/model"
	"gorm.io/gorm"
)

type StoreRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) *StoreRepository {
	return &StoreRepository{db: db}
}

// FindByGuids найти по guids
func (ur *StoreRepository) FindByGuids(guids []string) ([]model.Store, error) {
	var stores []model.Store

	if len(guids) == 0 {
		return stores, nil
	}

	if err := ur.db.
		Select(
			"store.id",
			"store.guid",
			"store.name",
			"store.price",
			"SUM(reserve.quantity) AS quantity",
		).
		Where("store.guid IN ?", guids).
		Joins("INNER JOIN reserve ON reserve.store_id = store.id").
		Group("store.id, store.guid, store.name, store.price").
		Limit(10).
		Find(&stores).Error; err != nil {
		return nil, err
	}

	return stores, nil
}

// GetAll выбрать все
func (ur *StoreRepository) GetAll(limit int) ([]model.Store, error) {
	var stores []model.Store

	if err := ur.db.
		Select(
			"store.id",
			"store.guid",
			"store.name",
			"store.price",
			"SUM(reserve.quantity) AS quantity",
		).
		Joins("INNER JOIN reserve ON reserve.store_id = store.id").
		Group("store.id, store.guid, store.name, store.price").
		Limit(limit).
		Find(&stores).Error; err != nil {
		return nil, err
	}

	return stores, nil
}

// FindByGuid найти по guid
func (ur *StoreRepository) FindByGuid(guid string) (*model.Store, error) {
	guids := []string{guid}

	stores, err := ur.FindByGuids(guids)
	if err != nil {
		return nil, err
	}

	if len(stores) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &stores[0], nil
}

// FindByName найти по name
func (ur *StoreRepository) FindByName(name string) (*model.Store, error) {
	var store model.Store

	if err := ur.db.Where("name = ?", name).First(&store).Error; err != nil {
		return nil, err
	}

	return &store, nil
}

func (ur *StoreRepository) Reserve(
	tx *gorm.DB,
	quantity int,
	storeId uint,
	orderGuid *uuid.UUID,
	userId uuid.UUID,
) (*model.Reserve, error) {
	reserve := &model.Reserve{
		StoreId:  storeId,
		Quantity: quantity,
		OrderId:  orderGuid,
		UserID:   userId,
	}

	result := tx.Create(reserve)
	if result.Error != nil {
		return nil, result.Error
	}

	return reserve, nil
}

func (ur *StoreRepository) GetBalance(storeId uint) (int, error) {
	var quantity int

	balance := model.Reserve{}

	result := ur.db.Table(balance.TableName()).
		Select("SUM(quantity) AS quantity").
		Where("store_id = ?", storeId).
		Scan(&quantity)
	if result.Error != nil {
		return 0, result.Error
	}

	return quantity, nil
}
