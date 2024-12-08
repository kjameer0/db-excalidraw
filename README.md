# db-excalidraw

This app is a server that stores [Excalidraw](https://excalidraw.com/) drawings for others to use. It's different from the Excalidraw library because it's not for sharing components for others to use in their drawings; it's meant for users to be able to share and create different versions of the drawings they create.

## Project tools

If you are using VSCode to edit, you can download the [Excalidraw extension](https://marketplace.visualstudio.com/items?itemName=pomdtr.excalidraw-editor) to view the system design for the project, or copy the raw JSON into the browser.

## Stack

1. Postgres
2. Golang
3. NextJS
4. Amazon S3
5. Excalidraw

## Setting up Environment

You can use `nodemon` from `npm` to reload the server on file changes:
`npm install -g nodemon`
`npx nodemon --exec go run ./cmd/web --signal SIGTERM -e go`

`test.sh` contains the above `npx` command you can run to start the server with `source test.sh` if your terminal is in the root directory of the project.

### Environment variables

### Dev

### Integration testing

### Production

## Setting Up Postgres

To set up the database you might need to download psql.
`psql -f internal/db/schema.sql`

You might have to play around with your user permissions to get this file to run.

You can seed the database with this command
`psql -d excalidb -f internal/database/seed.sql`

Full command:

`psql -f internal/database/schema.sql;psql -d excalidb -f internal/database/seed.sql`

## API usage

### Excalidraw drawing data model

```json
{
  "drawing": {},
  "name": "drawing-name"
}
```



