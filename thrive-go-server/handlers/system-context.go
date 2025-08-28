package handlers

import (
	"context"
	"thrive/server/db"
)

func GenerateSystemContext(dbClient *db.Client, ctx context.Context) string {
	wixServices, err := dbClient.GetWixServices(ctx)
	if err != nil {
		println("Error fetching Wix services:", err)
		return ""
	}

	if wixServices == nil || len(*wixServices) == 0 {
		println("No Wix services found.")
		return ""
	}

	// Generate the system context using the fetched Wix services
	systemContext := "Wix Services:\n"
	for _, service := range *wixServices {
		systemContext += service.ToString() + "\n"
	}
	return systemContext
}
