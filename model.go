package main

// Account is used to store lets encrypt registration info
type Account struct {
	//Email              string
	//Registration       *acme.RegistrationResource
	//PrivateKey         []byte
	DomainsCertificate DomainsCertificates
	//ChallengeCerts     map[string]*ChallengeCert
}

// DomainsCertificates stores a certificate for multiple domains
type DomainsCertificates struct {
	Certs []*DomainsCertificate
}

// DomainsCertificate contains a certificate for multiple domains
type DomainsCertificate struct {
	Domains     Domain
	Certificate *Certificate
	//tlsCert     *tls.Certificate
}

// Domain holds a domain name with SANs
type Domain struct {
	Main string
	SANs []string
}

// Certificate is used to store certificate info
type Certificate struct {
	Domain        string
	CertURL       string
	CertStableURL string
	PrivateKey    []byte
	Certificate   []byte
}
