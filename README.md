# db-excalidraw


## Copying Drawing

If you have the `jq` command installed on your Unix system, you can run something like:
`curl localhost:4000/drawing/hih | jq '.["drawingJson"]' | pbcopy`
when you have a dev server live to copy a drawing directly to your clipboard.
