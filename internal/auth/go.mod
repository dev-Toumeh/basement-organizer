module basement/auth

go 1.22.1

require (
	github.com/google/uuid v1.6.0
	internal/util v0.0.0-00010101000000-000000000000
)

require golang.org/x/crypto v0.23.0 // indirect

replace internal/util => ../util
