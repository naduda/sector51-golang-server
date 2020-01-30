migrate create -ext sql -dir migrations create_user_security

migrate -path migrations -database "postgres://postgres:12345678@localhost/sector51_test?sslmode=disable" up
migrate -path migrations -database "postgres://postgres:12345678@localhost/sector51?sslmode=disable" up