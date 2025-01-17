package graph

import (
	"context"

	"github.com/biFebriansyah/goraphql/graph/model"
)

func (r *mutationResolver) CreateProduct(ctx context.Context, input model.NewProduct) (*model.Products, error) {
	data, err := r.ProductService.CreateOne(input)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *queryResolver) Product(ctx context.Context, id string) (*model.Products, error) {
	data, err := r.ProductService.GetById(id)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *queryResolver) Products(ctx context.Context) ([]*model.Products, error) {
	data, err := r.ProductService.GetAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}
