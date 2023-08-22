// This file contains types that are used in the repository layer.
package repository

type RegistrationInput struct {
	Slug     string
	FullName string
	Phone    string
	Password string
}

type RegistrationOutput struct {
	Id int
}

type FindByPhoneInput struct {
	Phone string
}

type FindByPhoneOutput struct {
	Id       int
	Slug     string
	FullName string
	Phone    string
	Password string
}

type FindBySlugInput struct {
	Slug string
}

type FindBySlugOutput struct {
	Slug     string
	FullName string
	Phone    string
	Password string
}

type UpdateUserInput struct {
	Slug     string
	FullName string
	Phone    string
}
