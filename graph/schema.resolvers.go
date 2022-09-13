package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"

	"github.com/oscaralmgren/hackernews/graph/generated"
	"github.com/oscaralmgren/hackernews/graph/model"
	"github.com/oscaralmgren/hackernews/internal/auth"
	"github.com/oscaralmgren/hackernews/internal/links"
	"github.com/oscaralmgren/hackernews/internal/pkg/jwt"
	"github.com/oscaralmgren/hackernews/internal/users"
)

// CreateLink is the resolver for the createLink field.
func (r *mutationResolver) CreateLink(ctx context.Context, input model.NewLink) (*model.Link, error) {
	var link links.Link
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.Link{}, fmt.Errorf("access denied")
	}
	link.User = user
	link.Title = input.Title
	link.Address = input.Address
	linkID := link.Save()
	graphqlUser := &model.User{
		ID:   user.ID,
		Name: user.Username,
	}
	return &model.Link{
		ID:      strconv.FormatInt(linkID, 10),
		Title:   link.Title,
		Address: link.Address,
		User:    graphqlUser,
	}, nil
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	user.Create()
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	correct := user.Authenticate()
	if !correct {
		// 1
		return "", &users.WrongUsernameOrPasswordError{}
	}
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

// RefreshToken is the resolver for the refreshToken field.
func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	username, err := jwt.ParseToken(input.Token)
	if err != nil {
		return "", fmt.Errorf("access denied")
	}
	token, err := jwt.GenerateToken(username)
	if err != nil {
		return "", err
	}
	return token, nil
}

// Links is the resolver for the links field.
func (r *queryResolver) Links(ctx context.Context) ([]*model.Link, error) {
	var resultLinks []*model.Link
	var dbLinks []links.Link = links.GetAll()
	for _, link := range dbLinks {
		graphqlUser := &model.User{
			ID:   link.User.ID,
			Name: link.User.Username,
		}
		resultLinks = append(resultLinks, &model.Link{
			ID:      link.ID,
			Title:   link.Title,
			Address: link.Address,
			User:    graphqlUser})
	}
	return resultLinks, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
