package graph

import (
	"postsandcomments/internal/db"

	"github.com/sirupsen/logrus"
)

type Resolver struct {
	DataBase            db.Database
	SubscriptionManager *SubscriptionManager
	Logger              *logrus.Logger
}

type mutationResolver struct {
	*Resolver
}

type queryResolver struct {
	*Resolver
}

type subscriptionResolver struct {
	*Resolver
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}
