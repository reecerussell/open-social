//go:generate mockgen -package=mock -source=../users/client.go -destination=users/client.go
//go:generate mockgen -package=mock -source=../auth/client.go -destination=auth/client.go
//go:generate mockgen -package=mock -source=../media/client.go -destination=media/client.go

package mock
