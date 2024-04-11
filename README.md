# Svelte-Go

Example project setup with a `Go` server a `Svelte` UI and `GRPC` in between.

## Featurs

- [x] VSCode workspace setup
- 

## Global Installations
- Install protocol buffers compiler 
  - `brew install protobuf@3`
  - `echo 'export PATH="/opt/homebrew/opt/protobuf@3/bin:$PATH"' >> ~/.zshrc`
- Install GO plugins for protocol buffers 
  - `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`
  - `echo 'export PATH="$PATH:$HOME/.local/bin:$(go env GOPATH)/bin"' >> ~/.zshrc`

## Commands

- Generate GRPC code: `./proto.sh` (Execute every time you change the proto files)
- **[server]**:
 - Sync used modules from code with the `go.mod` declarations: `go mod tidy`
 - Start the server `go run main.go`
