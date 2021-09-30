package models

type PageDetailsResponse struct {
	HtmlVersion string `json:"html_version"`
	PageTitle   string `json:"page_title"`
	LinksCount  struct {
		InternalLinks int `json:"internal_links_count"`
		ExternalLinks int `json:"external_links_count"`
	} `json:"links_count"`
	InaccessibleLinks int  `json:"inaccessible_links"`
	LoginFormExists   bool `json:"login_form_exists"`
}

type PageDetailsRequest struct {
	Url string `json:"url"`
}
