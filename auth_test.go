package sdk

import (
	"context"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

func TestRegister_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/auth/generate").
		JSON(`{"nickname":"atennop", "description":"bio","interests":["programming"],"avatar":{"id":10,"accessHash":"hash"},"socialLink":"https://github.com/Atennop1"}`).
		Reply(200).
		JSON(`{"id":1,"token":"token","accessHash":"hash"}`)

	client := NewClient()
	auth, err := client.Register(context.Background(),
		MockNickname("atennop"),
		MockUserDescription("bio"),
		MockInterests([]Interest{MockInterest("programming")}),
		&FileDescriptor{Id: MockFileId(10), AccessHash: MockFileAccessHash("hash")},
		MockSocialLink("https://github.com/Atennop1"))

	require.NoError(t, err)
	require.Equal(t, &Authorization{
		Id:         MockUserId(1),
		Token:      MockToken("token"),
		AccessHash: MockUserAccessHash("hash"),
	}, auth)
}

func TestRegister_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/auth/generate").
		Reply(400)

	client := NewClient()
	_, err := client.Register(context.Background(),
		MockNickname("atennop"),
		MockUserDescription("bio"),
		MockInterests([]Interest{MockInterest("programming")}),
		&FileDescriptor{Id: MockFileId(10), AccessHash: MockFileAccessHash("hash")},
		MockSocialLink("https://github.com/Atennop1"))

	require.Error(t, err)
}

func TestSendLoginRequest_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/auth/email").
		JSON(`{"email": "example@example.com"}`).
		Reply(200)

	client := NewClient()
	err := client.SendLoginRequest(context.Background(), MockEmail("example@example.com"))
	require.NoError(t, err)
}

func TestSendLoginRequest_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/auth/email").
		Reply(400)

	client := NewClient()
	err := client.SendLoginRequest(context.Background(), MockEmail("example@example.com"))
	require.Error(t, err)
}

func TestConfirmLogin_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/auth/login").
		JSON(`{"email":"example@example.com", "code":11111111}`).
		Reply(200).
		JSON(`{"id":1,"token":"token","accessHash":"hash"}`)

	client := NewClient()
	auth, err := client.ConfirmLogin(context.Background(), MockEmail("example@example.com"), MockEmailCode(11111111))

	require.NoError(t, err)
	require.Equal(t, &Authorization{
		Id:         MockUserId(1),
		Token:      MockToken("token"),
		AccessHash: MockUserAccessHash("hash"),
	}, auth)
}

func TestConfirmLogin_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/auth/login").
		Reply(400)

	client := NewClient()
	_, err := client.ConfirmLogin(context.Background(), MockEmail("example@example.com"), MockEmailCode(11111111))
	require.Error(t, err)
}
