package wix

import (
	"fmt"
	"time"
)

type QueryServicesQuery struct {
	Filter interface{} `json:"filter"`
}

type QueryServicesRequest struct {
	Query QueryServicesQuery `json:"query"`
}

func NewQueryServicesRequest(filter interface{}) *QueryServicesRequest {
	return &QueryServicesRequest{
		Query: QueryServicesQuery{
			Filter: filter,
		},
	}
}

type Media struct {
	Items      []struct{} `json:"items"`
	MainMedia  Image      `json:"mainMedia"`
	CoverMedia Image      `json:"coverMedia"`
}

type CalculatedAddress struct {
	Country          string  `json:"country"`
	City             string  `json:"city"`
	FormattedAddress string  `json:"formattedAddress"`
	Geocode          Geocode `json:"geocode"`
}

type Geocode struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Business struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Default bool              `json:"default"`
	Address CalculatedAddress `json:"address"`
}

type Location struct {
	Type              string            `json:"type"`
	CalculatedAddress CalculatedAddress `json:"calculatedAddress"`
	Business          Business          `json:"business"`
}

type FixedPrice struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type Payment struct {
	RateType string     `json:"rateType"`
	Fixed    FixedPrice `json:"fixed"`
	Options  struct {
		Online      bool `json:"online"`
		InPerson    bool `json:"inPerson"`
		Deposit     bool `json:"deposit"`
		PricingPlan bool `json:"pricingPlan"`
	} `json:"options"`
	PricingPlanIds []string `json:"pricingPlanIds"`
}

type BookingPolicy struct {
	ID                      string `json:"id"`
	Revision                string `json:"revision"`
	CreatedDate             string `json:"createdDate"`
	UpdatedDate             string `json:"updatedDate"`
	Name                    string `json:"name"`
	CustomPolicyDescription struct {
		Enabled     bool   `json:"enabled"`
		Description string `json:"description"`
	} `json:"customPolicyDescription"`
	Default                 bool `json:"default"`
	LimitEarlyBookingPolicy struct {
		Enabled                  bool `json:"enabled"`
		EarliestBookingInMinutes int  `json:"earliestBookingInMinutes"`
	} `json:"limitEarlyBookingPolicy"`
	LimitLateBookingPolicy struct {
		Enabled                bool `json:"enabled"`
		LatestBookingInMinutes int  `json:"latestBookingInMinutes"`
	} `json:"limitLateBookingPolicy"`
	BookAfterStartPolicy struct {
		Enabled bool `json:"enabled"`
	} `json:"bookAfterStartPolicy"`
	CancellationPolicy struct {
		Enabled                     bool `json:"enabled"`
		LimitLatestCancellation     bool `json:"limitLatestCancellation"`
		LatestCancellationInMinutes int  `json:"latestCancellationInMinutes"`
	} `json:"cancellationPolicy"`
	ReschedulePolicy struct {
		Enabled                   bool `json:"enabled"`
		LimitLatestReschedule     bool `json:"limitLatestReschedule"`
		LatestRescheduleInMinutes int  `json:"latestRescheduleInMinutes"`
	} `json:"reschedulePolicy"`
	WaitlistPolicy struct {
		Enabled                  bool `json:"enabled"`
		Capacity                 int  `json:"capacity"`
		ReservationTimeInMinutes int  `json:"reservationTimeInMinutes"`
	} `json:"waitlistPolicy"`
	ParticipantsPolicy struct {
		Enabled                   bool `json:"enabled"`
		MaxParticipantsPerBooking int  `json:"maxParticipantsPerBooking"`
	} `json:"participantsPolicy"`
	ResourcesPolicy struct {
		Enabled           bool `json:"enabled"`
		AutoAssignAllowed bool `json:"autoAssignAllowed"`
	} `json:"resourcesPolicy"`
}

type Schedule struct {
	ID                string    `json:"id"`
	FirstSessionStart time.Time `json:"firstSessionStart"`
}

type SupportedSlug struct {
	Name        string `json:"name"`
	Custom      bool   `json:"custom"`
	CreatedDate string `json:"createdDate"`
}

type URL struct {
	RelativePath string `json:"relativePath"`
	URL          string `json:"url"`
}

type SEOData struct {
	Tags []struct {
		Type  string `json:"type"`
		Props struct {
			Name    string `json:"name"`
			Content string `json:"content"`
		} `json:"props"`
		Children string `json:"children"`
		Custom   bool   `json:"custom"`
		Disabled bool   `json:"disabled"`
	} `json:"tags"`
	Settings struct {
		PreventAutoRedirect bool     `json:"preventAutoRedirect"`
		Keywords            []any `json:"keywords"`
	} `json:"settings"`
}

type Service struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	SortOrder       int    `json:"sortOrder"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	TagLine         string `json:"tagLine"`
	DefaultCapacity int    `json:"defaultCapacity"`
	Media           Media  `json:"media"`
	Hidden          bool   `json:"hidden"`
	Category        struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		SortOrder int    `json:"sortOrder"`
	} `json:"category"`
	Form struct {
		ID             string `json:"id"`
		MobileSettings struct {
			Hidden bool `json:"hidden"`
		} `json:"mobileSettings"`
	} `json:"form"`
	Payment       Payment `json:"payment"`
	OnlineBooking struct {
		Enabled               bool `json:"enabled"`
		RequireManualApproval bool `json:"requireManualApproval"`
		AllowMultipleRequests bool `json:"allowMultipleRequests"`
	} `json:"onlineBooking"`
	Conferencing struct {
		Enabled bool `json:"enabled"`
	} `json:"conferencing"`
	Locations        []Location      `json:"locations"`
	BookingPolicy    BookingPolicy   `json:"bookingPolicy"`
	Schedule         Schedule        `json:"schedule"`
	StaffMemberIds   []string        `json:"staffMemberIds"`
	StaffMembers     []interface{}   `json:"staffMembers"`
	ResourceGroups   []interface{}   `json:"resourceGroups"`
	ServiceResources []interface{}   `json:"serviceResources"`
	SupportedSlugs   []SupportedSlug `json:"supportedSlugs"`
	MainSlug         SupportedSlug   `json:"mainSlug"`
	URLs             struct {
		ServicePage  URL `json:"servicePage"`
		BookingPage  URL `json:"bookingPage"`
		CalendarPage URL `json:"calendarPage"`
	} `json:"urls"`
	SEOData     SEOData `json:"seoData"`
	CreatedDate string  `json:"createdDate"`
	UpdatedDate string  `json:"updatedDate"`
	Revision    string  `json:"revision"`
}

func (s *Service) ToString() string {
	return fmt.Sprintf(`
		ID: %s
		Type: %s
		Name: %s
		Description: %s
		DefaultCapacity: %d
		ServicePage: %s
		BookingPage: %s`,
		s.ID, s.Type, s.Name, s.Description, s.DefaultCapacity, s.URLs.ServicePage.URL, s.URLs.BookingPage.URL)
}

type QueryServicesResponse struct {
	Services []Service `json:"services"`
}
