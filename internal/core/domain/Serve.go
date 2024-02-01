package domain

// Serve contains details about a server, including address, SSL grade, country, and owner.
type Serve struct {
	Address  string `json:"address"`
	SslGrade string `json:"ssl_grade"`
	Country  string `json:"country"`
	Owner    string `json:"owner"`
}
