mockgen:
	mockgen -source=./internal/post/repository.go -destination=./internal/post/mocks/rep_mock.go
	mockgen -source=./internal/user/repository.go -destination=./internal/user/mocks/rep_mock.go
	mockgen -source=./internal/session/repository.go -destination=./internal/session/mocks/rep_mock.go

	mockgen -source=./internal/post/usecase.go -destination=./internal/post/mocks/use_mock.go
	mockgen -source=./internal/user/usecase.go -destination=./internal/user/mocks/use_mock.go
	mockgen -source=./internal/session/usecase.go -destination=./internal/session/mocks/use_mock.go

test:
	go test ./...