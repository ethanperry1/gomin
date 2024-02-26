go test ./... -coverprofile=profile
ROOT=. PROFILE=profile NAME=gomin go run .