package rt

import (
	"encoding/json"
	"fmt"
	"time"
)

// TicketCreate fields for creating a ticket
type TicketCreate struct {
	Subject      string            `json:"Subject,omitempty"`
	Queue        string            `json:"Queue,omitempty"`
	Status       string            `json:"Status,omitempty"`
	Priority     string            `json:"Priority,omitempty"`
	Owner        string            `json:"Owner,omitempty"`
	Requestor    string            `json:"Requestor,omitempty"`
	Content      string            `json:"Content,omitempty"`
	ContentType  string            `json:"ContentType,omitempty"`
	Parent       string            `json:"Parent,omitempty"`
	CustomFields map[string]string `json:"CustomFields,omitempty"`
}

// TicketCreateResponse response for creating a ticket
type TicketCreateResponse struct {
	URL  string `json:"_url,omitempty"`
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
}

// Ticket representa un ticket en el sistema RT
type CustomField struct {
	Item
	Name   string   `json:"name,omitempty"`
	Values []string `json:"values,omitempty"`
}

type Ticket struct {
	ID              int           `json:"id,omitempty"`
	Subject         string        `json:"Subject,omitempty"`
	Queue           Queue         `json:"Queue,omitempty"`
	Status          string        `json:"Status,omitempty"`
	FinalPriority   string        `json:"FinalPriority,omitempty"`
	Owner           User          `json:"Owner,omitempty"`
	Requestor       []User        `json:"Requestor,omitempty"`
	Created         *time.Time    `json:"Created,omitempty"`
	Cc              []User        `json:"Cc,omitempty"`
	Creator         User          `json:"Creator,omitempty"`
	TimeLeft        string        `json:"TimeLeft,omitempty"`
	TimeEstimated   string        `json:"TimeEstimated,omitempty"`
	AdminCc         []User        `json:"AdminCc,omitempty"`
	Starts          *time.Time    `json:"Starts,omitempty"`
	Started         *time.Time    `json:"Started,omitempty"`
	LastUpdated     *time.Time    `json:"LastUpdated,omitempty"`
	InitialPriority string        `json:"InitialPriority,omitempty"`
	Due             *time.Time    `json:"Due,omitempty"`
	LastUpdatedBy   User          `json:"LastUpdatedBy,omitempty"`
	Priority        string        `json:"Priority,omitempty"`
	Resolved        *time.Time    `json:"Resolved,omitempty"`
	EffectiveID     Item          `json:"EffectiveID,omitempty"`
	CustomFields    []CustomField `json:"CustomFields,omitempty"`
}

// TicketUpdate representa los campos actualizables de un ticket
type TicketUpdate struct {
	Status       *string           `json:"Status,omitempty"`
	CustomFields map[string]string `json:"CustomFields,omitempty"`
}

// CreateTicket creates a ticket
func (c *Client) CreateTicket(ticket *TicketCreate) (*TicketCreateResponse, error) {
	resp, err := c.doRequest("POST", "ticket", ticket, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating ticket: %w", err)
	}

	var result TicketCreateResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &result, nil
}

// GetTicket obtiene un ticket por su ID
func (c *Client) GetTicket(id int) (*Ticket, error) {
	params := map[string]string{
		"fields[Queue]":   "Name",
		"fields[Owner]":   "Name,EmailAddress",
		"fields[Creator]": "Name,EmailAddress",
	}
	resp, err := c.doRequest("GET", fmt.Sprintf("ticket/%d", id), nil, params)
	if err != nil {
		return nil, fmt.Errorf("error getting ticket: %w", err)
	}

	var ticket Ticket
	if err := json.Unmarshal(resp, &ticket); err != nil {
		return nil, fmt.Errorf("error parsing ticket: %w", err)
	}
	// Iterate through requestors and fetch additional details
	for i := range ticket.Requestor {
		if ticket.Requestor[i].ID == "" {
			continue
		}

		user, err := c.GetUser(ticket.Requestor[i].ID)
		if err != nil {
			return nil, fmt.Errorf("error getting requestor details: %w", err)
		}
		ticket.Requestor[i].EmailAddress = user.EmailAddress
		ticket.Requestor[i].Name = user.Name
	}
	return &ticket, nil
}

// UpdateTicket actualiza un ticket existente
func (c *Client) UpdateTicket(id int, updates *TicketUpdate) error {
	_, err := c.doRequest("PUT", fmt.Sprintf("ticket/%d", id), updates, nil)
	if err != nil {
		return fmt.Errorf("error updating ticket: %w", err)
	}

	return nil
}
