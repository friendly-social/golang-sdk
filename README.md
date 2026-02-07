# Friendly Go SDK
This is the official Go SDK for interacting with Friendly API.

## Installation 
```bash
go get github.com/friendly-social/golang-sdk
```

## Quick Start
All API interactions are done through `sdk.Client` struct. For example:
```go
import (
	"context"
	"fmt"
	"log"

	sdk "github.com/friendly-social/golang-sdk"
)

func main() {
	client := sdk.NewClient("https://api.getfriend.ly")
	ctx := context.Background()

	auth, err := client.Register(ctx,
		"atennop",
		"the author of this SDK",
		Interests{"programming", "learning", "neovim", "debian"},
		&FileDescriptor{Id: 1, AccessHash: "very-long-hash"},
		"https://github.com/Atennop1")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Leaked Data:\nID: %d\nToken: %s\nAccessHash: %s\n", auth.Id, auth.Token, auth.AccessHash)
}
```

## Features
- Lightweight wrapper around plain net/http
- Streaming file uploads/downloads
- Context-aware requests
- Typed sentinel errors
- 100% test coverage

## Contributing
PRs are welcome! Please ensure that public APIs remain backward compatible and that all tests pass, and your changes will be merged :)
