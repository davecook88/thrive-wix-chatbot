package db

import (
	"context"
	"thrive/server/wix"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const WIX_SERVICES_DOC = "wix-services"

func (c *Client) GetWixServices(ctx context.Context) (*[]wix.Service, error) {
	doc, err := c.Collection(WIX_SERVICES_COLLECTION).Doc(WIX_SERVICES_DOC).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			println("No cached services found. Fetching from Wix in background...")
			go func() {
				wixClient := wix.NewWixClient()
				services, err := wixClient.QueryServices(wix.NewQueryServicesRequest(map[string]interface{}{"hidden": false}))
				if err != nil {
					println("Error fetching services from Wix:", err)
					return
				}
				c.saveWixServices(context.Background(), services)
			}()
			return nil, nil
		}
		return nil, err
	}

	var cachedServices CachedWixServices
	if err := doc.DataTo(&cachedServices); err != nil {
		return nil, err
	}

	lastUpdated, err := time.Parse(time.RFC3339, cachedServices.LastUpdated)
	if err != nil {
		return nil, err
	}

	if time.Since(lastUpdated).Hours() > 24 {
		println("Cached services are stale. Fetching from Wix in background...")
		go func() {
			wixClient := wix.NewWixClient()
			services, err := wixClient.QueryServices(wix.NewQueryServicesRequest(map[string]interface{}{"hidden": false}))
			if err != nil {
				println("Error fetching services from Wix:", err)
				return
			}
			c.saveWixServices(context.Background(), services)
		}()
	}

	return &cachedServices.Services, nil
}

func (c *Client) saveWixServices(ctx context.Context, services *[]wix.Service) {
	 doc := map[string]interface{}{
		"services":    services,
		"lastUpdated": time.Now().Format(time.RFC3339),
	}
	_, err := c.Collection(WIX_SERVICES_COLLECTION).Doc(WIX_SERVICES_DOC).Set(ctx, doc, firestore.MergeAll)
	if err != nil {
		println("Error saving wix services to cache:", err)
	}
	println("Successfully cached wix services")
}
