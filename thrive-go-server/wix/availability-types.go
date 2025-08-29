package wix

type AvailabilityQueryRequest struct {
	Query AvailabilityQuery `json:"query"`
}

type AvailabilityQuery struct {
	Filter AvailabilityFilter `json:"filter"`
}

type AvailabilityFilter struct {
	ServiceIDs []string `json:"serviceId"`
	Bookable   string   `json:"bookable"`
	StartDate  string   `json:"startDate"`
	EndDate    string   `json:"endDate"`
}

type AvailabilityQueryResponse struct {
	AvailabilityEntries []AvailabilityEntry `json:"availabilityEntries"`
}

type AvailabilityEntry struct {
	Slot                    Slot                    `json:"slot"`
	Bookable                bool                    `json:"bookable"`
	TotalSpots              int                     `json:"totalSpots"`
	OpenSpots               int                     `json:"openSpots"`
	WaitingList             map[string]interface{}  `json:"waitingList"`
	BookingPolicyViolations BookingPolicyViolations `json:"bookingPolicyViolations"`
}

type Slot struct {
	ServiceID  string   `json:"serviceId"`
	ScheduleID string   `json:"scheduleId"`
	StartDate  string   `json:"startDate"`
	EndDate    string   `json:"endDate"`
	Resource   Resource `json:"resource"`
	Location   Location `json:"location"`
}

func (s *Slot) ToString() string {
	return "Service ID: " + s.ServiceID + ", Start Date: " + s.StartDate + ", End Date: " + s.EndDate + ", Resource: " + s.Resource.Name
}

type Resource struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ScheduleID string `json:"scheduleId"`
}

type BookingPolicyViolations struct {
	TooEarlyToBook     bool `json:"tooEarlyToBook"`
	TooLateToBook      bool `json:"tooLateToBook"`
	BookOnlineDisabled bool `json:"bookOnlineDisabled"`
}
