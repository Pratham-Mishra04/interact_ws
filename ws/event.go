package ws

import (
	"encoding/json"
	"time"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"` // No type, will accept anything send from the frontend
}

type EventHandler func(event Event, c *Client) error

const (
	EventSendMessage = "send_message"
	EventNewMessage  = "new_message"
	EventChangeChat  = "change_chat"
)

type UserType struct {
	ID                string    `json:"id"`
	Email             string    `json:"email"`
	Name              string    `json:"name"`
	ProfilePic        string    `json:"profilePic"`
	CoverPic          string    `json:"coverPic"`
	Username          string    `json:"username"`
	PhoneNo           string    `json:"phoneNo"`
	Bio               string    `json:"bio"`
	Title             string    `json:"title"`
	Tagline           string    `json:"tagline"`
	Education         []string  `json:"education"`
	Achievements      []string  `json:"achievements"`
	Followers         []string  `json:"followers"`
	Following         []string  `json:"following"`
	Posts             []string  `json:"posts"`
	Projects          []string  `json:"projects"`
	NoFollowers       int       `json:"noFollowers"`
	NoFollowing       int       `json:"noFollowing"`
	IsFollowing       bool      `json:"isFollowing"`
	PasswordChangedAt time.Time `json:"passwordChangedAt"`
	LastViewed        []string  `json:"lastViewed"`
	Tags              []string  `json:"tags"`
}

type SendMessageEvent struct {
	Content string   `json:"content"`
	ChatID  string   `json:"chatID"`
	User    UserType `json:"user"`
	UserID  string   `json:"userID"`
}

type NewMessageEvent struct {
	SendMessageEvent
	CreatedAt time.Time `json:"createdAt"`
	Read      bool      `json:"read"`
}

type ChangeChatEvent struct {
	ID string `json:"id"`
}
