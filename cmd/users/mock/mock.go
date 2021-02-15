//go:generate mockgen -package=mock -source=../password/validator.go -destination=validator.go
//go:generate mockgen -package=repository -source=../repository/user_repository.go -destination=repository/user_repository.go
//go:generate mockgen -package=repository -source=../repository/follower_repository.go -destination=repository/follower_repository.go
//go:generate mockgen -package=mock -source=../provider/user_provider.go -destination=provider/user_provider.go

package mock
