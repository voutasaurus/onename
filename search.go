package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

type SearchResponse struct {
	OneNameErrorResponse
	Results []struct {
		Profile  Profile `json:"profile"`
		Username string  `json:"username"`
	} `json:"results"`
}

func (s SearchResponse) String() string {
	if s.Error != nil {
		return s.Error.String()
	}
	return fmt.Sprint(s.Results)
}

type Profile struct {
	Avatar  Url `json:"avatar"`
	Bio     string
	Bitcoin struct {
		Address string `json:"address"`
	} `json:"bitcoin"`
	Cover    Url       `json:"cover"`
	Facebook SocialID  `json:"facebook,omitempty"`
	Github   SocialID  `json:"github,omitempty"`
	Graph    Url       `json:"graph"`
	Location Formatted `json:"location"`
	Name     Formatted `json:"name"`
	Twitter  SocialID  `json:"twitter,omitempty"`
	Version  string    `json:"v"`
	Website  string    `json:"website"`
}

func (p Profile) String() string {
	ret := "\n\tProfile:\n"
	ret += "\n"
	ret += "\tName: " + p.Name.Formatted + "\n"
	ret += "\tLocation: " + p.Location.Formatted + "\n"
	ret += "\tBio: " + p.Bio + "\n"
	ret += "\tBitcoin Address: " + p.Bitcoin.Address + "\n"
	ret += "\tFacebook: " + p.Facebook.String() + "\n"
	ret += "\tTwitter: " + p.Twitter.String() + "\n"
	ret += "\tGithub: " + p.Github.String() + "\n"
	ret += "\tGraph: " + p.Graph.Url + "\n"
	ret += "\tAvatar: " + p.Avatar.Url + "\n"
	ret += "\tCover: " + p.Cover.Url + "\n"
	ret += "\tVersion: " + p.Version + "\n"
	ret += "\tWebsite: " + p.Website + "\n"
	return ret
}

type SocialID struct {
	Proof    Url    `json:"proof"`
	Username string `json:"username"`
}

func (s SocialID) String() string {
	ret := "\n"
	ret += "\t\tSocialID: " + s.Username + "\n"
	ret += "\t\tProof: " + fmt.Sprint(s.Proof)
	return ret
}

type Url struct {
	Url string `json:"url"`
}

type Formatted struct {
	Formatted string `json:"formatted"`
}

// GetProfileObjects parses the JSON responses matching the SearchResponse type
func (c *client) GetProfileObjects(url string) (SearchResponse, error) {
	var jsonObj SearchResponse
	jsonBytes, err := c.GetRequest(url)
	if err != nil {
		return jsonObj, err
	}
	err = json.Unmarshal(jsonBytes, &jsonObj)
	if err != nil {
		return jsonObj, JSONRead(err, jsonBytes)
	}
	if jsonObj.Error != nil {
		err = errors.New("Error: " + jsonObj.Error.Type + " - " + jsonObj.Error.Message)
	}
	return jsonObj, err
}

func (c *client) SearchUsers(query string) (SearchResponse, error) {
	url := c.BaseUrl + "/search?query=" + query
	return c.GetProfileObjects(url)
}
