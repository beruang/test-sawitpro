// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	FindByPhone(ctx context.Context, input FindByPhoneInput) (FindByPhoneOutput, error)
	FindBySlug(ctx context.Context, input FindBySlugInput) (FindBySlugOutput, error)
	Store(ctx context.Context, input RegistrationInput) (RegistrationOutput, error)
	Put(ctx context.Context, input UpdateUserInput) error
}
