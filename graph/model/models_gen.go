// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type AuthResponse struct {
	JwtToken   string    `json:"jwtToken"`
	ExpirateAt time.Time `json:"expirateAt"`
}

type EntitiesDraft struct {
	Tags  []string            `json:"tags"`
	Media []*MediaEntityDraft `json:"media"`
	Urls  []*URLEntityDraft   `json:"urls"`
}

type EventDraft struct {
	Title    string         `json:"title"`
	Message  string         `json:"message"`
	Actor    *PersonDraft   `json:"actor"`
	Kind     *string        `json:"kind"`
	ShareURL *string        `json:"shareURL"`
	Entities *EntitiesDraft `json:"entities"`
}

type EventStreamFilter struct {
	Sources []string `json:"sources"`
}

type MediaEntityDraft struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

type PersonDraft struct {
	Name     string  `json:"name"`
	Photo    *string `json:"photo"`
	Username *string `json:"username"`
}

type URLEntityDraft struct {
	URL        string `json:"url"`
	DisplayURL string `json:"displayURL"`
}
