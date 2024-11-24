# put all application flags after `go run <path>`
npx nodemon --exec go run ./cmd/web -env="development" --signal SIGTERM -e go
