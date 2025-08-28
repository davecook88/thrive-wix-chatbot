package wix

import "fmt"

type PricingPlanQueryResponse struct {
	Plans          []PricingPlan `json:"plans"`
	PagingMetadata PagingMetadata `json:"pagingMetadata"`
}

type PricingPlan struct {
	ID                        string                     `json:"id"`
	Revision                  string                     `json:"revision"`
	CreatedDate               string                     `json:"createdDate"`
	UpdatedDate               string                     `json:"updatedDate"`
	Name                      string                     `json:"name"`
	Description               string                     `json:"description"`
	Image                     Image                      `json:"image"`
	Slug                      string                     `json:"slug"`
	TermsAndConditions      string                     `json:"termsAndConditions"`
	MaxPurchasesPerBuyer      int                        `json:"maxPurchasesPerBuyer"`
	PricingVariants           []PricingVariant           `json:"pricingVariants"`
	Perks                     []Perk                     `json:"perks"`
	Visibility                string                     `json:"visibility"`
	Buyable                   bool                       `json:"buyable"`
	Status                    string                     `json:"status"`
	BuyerCanCancel            bool                       `json:"buyerCanCancel"`
	Archived                  bool                       `json:"archived"`
	Primary                   bool                       `json:"primary"`
	Currency                  string                     `json:"currency"`
	FormType                  string                     `json:"formType,omitempty"`
	TermsAndConditionsSettings TermsAndConditionsSettings `json:"termsAndConditionsSettings"`
	DisplayIndex              int                        `json:"displayIndex"`
	ThankYouPageSettings      *ThankYouPageSettings      `json:"thankYouPageSettings,omitempty"`
	FormID                    string                     `json:"formId,omitempty"`
}


type PricingVariant struct {
	ID                string              `json:"id"`
	Name              string              `json:"name"`
	FreeTrialDays     int                 `json:"freeTrialDays"`
	Fees              []interface{}       `json:"fees"`
	BillingTerms      BillingTerms        `json:"billingTerms"`
	PricingStrategies []PricingStrategy `json:"pricingStrategies"`
}

type BillingTerms struct {
	BillingCycle         BillingCycle         `json:"billingCycle"`
	StartType            string               `json:"startType"`
	EndType              string               `json:"endType"`
	CyclesCompletedDetails *CyclesCompletedDetails `json:"cyclesCompletedDetails,omitempty"`
}

type BillingCycle struct {
	Period string `json:"period"`
	Count  string `json:"count"`
}

type CyclesCompletedDetails struct {
	BillingCycleCount string `json:"billingCycleCount"`
}

type PricingStrategy struct {
	FlatRate FlatRate `json:"flatRate"`
}

type FlatRate struct {
	Amount string `json:"amount"`
}

type Perk struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type TermsAndConditionsSettings struct {
	AcceptRequired    bool `json:"acceptRequired"`
	AcceptedByDefault bool `json:"acceptedByDefault"`
}

type ThankYouPageSettings struct {
	Title      string `json:"title"`
	Message    string `json:"message"`
	ButtonText string `json:"buttonText"`
	ButtonLink string `json:"buttonLink,omitempty"`
}

type PagingMetadata struct {
	Count    int         `json:"count"`
	Cursors  Cursors     `json:"cursors"`
	HasNext  bool        `json:"hasNext"`
}

type Cursors struct {
}

func (p *PricingPlan) ToString() string {
	return fmt.Sprintf(`
		ID: %s
		Name: %s
		Description: %s
		Currency: %s
		Pricing Variant: %s - %s %s
	`, p.ID, p.Name, p.Description, p.Currency, p.PricingVariants[0].Name, p.PricingVariants[0].PricingStrategies[0].FlatRate.Amount, p.Currency)
}