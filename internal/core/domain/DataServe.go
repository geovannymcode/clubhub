package domain

// DataServe aggregates server data, SSL grades, and site status including if servers have changed.
type DataServe struct {
	Serves           []Serve `json:"servers"`
	ServersChanged   bool    `json:"servers_changed"`
	SslGrade         string  `json:"ssl_grade"`
	PreviousSslGrade string  `json:"previous_ssl_grade"`
	Logo             string  `json:"logo"`
	Title            string  `json:"title"`
	IsDown           bool    `json:"is_down"`
}
