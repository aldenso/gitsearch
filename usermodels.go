package main

//RespUser struct for http response
type RespUser struct {
	NextURL           string     `json:"-"`
	Count             int        `json:"total_count"`
	IncompleteResults bool       `json:"incomplete_results"`
	Items             []ItemUser `json:"items"`
}

//ItemUser struct for user request results
type ItemUser struct {
	Login             string  `json:"login"`
	ID                int     `json:"id"`
	AvatarURL         string  `json:"avatar_url"`
	GravatarID        string  `json:"gravatar_id"`
	URL               string  `json:"url"`
	HTMLURL           string  `json:"html_url"`
	FollowersURL      string  `json:"followers_url"`
	FollowingURL      string  `json:"following_url"`
	GistsURL          string  `json:"gists_url"`
	StarredURL        string  `json:"starred_url"`
	SubscriptionsURL  string  `json:"subscriptions_url"`
	OrganizationsURL  string  `json:"organizations_url"`
	ReposURL          string  `json:"repos_url"`
	EventsURL         string  `json:"events_url"`
	ReceivedEventsURL string  `json:"received_events_url"`
	Type              string  `json:"type"`
	SiteAdmin         bool    `json:"site_admin"`
	Score             float64 `json:"score"`
}
