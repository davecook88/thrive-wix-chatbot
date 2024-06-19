package handlers

import (
	"encoding/json"
	"fmt"
	"strings"
	"thrive/server/chatgpt"
	"thrive/server/wix"
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

func handleGetServices() *string {
	wixClient := wix.NewWixClient()
	fmt.Println("got wix client in get services")
	filter := map[string]interface{}{
		"hidden": false,
	}

	// Create a new QueryServicesRequest instance
	request := wix.NewQueryServicesRequest(filter)
	services, err := wixClient.QueryServices(request)
	if err != nil {
		fmt.Println("failed to get services", err)
		return new(string)
	}
	// print the json of the services
	responseStr := ""
	for _, service := range *services {
		responseStr += service.ToString() + "\n"
	}

	return &responseStr

}

func handleToolCallsWithResponse(toolCalls *[]chatgpt.ToolCall) *string {
	for _, toolCall := range *toolCalls {
		switch toolCall.Function.Name {
		case "getServices":
			return handleGetServices()
		}
	}
	return new(string)
}

func handleToolCalls(toolCalls *[]chatgpt.ToolCall, wixMember *wix.WixMember) {
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
			// Arguments should be a JSON string like this {\"estimatedLevel\":\"C1\"}
			arguments := map[string]string{}
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments); err != nil {
				fmt.Println("failed to unmarshal arguments", err)
				return
			}
			handleLevelEstimation(arguments["estimatedLevel"], contact)

		case "addNotes":
			// Arguments should be a JSON string like this {\"notes\":\"This user is interested in learning Spanish for travel.\"}
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
			handleGetServices()

		default:
			continue
		}
	}
	fmt.Println("updating contact")
	wixClient.UpdateContact(contact.ID, contact.Revision, contact.Info)

}
