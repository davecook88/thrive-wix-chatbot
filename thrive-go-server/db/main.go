package db

import (
	"context"
	"thrive/server/chatgpt"
	"thrive/server/wix"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var COLLECTION_NAME = "thrive-chats"

type Client struct {
	*firestore.Client
}

type SavedChatRecord struct {
	Messages    []chatgpt.Message `json:"messages"`
	MemberId    string            `json:"memberId"`
	MemberName  string            `json:"memberName"`
	LastUpdated string            `json:"lastUpdated"`
}

func NewClient(ctx context.Context, projectID string) (*Client, error) {
	opt := option.WithCredentialsFile("creds.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	return &Client{Client: client}, nil

}

func (c *Client) CreateChat(ctx context.Context, messages []chatgpt.Message, wixUser wix.WixMember) error {
	chatDoc := map[string]interface{}{
		"messages":   messages,
		"memberId":   wixUser.ID,
		"memberName": wixUser.Profile.Nickname,
		// Add any other fields you want to store in the document
	}
	_, _, err := c.Collection(COLLECTION_NAME).Add(ctx, chatDoc)
	return err
}

func (c *Client) UpdateChat(ctx context.Context, chatId string, savedChatRecord SavedChatRecord) error {
	chatDoc := map[string]interface{}{
		"messages":    savedChatRecord.Messages,
		"memberId":    savedChatRecord.MemberId,
		"memberName":  savedChatRecord.MemberName,
		"lastUpdated": savedChatRecord.LastUpdated,
	}
	_, err := c.Collection(COLLECTION_NAME).Doc(chatId).Set(ctx, chatDoc, firestore.MergeAll)
	return err
}

func (c *Client) GetChat(ctx context.Context, chatId string) (*[]chatgpt.Message, error) {
	doc, err := c.Collection(COLLECTION_NAME).Doc(chatId).Get(ctx)
	if err != nil {
		println("Error getting chat document:", err)
		return &[]chatgpt.Message{}, nil
	}

	var chat SavedChatRecord
	if err := doc.DataTo(&chat); err != nil {
		return nil, err
	}
	return &chat.Messages, nil
}

type ListChatsParams struct {
	Limit  int
	Offset int
}

func (c *Client) ListChats(ctx context.Context, params *ListChatsParams) *[]SavedChatRecord {
	// Implement this method to list chats
	snapshot := c.Collection(COLLECTION_NAME).Limit(params.Limit).Offset(params.Offset).OrderBy("lastUpdated", firestore.Desc).Documents(ctx)
	var chats []SavedChatRecord

	for {
		doc, err := snapshot.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			println("Error getting chat document:", err)
			return nil
		}
		var chat SavedChatRecord
		if err := doc.DataTo(&chat); err != nil {
			println("Error parsing chat document:", err)
			return nil
		}
		chats = append(chats, chat)
	}
	return &chats
}
