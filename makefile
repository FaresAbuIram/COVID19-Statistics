tests:
	go test ./...
mocks:
	@mockery@2.14.0 --name="SQLRepositoryInterface" --dir="./database" --output="./services/mocks"