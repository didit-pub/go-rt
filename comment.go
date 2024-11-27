package rt

import "fmt"

type Comment struct {
	Content     string `json:"Content,omitempty"`
	ContentType string `json:"ContentType,omitempty"`
}

// CommentTicket agrega un comentario a un ticket
func (c *Client) CommentTicket(ticketID int, comment *Comment) error {
	_, err := c.doRequest("POST", fmt.Sprintf("ticket/%d/comment", ticketID), comment, nil)
	if err != nil {
		return fmt.Errorf("error commenting ticket: %w", err)
	}
	return nil
}

// CorrespondTicket agrega una respuesta a un ticket
func (c *Client) CorrespondTicket(ticketID int, comment *Comment) error {
	_, err := c.doRequest("POST", fmt.Sprintf("ticket/%d/correspond", ticketID), comment, nil)
	if err != nil {
		return fmt.Errorf("error commenting ticket: %w", err)
	}
	return nil
}
