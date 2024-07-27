package ws

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Pratham-Mishra04/interactWS/initializers"
	"github.com/Pratham-Mishra04/interactWS/utils"
)

func ChatSetupHandler(event Event, c *Client) error {
	var setupEvent ChatSetupEvent

	if err := json.Unmarshal(event.Payload, &setupEvent); err != nil {
		return fmt.Errorf("bad payload in chat setup handler :%v", err)
	}

	c.chats = setupEvent.Chats
	if initializers.CONFIG.ENV == initializers.DevelopmentEnv {
		fmt.Println("Chats setup for user: " + c.userID)
	}
	return nil
}

func SendMessageHandler(event Event, c *Client) error {
	var chatEvent SendMessageEvent

	if err := json.Unmarshal(event.Payload, &chatEvent); err != nil {
		return fmt.Errorf("bad payload in send message handler :%v", err)
	}

	var broadMessage NewMessageEvent

	broadMessage.ID = chatEvent.ID
	broadMessage.CreatedAt = time.Now()
	broadMessage.Content = chatEvent.Content
	broadMessage.UserID = chatEvent.UserID
	broadMessage.User = chatEvent.User
	broadMessage.ChatID = chatEvent.ChatID
	broadMessage.Read = false

	data, err := json.Marshal(broadMessage)

	if err != nil {
		return fmt.Errorf("failed to marshall broadcast message")
	}

	outgoingEvent := Event{
		Type:    EventNewMessage,
		Payload: data,
	}

	for client := range c.manager.clients {
		if utils.Includes(client.chats, broadMessage.ChatID) {
			client.egress <- outgoingEvent
		}
	}

	return nil
}

func ReadMessageHandler(event Event, c *Client) error {
	var readEvent ReadEvent

	if err := json.Unmarshal(event.Payload, &readEvent); err != nil {
		return fmt.Errorf("bad payload in send message handler :%v", err)
	}

	var broadMessage UpdateReadEvent

	broadMessage.MessageID = readEvent.MessageID
	broadMessage.User = readEvent.User
	broadMessage.ChatID = readEvent.ChatID

	data, err := json.Marshal(broadMessage)

	if err != nil {
		return fmt.Errorf("failed to marshall broadcast message")
	}

	outgoingEvent := Event{
		Type:    UpdateRead,
		Payload: data,
	}

	for client := range c.manager.clients {
		if utils.Includes(client.chats, broadMessage.ChatID) {
			client.egress <- outgoingEvent
		}
	}

	return nil
}

func NotificationHandler(event Event, c *Client) error {
	var notificationEvent NotificationEvent

	if err := json.Unmarshal(event.Payload, &notificationEvent); err != nil {
		return fmt.Errorf("bad payload in send notification handler :%v", err)
	}

	var sendNotification NotificationEvent

	sendNotification.UserID = notificationEvent.UserID
	sendNotification.Content = notificationEvent.Content

	data, err := json.Marshal(sendNotification)

	if err != nil {
		return fmt.Errorf("failed to marshall send notification")
	}

	outgoingEvent := Event{
		Type:    EventReceiveNotification,
		Payload: data,
	}

	for client := range c.manager.clients {
		if client.userID == sendNotification.UserID {
			client.egress <- outgoingEvent
		}
	}

	return nil
}

// func ChatRoomHandler(event Event, c *Client) error {
// 	var changeChatEvent ChangeChatEvent

// 	if err := json.Unmarshal(event.Payload, &changeChatEvent); err != nil {
// 		return fmt.Errorf("bad payload in chat room handle :%v", err)
// 	}

// 	c.chatID = changeChatEvent.ID

// 	return nil
// }

func MeTypingHandler(event Event, c *Client) error {
	var meTypingEvent MeTypingEvent

	if err := json.Unmarshal(event.Payload, &meTypingEvent); err != nil {
		return fmt.Errorf("bad payload in me typing :%v", err)
	}

	var userTyping UserTypingEvent

	userTyping.User = meTypingEvent.User
	userTyping.ChatID = meTypingEvent.ChatID

	data, err := json.Marshal(userTyping)

	if err != nil {
		return fmt.Errorf("failed to marshall user typing message")
	}

	outgoingEvent := Event{
		Type:    UserTyping,
		Payload: data,
	}

	for client := range c.manager.clients {
		if utils.Includes(client.chats, userTyping.ChatID) {
			client.egress <- outgoingEvent
		}
	}

	return nil
}

func MeStopTypingHandler(event Event, c *Client) error {
	var meStopTypingEvent MeStopTypingEvent

	if err := json.Unmarshal(event.Payload, &meStopTypingEvent); err != nil {
		return fmt.Errorf("bad payload in me stop typing :%v", err)
	}

	var userStopTyping UserStopTypingEvent

	userStopTyping.User = meStopTypingEvent.User
	userStopTyping.ChatID = meStopTypingEvent.ChatID

	data, err := json.Marshal(userStopTyping)

	if err != nil {
		return fmt.Errorf("failed to marshall user stop typing message")
	}

	outgoingEvent := Event{
		Type:    UserStopTyping,
		Payload: data,
	}

	for client := range c.manager.clients {
		if utils.Includes(client.chats, userStopTyping.ChatID) {
			client.egress <- outgoingEvent
		}
	}

	return nil
}

func UpdateMembershipHandler(event Event, c *Client) error {
	var updateMembershipEvent UpdateMembership

	if err := json.Unmarshal(event.Payload, &updateMembershipEvent); err != nil {
		return fmt.Errorf("bad payload in update membership :%v", err)
	}

	var updateMembership UpdateMembership

	updateMembership.UserID = updateMembershipEvent.UserID
	updateMembership.OrganizationID = updateMembershipEvent.OrganizationID
	updateMembership.Role = updateMembershipEvent.Role

	data, err := json.Marshal(updateMembership)

	if err != nil {
		return fmt.Errorf("failed to marshall update membership")
	}

	outgoingEvent := Event{
		Type:    ReceiveUpdateMembershipEvent,
		Payload: data,
	}

	for client := range c.manager.clients {
		if client.userID == updateMembership.UserID {
			client.egress <- outgoingEvent
		}
	}

	return nil
}
