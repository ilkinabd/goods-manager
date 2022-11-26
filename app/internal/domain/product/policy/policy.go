package policy

import (
	"context"
	filter2 "github.com/ilkinabd/goods-manager/app/internal/domain/product/filter"

	"github.com/ilkinabd/goods-manager/app/internal/domain/product/model"
	"github.com/ilkinabd/goods-manager/app/internal/domain/product/service"
	"github.com/ilkinabd/goods-manager/app/pkg/errors"
)

type ProductPolicy struct {
	productService *service.ProductService
}

func NewProductPolicy(productService *service.ProductService) *ProductPolicy {
	return &ProductPolicy{productService: productService}
}

func (p *ProductPolicy) All(ctx context.Context, filtering []filter2.Criteria, sorting filter2.Sortable) ([]*model.Product, error) {
	products, err := p.productService.All(ctx, filtering, sorting)
	if err != nil {
		return nil, errors.Wrap(err, "productService.All")
	}

	return products, nil
}

func (p *ProductPolicy) CreateProduct(ctx context.Context, product *model.Product) (*model.Product, error) {
	return p.productService.Create(ctx, product)
}

func (p *ProductPolicy) One(ctx context.Context, id string) (*model.Product, error) {
	return p.productService.One(ctx, id)
}

func (p *ProductPolicy) Delete(ctx context.Context, id string) error {
	return p.productService.Delete(ctx, id)
}

func (p *ProductPolicy) Update(ctx context.Context, product *model.Product) error {
	return p.productService.Update(ctx, product)
}
