package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"project-z/pkg/graph/generated"
	"project-z/pkg/graph/model"
)

func (r *queryResolver) Products(ctx context.Context) ([]*model.Product, error) {

	dummyProduct := &model.Product{
		ID:            "123",
		Name:          "Dummy",
		Stock:         21,
		Description:   "Hello World",
		OriginalPrice: 100,
		ComparePrice:  100.32,
		Images:        []*model.Image{},
		Created:       122123123,
	}

	products := []*model.Product{}
	products = append(products, dummyProduct)

	return products, nil
	panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
