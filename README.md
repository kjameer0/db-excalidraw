# db-excalidraw

This app is a server that stores [Excalidraw](https://excalidraw.com/) drawings for others to use. It's different from the Excalidraw library because it's not for sharing components for others to use in their drawings; it's meant for users to be able to share and create different versions of the drawings they create.

## Project tools

If you are using VSCode to edit, you can download the [Excalidraw extension](https://marketplace.visualstudio.com/items?itemName=pomdtr.excalidraw-editor) to view the system design for the project, or copy the raw JSON into the browser.

## Stack

## Copying Drawing

If you have the `jq` command installed on your Unix system, you can run something like:
`curl localhost:4000/drawing/hih | jq '.["drawingJson"]' | pbcopy`
when you have a dev server live to copy a drawing directly to your clipboard.

## API usage

### Excalidraw drawing data model

```json
{
  "drawing": {},
  "name": "drawing-name"
}
```
