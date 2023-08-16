## WIP! Check commit logs 

## Usage

```bash
go run . ~/the/folder/you/want/to/monitor

# or build it 

go build . 
./helix-markdown-preview ~/the/folder/you/want/to/monitor
```

## Working: 

- rendering ms after file change
- support multiple files in same folder
- background file change monitoring 
- rendering using the gin web server library and web sockets
- mermaid diagrams 
- code highlight
- GitHub markdown support 
- GitHub markdown theme
- templates and scripts are now embedded into the binary during the build process

## Things I would like to change / implement

- Automatically start when helix opens a markdown file
- recursive folder monitoring
  - implement my own and add a file watcher for every `*.md` file
- scroll to last edited position
  - I don't know how to implement this since I don't know the current cursor position in helix. Also, since more than one place in the document can be changed since the last save, I can't just scroll to the difference between the two versions
