package models

var voteOptions = []string{
	"1",
	"2",
	"3",
	"5",
	"8",
	"13",
	"21",
	"34",
	"?",
}

func IsValidVote(vote string) bool {
	for _, option := range voteOptions {
		if option == vote {
			return true
		}
	}
	return false
}
