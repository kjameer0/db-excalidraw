# put all application flags after `go run <path>`
# TODO: add flags for development
# Flags:
# -env: "development or prod"
# -drawing-dir: path to save drawings to
# -addr: port to run server on, e.g -addr=":4000"
#
#
npx nodemon --exec go run ./cmd/web -env="development" --signal SIGTERM -e go
