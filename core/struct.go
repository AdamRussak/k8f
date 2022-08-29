package core

import log "github.com/sirupsen/logrus"

type PlainFormatter struct {
}
type myFormatter struct {
	log.TextFormatter
}
