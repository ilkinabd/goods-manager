package filter

import (
	sq "github.com/Masterminds/squirrel"
	pbProduct "github.com/ilkinabd/goods-contracts/gen/go/products/v1"
)

const (
	fieldName = "category_id"
)

func NewCategoryCriteriaFromPB(product *pbProduct.AllProductsRequest) Criteria {
	return criteria{
		Name:  fieldName,
		Value: product.CategoryId.GetValue(),
	}
}

func (c criteria) MeetCriteria(query sq.SelectBuilder) sq.SelectBuilder {
	if c.Value != "" {
		and := sq.And{}
		e := sq.Eq{c.Name: c.Value}
		and = append(and, e)
		query = query.Where(and)
	}
	return query
}
