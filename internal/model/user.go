package model

type ABMUser struct {
	ExternalID int32  `json:"externalId"`
	Run        string `json:"run"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
}

type UserAccess struct {
	Location int8    `json:"location"`
	EntryAt  *string `json:"entryAt"`
	ExitAt   *string `json:"exitAt"`
}

type User struct {
	ImageURL      *string       `json:"imageUrl"`
	Run           string        `json:"run"`
	FirstName     string        `json:"firstName"`
	LastName      string        `json:"lastName"`
	AccessHistory []*UserAccess `json:"accessHistory"`
}
