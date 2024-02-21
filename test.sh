go test ./... -coverprofile=profile
ROOT=. PROFILE=profile NAME=gobar go run .