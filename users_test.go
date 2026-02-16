package sdk

import (
	"context"
	"strconv"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

func TestGetSelfDetails_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/users/details").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		Reply(200).
		JSON(`{"id":1,"accessHash":"hash","nickname":"atennop","description":"something","interests":["vim"],"avatar":{"id":2,"accessHash":"hash2"}}`)

	client := NewClient()
	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	self, err := client.GetSelfDetails(context.Background(), auth)

	require.NoError(t, err)
	require.Equal(t, &UserDetails{
		Id:          MockUserId(1),
		AccessHash:  MockUserAccessHash("hash"),
		Nickname:    MockNickname("atennop"),
		Description: MockUserDescription("something"),
		Interests: MockInterests([]Interest{
			MockInterest("vim"),
		}),
		Avatar: &FileDescriptor{
			Id:         MockFileId(2),
			AccessHash: MockFileAccessHash("hash2"),
		},
	}, self)
}

func TestGetSelfDetails_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/users/details").
		Reply(400)

	client := NewClient()
	_, err := client.GetSelfDetails(context.Background(), nil)
	require.Error(t, err)
}

func TestGetUserDetails_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/users/details/2/hash2").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		Reply(200).
		JSON(`{"id":2,"accessHash":"hash2","nickname":"tr3ble","description":"something2","interests":["mac"],"avatar":{"id":3,"accessHash":"hash3"}}`)

	client := NewClient()
	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	user, err := client.GetUserDetails(context.Background(), auth, MockUserId(2), MockUserAccessHash("hash2"))

	require.NoError(t, err)
	require.Equal(t, &UserDetails{
		Id:          MockUserId(2),
		AccessHash:  MockUserAccessHash("hash2"),
		Nickname:    MockNickname("tr3ble"),
		Description: MockUserDescription("something2"),
		Interests: MockInterests([]Interest{
			MockInterest("mac"),
		}),
		Avatar: &FileDescriptor{
			Id:         MockFileId(3),
			AccessHash: MockFileAccessHash("hash3"),
		},
	}, user)
}

func TestGetUserDetails_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/users/details/2/hash2").
		Reply(400)

	client := NewClient()
	_, err := client.GetUserDetails(context.Background(), nil, MockUserId(2), MockUserAccessHash("hash2"))
	require.Error(t, err)
}

func TestOptions(t *testing.T) {
	cases := []struct {
		name           string
		options        []editAccountOption
		expectedBodies []string
	}{
		{
			name: "NicknameOption",
			options: []editAccountOption{
				EditNicknameOption(MockNickname("atennop")),
				EditNicknameOption(MockNickname("")),
			},
			expectedBodies: []string{
				`{"nickname":{"value":"atennop"}}`,
				`{"nickname":{"value":""}}`,
			},
		},
		{
			name: "DescriptionOption",
			options: []editAccountOption{
				EditDescriptionOption(MockUserDescription("bio")),
				EditDescriptionOption(MockUserDescription("")),
			},
			expectedBodies: []string{
				`{"description":{"value":"bio"}}`,
				`{"description":{"value":""}}`,
			},
		},
		{
			name: "InterestsOption",
			options: []editAccountOption{
				EditInterestsOption(MockInterests([]Interest{MockInterest("neovim"), MockInterest("coding")})),
				EditInterestsOption(MockInterests([]Interest{})),
			},
			expectedBodies: []string{
				`{"interests":{"value":["neovim", "coding"]}}`,
				`{"interests":{"value":[]}}`,
			},
		},
		{
			name: "AvatarOption",
			options: []editAccountOption{
				EditAvatarOption(&FileDescriptor{Id: MockFileId(10), AccessHash: MockFileAccessHash("hash")}),
				EditAvatarOption(nil),
			},
			expectedBodies: []string{
				`{"avatar":{"value":{"id":10, "accessHash":"hash"}}}`,
				`{"avatar":{"value":null}}`,
			},
		},
		{
			name: "SocialLinkOption",
			options: []editAccountOption{
				EditSocialLinkOption(MockSocialLink("https://example.com")),
				EditSocialLinkOption(MockSocialLink("")),
			},
			expectedBodies: []string{
				`{"socialLink":{"value":"https://example.com"}}`,
				`{"socialLink":{"value":""}}`,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for i := 0; i < len(tc.options); i++ {
				t.Run(tc.name+"_"+strconv.Itoa(i), func(t *testing.T) {
					defer gock.Off()

					gock.New("https://api.getfriend.ly").
						JSON(tc.expectedBodies[i]).
						Reply(200)

					c := NewClient()
					auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
					err := c.EditAccount(context.Background(), auth, tc.options[i])
					require.NoError(t, err)
				})
			}
		})
	}
}

func TestEditAccount_Empty(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Patch("/users/edit").
		Reply(200)

	c := NewClient()
	err := c.EditAccount(context.Background(), nil)

	require.NoError(t, err)
	require.False(t, gock.IsDone())
}

func TestEditAccount_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Patch("/users/edit").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		JSON(`{"nickname":{"value":"atennop"},"description":{"value":"bio"}}`).
		Reply(200)

	c := NewClient()
	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	err := c.EditAccount(context.Background(), auth, EditNicknameOption(MockNickname("atennop")), EditDescriptionOption(MockUserDescription("bio")))
	require.NoError(t, err)
}

func TestEditAccount_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Patch("/users/edit").
		Reply(400)

	c := NewClient()
	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	err := c.EditAccount(context.Background(), auth, EditNicknameOption(MockNickname("atennop")), EditDescriptionOption(MockUserDescription("bio")))
	require.Error(t, err)
}
