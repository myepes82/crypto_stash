package models

import "time"

type Secrets struct {
	Changelog time.Time         `json:"changelog"`
	Content   map[string]string `json:"secrets"`
}

func NewSecrets() *Secrets {
	return &Secrets{}
}

func (s *Secrets) UpdateChangelog() {
	s.Changelog = time.Now()
}

func (s *Secrets) AddSecret(key, value string) {
	if s.Content == nil {
		s.Content = make(map[string]string)
	}
	s.Content[key] = value
}

func (s *Secrets) GetSecret(key string) (string, bool) {
	value, ok := s.Content[key]
	return value, ok
}
