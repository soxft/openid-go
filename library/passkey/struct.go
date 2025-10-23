package passkey

type KeyCreatePrepare struct {
	Challenge              string                 `json:"challenge"`
	Rp                     Rp                     `json:"rp"`
	User                   User                   `json:"user"`
	PubKeyCredParams       []PubKeyCredParams     `json:"pubKeyCredParams"`
	Timeout                int                    `json:"timeout"`
	Attestation            string                 `json:"attestation"`
	AuthenticatorSelection AuthenticatorSelection `json:"authenticatorSelection"`
	Extensions             Extensions             `json:"extensions"`
}

type Rp struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type User struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type PubKeyCredParams struct {
	Alg  int    `json:"alg"`
	Type string `json:"type"`
}

type AuthenticatorSelection struct {
	AuthenticatorAttachment string `json:"authenticatorAttachment"`
	RequireResidentKey      bool   `json:"requireResidentKey"`
	ResidentKey             string `json:"residentKey"`
}

type Extensions struct {
	CredProps bool `json:"credProps"`
}
