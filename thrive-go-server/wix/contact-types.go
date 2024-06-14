package wix

type GetContactAPIResponse struct {
	Contact      Contact `json:"contact"`
	ResponseType string  `json:"responseType"`
}

type Contact struct {
	ID           string       `json:"id"`
	Revision     int          `json:"revision"`
	Source       Source       `json:"source"`
	CreatedDate  string       `json:"createdDate"`
	UpdatedDate  string       `json:"updatedDate"`
	LastActivity LastActivity `json:"lastActivity"`
	PrimaryInfo  PrimaryInfo  `json:"primaryInfo"`
	Picture      Picture      `json:"picture"`
	Info         ContactInfo  `json:"info"`
	PrimaryEmail PrimaryEmail `json:"primaryEmail"`
	PrimaryPhone PrimaryPhone `json:"primaryPhone"`
}

type Source struct {
	SourceType string `json:"sourceType"`
	WixAppId   string `json:"wixAppId"`
	AppId      string `json:"appId"`
}

type LastActivity struct {
	ActivityDate string `json:"activityDate"`
	ActivityType string `json:"activityType"`
	Date         string `json:"date"`
	Description  string `json:"description"`
	Icon         Icon   `json:"icon"`
}

type Icon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PrimaryInfo struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type Picture struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type ContactInfo struct {
	Name           Name           `json:"name"`
	Emails         Emails         `json:"emails"`
	Phones         Phones         `json:"phones"`
	Locale         string         `json:"locale"`
	LabelKeys      LabelKeys      `json:"labelKeys"`
	ExtendedFields ExtendedFields `json:"extendedFields"`
	Picture        Picture        `json:"picture"`
}

type Name struct {
	First string `json:"first"`
	Last  string `json:"last"`
}

type Emails struct {
	Items []EmailItem `json:"items"`
}

type EmailItem struct {
	ID      string `json:"id"`
	Tag     string `json:"tag"`
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
}

type Phones struct {
	Items []PhoneItem `json:"items"`
}

type PhoneItem struct {
	ID          string `json:"id"`
	Tag         string `json:"tag"`
	CountryCode string `json:"countryCode"`
	Phone       string `json:"phone"`
	E164Phone   string `json:"e164Phone"`
	Primary     bool   `json:"primary"`
}

type LabelKeys struct {
	Items []string `json:"items"`
}

type ExtendedFields struct {
	Items map[string]string `json:"items"`
}

type PrimaryEmail struct {
	Email                string `json:"email"`
	SubscriptionStatus   string `json:"subscriptionStatus"`
	DeliverabilityStatus string `json:"deliverabilityStatus"`
}

type PrimaryPhone struct {
	CountryCode          string `json:"countryCode"`
	E164Phone            string `json:"e164Phone"`
	FormattedPhone       string `json:"formattedPhone"`
	SubscriptionStatus   string `json:"subscriptionStatus"`
	DeliverabilityStatus string `json:"deliverabilityStatus"`
	Phone                string `json:"phone"`
}
