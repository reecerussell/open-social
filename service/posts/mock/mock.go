//go:generate mockgen -package=mock -source=../repository/post_repository.go -destination=repository/post_repository.go
//go:generate mockgen -package=mock -source=../repository/like_repository.go -destination=repository/like_repository.go
//go:generate mockgen -package=mock -source=../provider/post_provider.go -destination=provider/post_provider.go

package mock
