package ws

import (
	"encoding/json"
	"time"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"` // Will accept anything send from the frontend
}

type EventHandler func(event Event, c *Client) error

const (
	ChatSetup                    = "chat_setup"
	EventSendMessage             = "send_message"
	EventNewMessage              = "new_message"
	EventSendNotification        = "send_notification"
	EventReceiveNotification     = "receive_notification"
	GetRead                      = "send_read_message"
	UpdateRead                   = "read_message"
	EventChangeChat              = "change_chat"
	MeTyping                     = "me_typing"
	MeStopTyping                 = "me_stop_typing"
	UserTyping                   = "user_typing"
	UserStopTyping               = "user_stop_typing"
	SendUpdateMembershipEvent    = "send_update_membership"
	ReceiveUpdateMembershipEvent = "receive_update_membership"
	HackathonSetup               = "hackathon_setup"
	HackathonUpdateSendEvent     = "send_hackathon_update"
	HackathonUpdateReceiveEvent  = "receive_hackathon_update"
	SendHackathonAnnouncement    = "send_new_hackathon_announcement"
	ReceiveHackathonAnnouncement = "receive_new_hackathon_announcement"
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

type HackathonType struct {
	ID                      string    `json:"id"`
	OrganizationID          string    `json:"organizationID"`
	MinTeamSize             int8      `json:"minTeamSize"`
	MaxTeamSize             int8      `json:"maxTeamSize"`
	TeamFormationStartTime  time.Time `json:"teamFormationStartTime"`
	TeamFormationEndTime    time.Time `json:"teamFormationEndTime"`
	StartTime               time.Time `json:"startTime"`
	EndTime                 time.Time `json:"endTime"`
	IsEnded                 bool      `json:"isEnded"`
	AllowEditDuringJudging  bool      `json:"allowEditDuringJudging"`
	EnableGithubIntegration bool      `json:"enableGithubIntegration"`
	EnableFigmaIntegration  bool      `json:"enableFigmaIntegration"`
	EnableAutoCodeReviews   bool      `json:"enableAutoCodeReviews"`
	MakeProjectsPublic      bool      `json:"makeProjectsPublic"`
}

type AnnouncementType struct {
	ID             string    `json:"id"`
	HackathonID    string    `json:"hackathonID"`
	OrganizationID string    `json:"organizationID"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"createdAt"`
}

type ChatSetupEvent struct {
	Chats []string `json:"chats"`
}

type SendMessageEvent struct {
	ID      string   `json:"id"`
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

type MeTypingEvent struct {
	User   UserType `json:"user"`
	ChatID string   `json:"chatID"`
}

type MeStopTypingEvent struct {
	MeTypingEvent
}

type UserTypingEvent struct {
	MeTypingEvent
}

type UserStopTypingEvent struct {
	MeTypingEvent
}

type NotificationEvent struct {
	UserID  string `json:"userID"`
	Content string `json:"content"`
}

type ReadEvent struct {
	User      UserType `json:"user"`
	MessageID string   `json:"messageID"`
	ChatID    string   `json:"chatID"`
}

type UpdateReadEvent struct {
	User      UserType `json:"user"`
	MessageID string   `json:"messageID"`
	ChatID    string   `json:"chatID"`
}

type UpdateMembership struct {
	UserID         string `json:"userID"`
	OrganizationID string `json:"organizationID"`
	Role           string `json:"role"`
}

type UpdateHackathonEvent struct {
	Hackathon HackathonType `json:"hackathon"`
}

type NewHackathonAnnouncementEvent struct {
	Announcement AnnouncementType `json:"announcement"`
}
