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
	client := sdk.NewClient()
	ctx := context.Background()

	nickname, _ := NewNickname("atennop")
	description, _ := NewUserDescription("the author of this SDK")
	socialLink, _ := NewSocialLink("https://github.com/Atennop1")

	interest1, _ := NewInterest("programming")
	interest2, _ := NewInterest("learning")
	interests, _ := NewInterests(interest1, interest2)

	avatarId := NewFileId(1)
	avatarHash, _ := NewFileAccessHash("very-long-hash")
	avatar := &FileDescriptor{Id: avatarId, AccessHash: avatarHash}

	auth, err := client.Register(ctx, nickname, description, interests, avatar, socialLink)
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
