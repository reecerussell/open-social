//go:generate mockgen -package=mock -source=../password/validator.go -destination=validator.go
//go:generate mockgen -package=repository -source=../repository/user_repository.go -destination=repository/user_repository.go

package mock
