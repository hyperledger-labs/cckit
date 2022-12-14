package testdata

import (
	idtestdata "github.com/hyperledger-labs/cckit/identity/testdata"
)

var (
	Certificates = idtestdata.Certs{{
		CertFilename: `admin.pem`, PKeyFilename: `admin.key.pem`,
	}}.
		UseReadFile(idtestdata.ReadLocal())
)
