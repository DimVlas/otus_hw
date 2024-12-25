package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"

	easyjson "github.com/mailru/easyjson"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return domainsStat(r, domain)
}

func domainsStat(r io.Reader, domain string) (DomainStat, error) {
	rd := bufio.NewReader(r)
	var br bool
	var user User
	stat := make(DomainStat)

	for {
		line, err := rd.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				br = true
			} else {
				return nil, err
			}
		}

		if err = easyjson.Unmarshal(line, &user); err != nil {
			return nil, err
		}

		if strings.HasSuffix(user.Email, "."+domain) {
			d := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			stat[d]++
		}

		if br {
			break
		}
	}

	return stat, nil
}
