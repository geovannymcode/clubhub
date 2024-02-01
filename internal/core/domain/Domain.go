package domain

import "time"

// Domain represents an internet domain with its address and the last consultation time.
type Domain struct {
	ID          int64
	Address     string
	LastConsult time.Time
}
