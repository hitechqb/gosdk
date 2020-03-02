package sdkcm

type IPagingFilter interface {
	SetSkip(skip int)
	SetPerPage(perPage int)
}

type Paging struct {
	Total        int     `json:"total"`
	CurrentPage  int     `json:"current_page"`
	PerPage      int     `json:"per_page"`
	LastPage     int     `json:"last_page"`
	FirstPageURL *string `json:"first_page_url"`
	LastPageURL  *string `json:"last_page_url"`
	NextPageURL  *string `json:"next_page_url"`
	PrevPageURL  *string `json:"prev_page_url"`
	From         int     `json:"from"`
	To           int     `json:"to"`
}
