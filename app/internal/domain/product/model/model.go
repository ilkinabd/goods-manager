package model

import (
	"encoding/json"
	"github.com/google/uuid"
	pbProducts "github.com/ilkinabd/goods-contracts/gen/go/products/v1"
	"github.com/ilkinabd/goods-manager/app/internal/domain/product/dao"
	"github.com/ilkinabd/goods-manager/app/pkg/errors"
	"github.com/ilkinabd/goods-manager/app/pkg/logging"
	"github.com/mitchellh/mapstructure"
	"time"
)

type Product struct {
	ID            string                 `mapstructure:"id"`
	Name          string                 `mapstructure:"name"`
	Description   string                 `mapstructure:"description"`
	ImageID       *string                `mapstructure:"image_id"`
	Price         uint64                 `mapstructure:"price"`
	CurrencyID    uint32                 `mapstructure:"currency_id"`
	Rating        uint32                 `mapstructure:"rating"`
	CategoryID    uint32                 `mapstructure:"category_id"`
	Specification map[string]interface{} `mapstructure:"specification"`
	CreatedAt     time.Time              `mapstructure:"created_at"`
	UpdatedAt     *time.Time             `mapstructure:"updated_at"`
}

func (p *Product) ToMap() (map[string]interface{}, error) {
	var updateProductMap map[string]interface{}
	err := mapstructure.Decode(p, &updateProductMap)
	if err != nil {
		return updateProductMap, errors.Wrap(err, "mapstructure.Decode(product)")
	}

	return updateProductMap, nil
}

func (p *Product) UpdateFromPB(productPB *pbProducts.UpdateProductRequest) {
	if productPB.Name != nil {
		p.Name = productPB.GetName()
	}
	if productPB.Description != nil {
		p.Description = productPB.GetDescription()
	}
	if productPB.ImageId != nil {
		p.ImageID = productPB.ImageId
	}
	if productPB.Price != nil {
		p.Price = productPB.GetPrice()
	}
	if productPB.CurrencyId != nil {
		p.CurrencyID = productPB.GetCurrencyId()
	}
	if productPB.Rating != nil {
		p.Rating = productPB.GetRating()
	}
	if productPB.CategoryId != nil {
		p.CategoryID = productPB.GetCategoryId()
	}

}

func (p *Product) ToProto() *pbProducts.Product {
	var updatedAt int64
	if p.UpdatedAt != nil {
		updatedAt = p.UpdatedAt.UnixMilli()
	}

	specBytes, err := json.Marshal(p.Specification)
	if err != nil {
		logging.GetLogger().Warnf("failed to marshal product specification %v", err)
		logging.GetLogger().Trace(p.Specification)
	}

	return &pbProducts.Product{
		Id:            p.ID,
		Name:          p.Name,
		Description:   p.Description,
		ImageId:       p.ImageID,
		Price:         p.Price,
		CurrencyId:    p.CurrencyID,
		Rating:        p.Rating,
		CategoryId:    p.CategoryID,
		Specification: string(specBytes),
		UpdatedAt:     updatedAt,
		CreatedAt:     p.CreatedAt.UnixMilli(),
	}
}

func NewProductFromPB(productPB *pbProducts.CreateProductRequest) (*Product, error) {
	spec, err := parseSpecificationFromPB(productPB.Specification)
	if err != nil {
		return nil, err
	}

	return &Product{
		ID:            uuid.New().String(),
		Name:          productPB.GetName(),
		Description:   productPB.GetDescription(),
		ImageID:       productPB.ImageId,
		Price:         productPB.GetPrice(),
		CurrencyID:    productPB.GetCurrencyId(),
		Rating:        productPB.GetRating(),
		CategoryID:    productPB.GetCategoryId(),
		Specification: spec,
		CreatedAt:     time.Now(),
	}, nil
}

func parseSpecificationFromPB(specFromPB string) (spec map[string]interface{}, unmarshalErr error) {
	if specFromPB != "" {
		return nil, errors.New("specification is empty")
	}

	if unmarshalErr = json.Unmarshal([]byte(specFromPB), &spec); unmarshalErr == nil {
		return spec, nil
	}

	return nil, errors.Wrap(unmarshalErr, "failed to parse specification")
}

func NewProductFromDAO(sp *dao.Product) *Product {
	var imageID *string
	if sp.ImageID.Valid {
		imageID = &sp.ImageID.String
	}

	createdAt, err := time.Parse(time.RFC3339, sp.CreatedAt.String)
	if err != nil {
		logging.GetLogger().WithError(err).Error("time.Parse(sp.CreatedAt)")
	}

	var updatedAt time.Time
	if sp.UpdatedAt.Valid {
		updatedAt, err = time.Parse(time.RFC3339, sp.UpdatedAt.String)
		if err != nil {
			logging.GetLogger().WithError(err).Error("time.Parse(sp.CreatedAt)")
		}
	}

	return &Product{
		ID:            sp.ID,
		Name:          sp.Name,
		Description:   sp.Description,
		ImageID:       imageID,
		Price:         sp.Price,
		CurrencyID:    sp.CurrencyID,
		Rating:        sp.Rating,
		CategoryID:    sp.CategoryID,
		Specification: sp.Specification,
		CreatedAt:     createdAt,
		UpdatedAt:     &updatedAt,
	}
}
