package filter

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	pbProduct "github.com/ilkinabd/goods-contracts/gen/go/products/v1"
)

type criteria struct {
	Name  string
	Value string
	Type  string
}

type sort struct {
	field string
	order string
}

type Criteria interface {
	MeetCriteria(query sq.SelectBuilder) sq.SelectBuilder
}

type Sortable interface {
	Sort(query sq.SelectBuilder) sq.SelectBuilder
}

func NewSortFromPB(product *pbProduct.AllProductsRequest) Sortable {
	return sort{
		field: product.GetSort().GetField(),
		order: product.GetSort().GetOrder(),
	}
}

func (s sort) Sort(query sq.SelectBuilder) sq.SelectBuilder {
	if s.field == "" || s.order == "" {
		return query
	}
	field := s.field
	return query.OrderBy(fmt.Sprintf("%s %s", field, s.order))
}
