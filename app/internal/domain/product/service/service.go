package service

import (
	"context"
	"github.com/ilkinabd/goods-manager/app/internal/domain/product/filter"

	"github.com/ilkinabd/goods-manager/app/internal/domain/product/dao"
	"github.com/ilkinabd/goods-manager/app/internal/domain/product/model"
	"github.com/ilkinabd/goods-manager/app/pkg/errors"
)

type ProductService struct {
	repository dao.ProductDAO
}

func NewProductService(repository dao.ProductDAO) *ProductService {
	return &ProductService{repository: repository}
}

func (s *ProductService) All(ctx context.Context, filtering []filter.Criteria, sorting filter.Sortable) ([]*model.Product, error) {
	dbProducts, err := s.repository.All(ctx, filtering, sorting)
	if err != nil {
		return nil, errors.Wrap(err, "repository.All")
	}

	var products []*model.Product
	for _, dbP := range dbProducts {
		products = append(products, model.NewProductFromDAO(dbP))
	}

	return products, nil
}

func (s *ProductService) Create(ctx context.Context, product *model.Product) (*model.Product, error) {
	productStorageMap, err := product.ToMap()
	if err != nil {
		return nil, err
	}

	err = s.repository.Create(ctx, productStorageMap)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) One(ctx context.Context, id string) (*model.Product, error) {
	one, err := s.repository.One(ctx, id)
	if err != nil {
		return nil, err
	}

	return model.NewProductFromDAO(one), nil
}

func (s *ProductService) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

func (s *ProductService) Update(ctx context.Context, product *model.Product) error {
	productStorageMap, err := product.ToMap()
	if err != nil {
		return err
	}

	return s.repository.Update(ctx, product.ID, productStorageMap)
}
