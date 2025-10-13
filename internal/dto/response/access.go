package response

type AccessDTO struct {
	ExternalID string  `json:"externalId"`
	Run        string  `json:"run"`
	FullName   string  `json:"fullName"`
	Location   string  `json:"location"`
	EntryAt    string  `json:"entryAt"`
	ExitAt     *string `json:"exitAt"`
}
