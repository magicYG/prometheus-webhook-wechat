package models

type Config struct {
	Targets          map[string]Target `yaml:"targets"`
	MaxContentLength int               `yaml:"maxContentLength"`
}

type Target struct {
	URL     string         `yaml:"url"`
	Mention *MentionTarget `yaml:"mention,omitempty"`
}

type MentionTarget struct {
	ALL     bool     `yaml:"all,omitempty"`
	Mobiles []string `yaml:"mobiles,omitempty"`
}
