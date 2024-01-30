package domain

import "time"

type DetailDomain struct {
	ID         int64
	DomainID   int64
	IpAddress  string
	ServerName string
	Grade      string
	Date       time.Time
}
