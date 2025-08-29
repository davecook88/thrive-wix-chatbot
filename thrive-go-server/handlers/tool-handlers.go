package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"thrive/server/chatgpt"
	"thrive/server/db"
	"thrive/server/wix"
	"time"
)

func handleLevelEstimation(level string, wixContact *wix.Contact) {
	existingLabels := wixContact.Info.LabelKeys
	// remove any existing level labels, they begin with custom.level-
	filteredLabels := []string{}
	for _, label := range existingLabels.Items {
		if !strings.HasPrefix(label, "custom.level-") {
			filteredLabels = append(filteredLabels, label)
		}
	}
	// append new level label
	filteredLabels = append(filteredLabels, "custom.level-"+strings.ToLower(level))
	wixContact.Info.LabelKeys.Items = filteredLabels
}

func handleGetServices(ctx context.Context, dbClient *db.Client) *string {
	services, err := dbClient.GetWixServices(ctx)
	if err != nil {
		fmt.Println("failed to get services", err)
		return new(string)
	}
	// print the json of the services
	responseStr := ""
	for _, service := range *services {
		responseStr += service.ToString() + "\n ##SERVICE_SEPARATOR \n"
	}

	return &responseStr

}

func handleGetAvailability(ctx context.Context, dbClient *db.Client, serviceIDs []string) *string {
	wixClient := wix.NewWixClient()

	if serviceIDs == nil {
		services, err := dbClient.GetWixServices(ctx)
		if err != nil {
			fmt.Println("failed to get services", err)
			return new(string)
		}
		for _, service := range *services {
			serviceIDs = append(serviceIDs, service.ID)
		}
	}

	request := &wix.AvailabilityQueryRequest{
		Query: wix.AvailabilityQuery{
			Filter: wix.AvailabilityFilter{
				ServiceIDs: serviceIDs,
				Bookable:   "true",
				StartDate:  time.Now().Format(time.RFC3339),
				EndDate:    time.Now().AddDate(0, 0, 7).Format(time.RFC3339),
			},
		},
	}

	availability, err := wixClient.QueryAvailability(request)
	if err != nil {
		fmt.Println("failed to get availability", err)
		return new(string)
	}

	responseStr := ""
	for _, entry := range availability.AvailabilityEntries {
		responseStr += entry.Slot.ToString() + "\n"
	}

	return &responseStr
}

func handleToolCallsWithResponse(ctx context.Context, toolCalls *[]chatgpt.ToolCall, dbClient *db.Client) *string {
	for _, toolCall := range *toolCalls {
		switch toolCall.Function.Name {
		case "getServices":
			return handleGetServices(ctx, dbClient)
		case "getAvailability":
			arguments := make(map[string][]string)
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments); err != nil {
				fmt.Println("failed to unmarshal arguments", err)
				return new(string)
			}
			return handleGetAvailability(ctx, dbClient, arguments["serviceIds"])
		}
	}
	return new(string)
}

func handleToolCalls(toolCalls *[]chatgpt.ToolCall, wixMember *wix.WixMember, dbClient *db.Client) {
	wixClient := wix.NewWixClient()
	fmt.Println("got wix client")
	contact, err := wixClient.GetContact(wixMember.ContactId)
	if err != nil {
		fmt.Println("failed to get contact", err)
		return
	}
	for _, toolCall := range *toolCalls {
		switch toolCall.Function.Name {
		case "estimateUserLevel":
			// Arguments should be a JSON string like this {"estimatedLevel":"C1"}
			arguments := map[string]string{}
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments); err != nil {
				fmt.Println("failed to unmarshal arguments", err)
				return
			}
			handleLevelEstimation(arguments["estimatedLevel"], contact)

		case "addNotes":
			// Arguments should be a JSON string like this {"notes":"This user is interested in learning Spanish for travel."}
			arguments := map[string]string{}
			fmt.Println("addNotes", toolCall.Function.Arguments)
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments); err != nil {
				fmt.Println("failed to unmarshal arguments", err)
				return
			}
			// Add notes to the contact
			existingNotes := contact.Info.ExtendedFields.Items["custom.notes"]
			newNotes := existingNotes + "\n*" + arguments["notes"] + "*"
			contact.Info.ExtendedFields.Items["custom.notes"] = newNotes

		case "getServices":
			handleGetServices(context.Background(), dbClient)

		default:
			continue
		}
	}
	fmt.Println("updating contact")
	wixClient.UpdateContact(contact.ID, contact.Revision, contact.Info)

}
