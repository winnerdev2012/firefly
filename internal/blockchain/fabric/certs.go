// Copyright © 2021 Kaleido, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fabric

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
)

func getDNFromCertString(certString string) (string, error) {
	cert, err := getCertificateFromBytes(certString)
	if err != nil {
		return "", err
	}
	return getDN(&cert.Subject), nil
}

// borrowed from fabric-chaincode-go to guarantee the same
// resolution of "DN" string from x509 certs
func getDN(name *pkix.Name) string {
	r := name.ToRDNSequence()
	return r.String()
}

func getCertificateFromBytes(certString string) (*x509.Certificate, error) {
	idbytes, err := base64.StdEncoding.DecodeString(certString)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(idbytes)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}
