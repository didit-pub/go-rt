package rt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client contiene la lógica principal del cliente RT
type ClientConfiguration struct {
	APIURL   string
	Timeout  time.Duration `default:"30s"`
	Username string
	Password string
	Token    string
	Debug    bool `default:"false"`
}

type Client struct {
	ClientConfiguration
	httpClient *http.Client
}

// NewClient crea una nueva instancia del cliente
func NewClient(cfg *ClientConfiguration) (*Client, error) {
	if cfg.APIURL == "" {
		return nil, fmt.Errorf("RT API URL is required")
	}

	// Crear el cliente HTTP con timeout
	httpClient := &http.Client{
		Timeout: cfg.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	// Crear el cliente base
	client := &Client{
		ClientConfiguration: *cfg,
		httpClient:          httpClient,
	}

	return client, nil
}

// doRequest realiza una petición HTTP y maneja la respuesta
func (c *Client) doRequest(method, endpoint string, body interface{}, params map[string]string) ([]byte, error) {
	var err error
	// Construir URL completa
	baseURL := fmt.Sprintf("%s/%s", c.APIURL, endpoint)

	// Crear URL con parámetros
	requestURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}

	// Añadir parámetros si existen
	if len(params) > 0 {
		q := requestURL.Query()
		for key, value := range params {
			q.Add(key, value)
		}
		requestURL.RawQuery = q.Encode()
	}

	var req *http.Request
	if body != nil {
		jsonBody, errb := json.Marshal(body)
		if errb != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", errb)
		}

		// Debug logging del body si está habilitado
		if c.Debug {
			fmt.Printf("Request body:\n%s\n", string(jsonBody))
		}

		req, err = http.NewRequest(method, requestURL.String(), bytes.NewBuffer(jsonBody))
	} else {
		req, err = http.NewRequest(method, requestURL.String(), nil)
	}

	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	// Añadir headers
	req.Header.Set("Content-Type", "application/json")

	// Configurar autenticación
	if c.Token != "" {
		req.Header.Set("Authorization", "token "+c.Token)
	} else {
		req.SetBasicAuth(c.Username, c.Password)
	}

	// Debug logging
	if c.Debug {
		fmt.Printf("Request: %s %s\n", method, requestURL.String())
	}

	// Realizar la petición
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Leer la respuesta
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Verificar el código de estado
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s",
			resp.StatusCode, string(respBody))
	}
	if c.Debug {
		fmt.Printf("Response: %s\n", string(respBody))
	}
	return respBody, nil
}
