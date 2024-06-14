package wix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type WixClient struct {
	Token  string
	SiteId string
	*http.Client
}

func NewWixClient() *WixClient {
	return &WixClient{
		Token:  os.Getenv("WIX_TOKEN"),
		SiteId: os.Getenv("WIX_SITE_ID"),
		Client: &http.Client{},
	}
}

func (c *WixClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Wix-Site-Id", c.SiteId)
	return c.Client.Do(req)
}

func (c *WixClient) GetMember(memberId string) (*WixMember, error) {
	url := "https://www.wixapis.com/members/v1/members/" + memberId
	println("url", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		println("failed to create request", err)
		return nil, err
	}
	resp, err := c.Do(req)
	if err != nil {
		println("failed to make request", err)
		return nil, err
	}

	defer resp.Body.Close()

	var memberResp WixGetMemberResponse
	if err := json.NewDecoder(resp.Body).Decode(&memberResp); err != nil {
		println("failed to parse", err)
		return nil, err
	}
	fmt.Println("member", memberResp.Member)
	return &memberResp.Member, nil

}

func (c *WixClient) GetContact(contactId string) (*Contact, error) {
	url := "https://www.wixapis.com/contacts/v4/contacts/" + contactId
	fmt.Println("url", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req)
	if err != nil {
		fmt.Println("failed to make request", err)
		return nil, err
	}
	defer resp.Body.Close()

	var contact GetContactAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&contact); err != nil {
		return nil, err
	}
	return &contact.Contact, nil
}

type PatchInfo struct {
	Revision int         `json:"revision"`
	Info     ContactInfo `json:"info"`
}

func (c *WixClient) UpdateContact(contactId string, revision int, info ContactInfo) (*Contact, error) {
	url := "https://www.wixapis.com/contacts/v4/contacts/" + contactId
	// turn info into json
	putInfo := PatchInfo{
		Revision: revision,
		Info:     info,
	}
	jsonInfo, err := json.Marshal(putInfo)
	if err != nil {
		fmt.Println("failed to marshal info", err)
	}
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonInfo))
	if err != nil {
		fmt.Println("failed to create request", err)
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		fmt.Println("failed to make request", err)
		return nil, err
	}
	defer resp.Body.Close()

	var contactResp GetContactAPIResponse
	if err := json.NewDecoder(req.Body).Decode(&contactResp); err != nil {
		return nil, err
	}
	return &contactResp.Contact, nil

}
