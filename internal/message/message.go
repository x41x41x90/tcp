package message

import (
	"regexp"
	"strings"
)

var (
	expTestExp  = `^test$`
	expShellExp = `^shell `
	updateExp   = `^update `
	fileExp     = `^file `

	TEST           = `test`
	SHELL          = `shell`
	UPDATE         = `update`
	FILE           = `file`
	unknownMessage = `unknown`
)

// Message for work with message.
type Message struct {
	testConnection *regexp.Regexp
	shellCommand   *regexp.Regexp
	updateCommand  *regexp.Regexp
	fileCommand    *regexp.Regexp
}

func New() *Message {
	return &Message{
		testConnection: regexp.MustCompile(expTestExp),
		shellCommand:   regexp.MustCompile(expShellExp),
		updateCommand:  regexp.MustCompile(updateExp),
		fileCommand:    regexp.MustCompile(fileExp),
	}
}

func (m *Message) Processing(message string) (t string, args []string) {
	t = unknownMessage
	if m.testConnection.MatchString(message) {
		t = TEST
	}
	if m.shellCommand.MatchString(message) {
		t = SHELL
	}
	if m.fileCommand.MatchString(message) {
		t = FILE
	}
	if m.updateCommand.MatchString(message) {
		t = UPDATE
	}
	commands := strings.Split(message, " ")
	if len(commands) > 1 {
		args = commands[1:]
	}
	return
}
