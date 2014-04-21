package backend

import (
	"strings"
)

type stub struct {
}

func (s *stub) RevcCommit(commitMessage []byte) error {
	message := string(commitMessage)
	items := strings.Split(message, ",")
}
