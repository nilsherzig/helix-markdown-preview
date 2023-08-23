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

- Custom kanban renderer / parser
  - run parser as marked extension in client side js 
  - using a code block (?)
  - Somehow add kanban auto formatting to helix (pipe to external script?)
- Automatically start when helix opens a markdown file
- recursive folder monitoring
  - implement my own and add a file watcher for every `*.md` file
- scroll to current position
  - I don't know how to get the current line number from helix, maybe scroll to last edit? 
