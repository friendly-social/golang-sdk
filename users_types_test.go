package sdk

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func MockNickname(s string) Nickname {
	return Nickname{value: s}
}

func MockUserDescription(s string) UserDescription {
	return UserDescription{value: s}
}

func MockInterest(s string) Interest {
	return Interest{value: s}
}

func MockInterests(i []Interest) Interests {
	return Interests{value: i}
}

func MockSocialLink(s string) SocialLink {
	return SocialLink{value: s}
}

func TestUsersTypes(t *testing.T) {
	t.Run("Nickname", func(t *testing.T) {
		nickname, err := NewNickname("atennop")
		require.NoError(t, err)
		require.Equal(t, MockNickname("atennop"), nickname)
		require.Equal(t, "atennop", nickname.Value())

		_, err = NewNickname(strings.Repeat("1", 4096))
		require.Error(t, err)
		require.ErrorIs(t, err, ErrTooLongNickname)

		_, err = NewNickname("")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrEmptyNickname)

		data, err := json.Marshal(nickname)
		require.NoError(t, err)
		require.Equal(t, `"atennop"`, string(data))

		var loadedNickname Nickname
		err = json.Unmarshal(data, &loadedNickname)
		require.NoError(t, err)
		require.Equal(t, nickname, loadedNickname)
	})

	t.Run("Description", func(t *testing.T) {
		desc, err := NewUserDescription("something")
		require.NoError(t, err)
		require.Equal(t, MockUserDescription("something"), desc)
		require.Equal(t, "something", desc.Value())

		_, err = NewUserDescription(strings.Repeat("1", 4096))
		require.Error(t, err)
		require.ErrorIs(t, err, ErrTooLongUserDescription)

		_, err = NewUserDescription("")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrEmptyUserDescription)

		data, err := json.Marshal(desc)
		require.NoError(t, err)
		require.Equal(t, `"something"`, string(data))

		var loadedDesc UserDescription
		err = json.Unmarshal(data, &loadedDesc)
		require.NoError(t, err)
		require.Equal(t, desc, loadedDesc)
	})

	t.Run("Interest", func(t *testing.T) {
		interest, err := NewInterest("vim")
		require.NoError(t, err)
		require.Equal(t, MockInterest("vim"), interest)
		require.Equal(t, "vim", interest.Value())

		_, err = NewInterest(strings.Repeat("1", 4096))
		require.Error(t, err)
		require.ErrorIs(t, err, ErrTooLongInterest)

		_, err = NewInterest("")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrEmptyInterest)

		data, err := json.Marshal(interest)
		require.NoError(t, err)
		require.Equal(t, `"vim"`, string(data))

		var loadedInterest Interest
		err = json.Unmarshal(data, &loadedInterest)
		require.NoError(t, err)
		require.Equal(t, interest, loadedInterest)
	})

	t.Run("Interests", func(t *testing.T) {
		interests, err := NewInterests(MockInterest("vim"), MockInterest("debian"))
		require.NoError(t, err)
		require.Equal(t, MockInterests([]Interest{MockInterest("vim"), MockInterest("debian")}), interests)
		require.Equal(t, []Interest{MockInterest("vim"), MockInterest("debian")}, interests.Value())

		_, err = NewInterests(slices.Repeat([]Interest{MockInterest("vim")}, 1000)...)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrTooMuchInterests)

		_, err = NewInterests()
		require.Error(t, err)
		require.ErrorIs(t, err, ErrEmptyInterests)

		data, err := json.Marshal(interests)
		require.NoError(t, err)
		require.Equal(t, `["vim","debian"]`, string(data))

		var loadedInterests Interests
		err = json.Unmarshal(data, &loadedInterests)
		require.NoError(t, err)
		require.Equal(t, interests, loadedInterests)
	})

	t.Run("SocialLink", func(t *testing.T) {
		link, err := NewSocialLink("https://github.com/Atennop1")
		require.NoError(t, err)
		require.Equal(t, MockSocialLink("https://github.com/Atennop1"), link)
		require.Equal(t, "https://github.com/Atennop1", link.Value())

		_, err = NewSocialLink(strings.Repeat("1", 4096))
		require.Error(t, err)
		require.ErrorIs(t, err, ErrTooLongSocialLink)

		_, err = NewSocialLink("")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrEmptySocialLink)

		data, err := json.Marshal(link)
		require.NoError(t, err)
		require.Equal(t, `"https://github.com/Atennop1"`, string(data))

		var loadedLink SocialLink
		err = json.Unmarshal(data, &loadedLink)
		require.NoError(t, err)
		require.Equal(t, link, loadedLink)
	})
}
