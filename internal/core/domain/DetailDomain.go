package domain

import "time"

// DetailDomain details technical aspects of a Domain, like IP address, server name, and grade.
type DetailDomain struct {
	ID         int64
	DomainID   int64
	IpAddress  string
	ServerName string
	Grade      string
	Date       time.Time
}
