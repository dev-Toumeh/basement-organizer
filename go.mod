module basement/main

go 1.22.1

require (
	internal/auth v0.0.0-00010101000000-000000000000
	internal/util v0.0.0-00010101000000-000000000000
)

require (
	github.com/google/uuid v1.6.0 // indirect
	golang.org/x/crypto v0.23.0 // indirect
)

replace internal/util => ./internal/util

replace internal/auth => ./internal/auth
