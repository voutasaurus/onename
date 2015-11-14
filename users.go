package main

type Users map[string]User

type User struct {
	Profile       Profile
	Verifications []Verification
}

type Verification struct {
	Identifier string `json:"identifier"`
	ProofURL   string `json:"proof_url"`
	Service    string `json:"service"`
	Valid      bool   `json:"valid"`
}

type Url struct {
	Url string `json:"url"`
}

type SocialID struct {
	Proof    Url    `json:"proof"`
	Username string `json:"username"`
}

type Formatted struct {
	Formatted string `json:"formatted"`
}

type Profile struct {
	Avatar  Url `json:"avatar"`
	Bio     string
	Bitcoin struct {
		Address string `json:"address"`
	} `json:"bitcoin"`
	Cover    Url       `json:"cover"`
	Facebook SocialID  `json:"facebook"`
	Graph    Url       `json:"graph"`
	Location Formatted `json:"location"`
	Name     Formatted `json:"name"`
	Twitter  SocialID  `json:"twitter"`
	Version  string    `json:"v"`
	Website  string    `json:"website"`
}
