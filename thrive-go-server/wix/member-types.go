package wix

type WixMember struct {
	ID             string  `json:"id"`
	ContactId      string  `json:"contactId"`
	Status         string  `json:"status"`
	Profile        Profile `json:"profile"`
	PrivacyStatus  string  `json:"privacyStatus"`
	ActivityStatus string  `json:"activityStatus"`
	CreatedDate    string  `json:"createdDate"`
	UpdatedDate    string  `json:"updatedDate"`
}

type Profile struct {
	Nickname string `json:"nickname"`
	Slug     string `json:"slug"`
	Photo    Image  `json:"photo"`
	Cover    Image  `json:"cover"`
	Title    string `json:"title"`
}

type Image struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type WixGetMemberResponse struct {
	Member WixMember `json:"member"`
}
