package ws

import (
	"encoding/json"
	"fmt"
	"time"
)

func SendMessageHandler(event Event, c *Client) error {
	var chatEvent SendMessageEvent

	if err := json.Unmarshal(event.Payload, &chatEvent); err != nil {
		return fmt.Errorf("Bad Payload :%v", err)
	}

	var broadMessage NewMessageEvent

	broadMessage.CreatedAt = time.Now()
	broadMessage.Content = chatEvent.Content
	broadMessage.UserID = chatEvent.UserID
	broadMessage.User = chatEvent.User
	broadMessage.ChatID = chatEvent.ChatID
	broadMessage.Read = false

	data, err := json.Marshal(broadMessage)

	if err != nil {
		return fmt.Errorf("Failed to Marshall BroadCast Message")
	}

	outgoingEvent := Event{
		Type:    EventNewMessage,
		Payload: data,
	}

	for client := range c.manager.clients {
		if client.chatID == broadMessage.ChatID {
			client.egress <- outgoingEvent
		}
	}

	return nil

}

func ChatRoomHandler(event Event, c *Client) error {
	var changeChatEvent ChangeChatEvent

	if err := json.Unmarshal(event.Payload, &changeChatEvent); err != nil {
		return fmt.Errorf("Bad Payload :%v", err)
	}

	c.chatID = changeChatEvent.ID

	return nil
}

func MeTypingHandler(event Event, c *Client) error {
	var meTypingEvent MeTypingEvent

	if err := json.Unmarshal(event.Payload, &meTypingEvent); err != nil {
		return fmt.Errorf("Bad Payload :%v", err)
	}

	var userTyping UserTypingEvent

	userTyping.User = meTypingEvent.User
	userTyping.ChatID = meTypingEvent.ChatID

	data, err := json.Marshal(userTyping)

	if err != nil {
		return fmt.Errorf("Failed to Marshall User Typing Message")
	}

	outgoingEvent := Event{
		Type:    UserTyping,
		Payload: data,
	}

	for client := range c.manager.clients {
		if client.chatID == userTyping.ChatID {
			client.egress <- outgoingEvent
		}
	}

	return nil
}

func MeStopTypingHandler(event Event, c *Client) error {
	var meStopTypingEvent MeStopTypingEvent

	if err := json.Unmarshal(event.Payload, &meStopTypingEvent); err != nil {
		return fmt.Errorf("Bad Payload :%v", err)
	}

	var userStopTyping UserStopTypingEvent

	userStopTyping.User = meStopTypingEvent.User
	userStopTyping.ChatID = meStopTypingEvent.ChatID

	data, err := json.Marshal(userStopTyping)

	if err != nil {
		return fmt.Errorf("Failed to Marshall User Stop Typing Message")
	}

	outgoingEvent := Event{
		Type:    UserStopTyping,
		Payload: data,
	}

	for client := range c.manager.clients {
		if client.chatID == userStopTyping.ChatID {
			client.egress <- outgoingEvent
		}
	}

	return nil
}
