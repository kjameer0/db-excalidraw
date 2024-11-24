# put all application flags after `go run <path>`
# TODO: add flags for development
npx nodemon --exec go run ./cmd/web -env="development" --signal SIGTERM -e go
