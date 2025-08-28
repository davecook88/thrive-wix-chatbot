package db

import (
	"context"
	"thrive/server/wix"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const WIX_PRICING_PLANS_DOC = "wix-pricing-plans"
const WIX_PRICING_PLANS_COLLECTION = "wix-pricing-plans"

type CachedWixPricingPlans struct {
	PricingPlans []wix.PricingPlan `json:"pricingPlans"`
	LastUpdated  string            `json:"lastUpdated"`
}

func (c *Client) GetWixPricingPlans(ctx context.Context) (*[]wix.PricingPlan, error) {
	doc, err := c.Collection(WIX_PRICING_PLANS_COLLECTION).Doc(WIX_PRICING_PLANS_DOC).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			println("No cached pricing plans found. Fetching from Wix in background...")
			go func() {
				wixClient := wix.NewWixClient()
				pricingPlans, err := wixClient.QueryPricingPlans()
				if err != nil {
					println("Error fetching pricing plans from Wix:", err)
					return
				}
				c.saveWixPricingPlans(context.Background(), pricingPlans)
			}()
			return nil, nil
		}
		return nil, err
	}

	var cachedPricingPlans CachedWixPricingPlans
	if err := doc.DataTo(&cachedPricingPlans); err != nil {
		return nil, err
	}

	lastUpdated, err := time.Parse(time.RFC3339, cachedPricingPlans.LastUpdated)
	if err != nil {
		return nil, err
	}

	if time.Since(lastUpdated).Hours() > 24 {
		println("Cached pricing plans are stale. Fetching from Wix in background...")
		go func() {
			wixClient := wix.NewWixClient()
			pricingPlans, err := wixClient.QueryPricingPlans()
			if err != nil {
				println("Error fetching pricing plans from Wix:", err)
				return
			}
			c.saveWixPricingPlans(context.Background(), pricingPlans)
		}()
	}

	return &cachedPricingPlans.PricingPlans, nil
}

func (c *Client) saveWixPricingPlans(ctx context.Context, pricingPlans *[]wix.PricingPlan) {
	 doc := map[string]interface{}{
		"pricingPlans": pricingPlans,
		"lastUpdated":  time.Now().Format(time.RFC3339),
	}
	_, err := c.Collection(WIX_PRICING_PLANS_COLLECTION).Doc(WIX_PRICING_PLANS_DOC).Set(ctx, doc, firestore.MergeAll)
	if err != nil {
		println("Error saving wix pricing plans to cache:", err)
	}
	println("Successfully cached wix pricing plans")
}
