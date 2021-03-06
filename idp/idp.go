/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

                 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package idp contains Identity Provider APIs as used by Hyperledger Fabric
package idp

import (
	"time"

	"github.com/cloudflare/cfssl/csr"
)

/*
 * The identity provider APIs are split into 3 main interfaces:
 * 1) ClientAPI - the interface used by a client SDK;
 * 2) PeerAPI - the interface used by a peer;
 * 3) ChaincodeAPI - the interface used by a chaincode.
 */

// ClientAPI is the API used by the client SDK to interface with the IDP
type ClientAPI interface {

	// Capabilities returns the capabilities of this IDP client
	Capabilities() []Capability

	// Register a new identity
	// @param req The registration request
	Register(req *RegistrationRequest) (*RegistrationResponse, error)

	// Enroll a new identity
	// @param req The enrollment request
	Enroll(req *EnrollmentRequest) (Identity, error)

	// RegisterAndEnroll registers and enrolls a new identity
	// @param req The registration request
	RegisterAndEnroll(req *RegistrationRequest) (Identity, error)

	// ImportSigner imports a signer from an external CA
	// @param req The import request
	ImportSigner(req *ImportSignerRequest) (Signer, error)

	// DeserializeIdentity deserializes an identity
	DeserializeIdentity([]byte) (Identity, error)
}

// PeerAPI is the API used by the peer pertaining to the IDP
type PeerAPI interface {

	// A peer can sign and verify
	Signer
}

// ChaincodeAPI is the API used by the chaincode pertaining to the IDP
type ChaincodeAPI interface {

	// TODO: For ZRL to define

}

// Identity represents an arbitrary identity
type Identity interface {

	// Name returns the identity name
	GetName() string

	// GetPublicSigner returns the public signer for this identity
	GetPublicSigner() TemporalSigner

	// GetPrivateSigners returns private signers for this identity
	GetPrivateSigners(req *GetPrivateSignersRequest) ([]TemporalSigner, error)

	// GetAttributeNames returns the names of all attributes associated with this identity
	GetAttributeNames() ([]string, error)

	// Delete this identity completely and revoke all of it's signers
	Delete() error

	// Serialize an identity
	Serialize() ([]byte, error)
}

// TemporalSigner is a signer which can be renewed and revoked
type TemporalSigner interface {

	// Extends Signer
	Signer

	// Renew this identity
	Renew() error

	// Revoke this identity
	Revoke() error
}

// Signer interface
type Signer interface {

	// Extends Verifier
	Verifier

	// Sign the message
	Sign(msg []byte) ([]byte, error)

	// SignOpts the message with options
	SignOpts(msg []byte, opts SignatureOpts) ([]byte, error)

	// NewAttributeProof creates a proof for an attribute
	NewAttributeProof(spec *AttributeProofSpec) (proof []byte, err error)
}

// Verifier interface
type Verifier interface {

	// Verify myself
	VerifySelf() error

	// Verify a message given a signature over the message
	Verify(msg []byte, sig []byte) error

	// Verify a signature over some message with specific options
	VerifyOpts(msg []byte, sig []byte, opts SignatureOpts) error

	// VerifyAttributes verifies attributes given proofs
	VerifyAttributes(proof [][]byte, spec *AttributeProofSpec) error

	// Serialize verifier
	Serialize() ([]byte, error)
}

// RegistrationRequest for a new identity
type RegistrationRequest struct {
	// Name is the unique name of the identity
	Name string `json:"id"`
	// Type of identity being registered (e.g. "peer, app, user")
	Type string `json:"type"`
	// Group name associated with the identity
	Group string `json:"group"`
	// Attributes associated with this identity
	Attributes []Attribute `json:"attrs,omitempty"`
	// Registrar is the identity that is performing the registration
	Registrar Identity `json:"registrar"`
}

// RegistrationResponse is a registration response
type RegistrationResponse struct {
	// The optional secret returned from a registration response
	Secret string `json:"credential,omitempty"`
}

// EnrollmentRequest is a request to enroll an identity
type EnrollmentRequest struct {
	// The identity name to enroll
	Name string `json:"name"`
	// The secret returned via Register
	Secret string `json:"secret,omitempty"`
	// The hosts
	Hosts []string `json:"hosts,omitempty"`
	// The certificate request
	CR *csr.CertificateRequest `json:"cr,omitempty"`
	// The attributes to include in the enrollment certificate (optional)
	Attributes []Attribute `json:"attrs,omitempty"`
}

// ImportSignerRequest is required when importing a signer from an external CA
type ImportSignerRequest struct {
	// The certificate to import
	Cert []byte `json:"cert"`
	// The private key to import (optional)
	Key []byte `json:"key,omitempty"`
}

// GetPrivateSignersRequest is input provided to get private signers
type GetPrivateSignersRequest struct {
	Count          int       `json:"count"`
	AttrNames      []string  `json:"attrNames,omitempty"`
	ValidityPeriod time.Time `json:"validityPeriod,omitempty"`
}

// SignatureOpts are signature options
type SignatureOpts interface {
	Policy() []string
	Label() string
}

// Attribute is an arbitrary name/value pair
type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// AttributeProofSpec is an attribute proof specification
type AttributeProofSpec struct {
	Attributes []string
	Message    []byte
}

// Capability is a capability of a CA
type Capability int

// The capabilities of a CA relative to the CA API
const (
	REGISTRATION  Capability = iota // CA has registrar capability
	ENROLLMENT                      // CA has enrollment capability
	ATTRIBUTES                      // CA has attributes capability
	ANONYMITY                       // CA support anonymous identities
	UNLINKABILITY                   // CA support unlinkable identities
)
