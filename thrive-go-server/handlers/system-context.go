package handlers

import (
	"context"
	"fmt"
	"thrive/server/db"
	"thrive/server/wix"
)

func GenerateSystemContext(dbClient *db.Client, ctx context.Context) string {
	wixServices, err := dbClient.GetWixServices(ctx)
	if err != nil {
		println("Error fetching Wix services:", err)
		return ""
	}

	wixPricingPlans, err := dbClient.GetWixPricingPlans(ctx)
	if err != nil {
		println("Error fetching Wix pricing plans:", err)
		return ""
	}

	if wixServices == nil || len(*wixServices) == 0 {
		println("No Wix services found.")
		return ""
	}

	// Create a map of pricing plans by ID for easy lookup
	pricingPlansMap := make(map[string]wix.PricingPlan)
	if wixPricingPlans != nil {
		for _, plan := range *wixPricingPlans {
			pricingPlansMap[plan.ID] = plan
		}
	}

	// Generate the system context using the fetched Wix services and pricing plans
	systemContext := "Wix Services and Pricing Plans:\n"
	for _, service := range *wixServices {
		systemContext += service.ToString() + "\n"
		if service.Payment.PricingPlanIds != nil {
			systemContext += "\tPricing Plans Available:\n"
		}
		for _, servicePricingPlan := range service.Payment.PricingPlanIds {
			if plan, ok := pricingPlansMap[servicePricingPlan]; ok {
				systemContext += fmt.Sprintf("\t%s\n", plan.ToString())
			}
		}
	}

	return systemContext
}
