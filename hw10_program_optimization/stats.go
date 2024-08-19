package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		var user User
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, err
		}

		if strings.HasSuffix(user.Email, "."+domain) {
			key := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[key]++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
