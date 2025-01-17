package graph

import (
	"context"

	"github.com/biFebriansyah/goraphql/graph/model"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.Users, error) {
	data, err := r.UserService.CreateOne(input)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.Users, error) {
	data, err := r.UserService.GetById(id)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.Users, error) {
	data, err := r.UserService.GetAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}
