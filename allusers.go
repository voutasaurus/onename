package main

type AllUsersResponse struct {
	Stats     Stats              `json:"stats"`
	Usernames []string           `json:"usernames"`
	Profiles  map[string]Profile `json:"profiles"`
}

// GetAllUsers
func (c *client) GetAllUsers() ([]byte, error) {
	url := c.BaseUrl + "/users"
	return c.GetRequest(url)
}
