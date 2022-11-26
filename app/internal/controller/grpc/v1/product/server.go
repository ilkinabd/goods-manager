package product

import (
	"context"
	pbProducts "github.com/ilkinabd/goods-contracts/gen/go/products/v1"
	"github.com/ilkinabd/goods-manager/app/internal/domain/product/filter"
	"github.com/ilkinabd/goods-manager/app/internal/domain/product/model"
	"github.com/ilkinabd/goods-manager/app/internal/domain/product/policy"
	"github.com/ilkinabd/goods-manager/app/pkg/logging"
)

type Server struct {
	policy *policy.ProductPolicy
	pbProducts.UnimplementedProductServiceServer
}

func NewServer(
	policy *policy.ProductPolicy,
	srv pbProducts.UnimplementedProductServiceServer,
) *Server {
	return &Server{
		policy:                            policy,
		UnimplementedProductServiceServer: srv,
	}
}

func (s *Server) AllProducts(
	ctx context.Context,
	request *pbProducts.AllProductsRequest,
) (*pbProducts.AllProductsResponse, error) {

	sort := filter.NewSortFromPB(request)
	criteria := make([]filter.Criteria, 1)

	categoryCriteria := filter.NewCategoryCriteriaFromPB(request)
	criteria = append(criteria, categoryCriteria)

	all, err := s.policy.All(ctx, criteria, sort)
	if err != nil {
		return nil, err
	}

	productsProto := make([]*pbProducts.Product, len(all))
	for i, p := range all {
		productsProto[i] = p.ToProto()
	}

	return &pbProducts.AllProductsResponse{
		Products: productsProto,
	}, nil
}

func (s *Server) ProductByID(
	ctx context.Context,
	req *pbProducts.ProductByIDRequest,
) (*pbProducts.ProductByIDResponse, error) {
	one, err := s.policy.One(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pbProducts.ProductByIDResponse{
		Product: one.ToProto(),
	}, nil
}

func (s *Server) UpdateProduct(
	ctx context.Context,
	req *pbProducts.UpdateProductRequest,
) (*pbProducts.UpdateProductResponse, error) {
	product, err := s.policy.One(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	product.UpdateFromPB(req)

	err = s.policy.Update(ctx, product)
	if err != nil {
		return nil, err
	}

	return &pbProducts.UpdateProductResponse{}, nil
}

func (s *Server) DeleteProduct(
	ctx context.Context,
	req *pbProducts.DeleteProductRequest,
) (*pbProducts.DeleteProductResponse, error) {
	err := s.policy.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pbProducts.DeleteProductResponse{}, nil
}

func (s *Server) CreateProduct(
	ctx context.Context,
	req *pbProducts.CreateProductRequest,
) (*pbProducts.CreateProductResponse, error) {
	p, err := model.NewProductFromPB(req)
	if err != nil {
		logging.WithError(ctx, err).WithField("product in pb", req).Error("model.NewProductFromPB")
		return nil, err
	}

	product, err := s.policy.CreateProduct(ctx, p)
	if err != nil {
		return nil, err
	}

	return &pbProducts.CreateProductResponse{
		Product: product.ToProto(),
	}, nil
}
