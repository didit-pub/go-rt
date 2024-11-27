package rt

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type User struct {
	ID           string `json:"id,omitempty"`
	URL          string `json:"_url,omitempty"`
	Name         string `json:"name,omitempty"`
	EmailAddress string `json:"EmailAddress,omitempty"`
}

func (c *Client) GetUser(userName string) (User, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("user/%s", userName), nil, nil)
	if err != nil {
		return User{}, fmt.Errorf("error getting ticket: %w", err)
	}

	var user User
	if err := json.Unmarshal(resp, &user); err != nil {
		return User{}, fmt.Errorf("error parsing user: %w", err)
	}

	return user, nil
}

func (u *User) UnmarshalJSON(data []byte) error {
	// Creamos una estructura temporal que usaremos para el unmarshaling inicial
	type UserAlias User
	type UserRaw struct {
		*UserAlias
		ID interface{} `json:"id,omitempty"`
	}

	raw := &UserRaw{UserAlias: (*UserAlias)(u)}
	if err := json.Unmarshal(data, raw); err != nil {
		return err
	}

	// Convertimos ID a string según el tipo que venga
	switch v := raw.ID.(type) {
	case string:
		u.ID = v
	case float64: // JSON numbers se decodifican como float64
		u.ID = strconv.FormatInt(int64(v), 10)
	case nil:
		u.ID = ""
	default:
		return fmt.Errorf("invalid type for ID: %T", v)
	}

	return nil
}
