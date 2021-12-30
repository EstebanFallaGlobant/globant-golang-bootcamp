package server

import (
	"fmt"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type CtmStringService struct {
	logger log.Logger
}

func NewService(logger log.Logger) *CtmStringService {
	result := new(CtmStringService)
	result.logger = logger

	return result
}

func (svc *CtmStringService) IsPal(s string) string {
	level.Info(svc.logger).Log(fmt.Sprintf("string: %s", s))
	reverse := svc.Reverse(s)

	if strings.ToLower(s) != reverse {
		return "not a palindrome"
	}

	return "is a palindrome"
}

func (svc *CtmStringService) Reverse(s string) string {
	level.Info(svc.logger).Log(fmt.Sprintf("string: %s", s))
	rns := []rune(s)

	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		rns[i], rns[j] = rns[j], rns[i]
	}

	return strings.ToLower(string(rns))
}
