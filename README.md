```bash
go run . ~/the/folder/you/want/to/monitor
```

## WIP! Check commit logs 

## Working: 

- rendering ms after file change
- support multiple files in same folder
- background file change monitoring 
- rendering using the gin web server library and web sockets
- mermaid diagrams 
- code highlight
- GitHub theme

## Not working right now: 

- recursive folder monitoring
  - [ ] implement my own and add a file watcher for every `*.md` file
- scroll to last edited position
  - I don't know how to solve this, since I have no way to get the current cursor position
