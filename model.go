package main

// Domain holds a domain name with SANs
type Domain struct {
	Main string
	SANs []string
}

// Certificate is used to store certificate info
type Certificate struct {
	Domain      Domain
	Certificate []byte
	Key         []byte
}

// Certificates holds one or more certificates
type Certificates struct {
	Certificates []*Certificate
}
