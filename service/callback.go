package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type CallbackService struct {
	BackURL string
}

type ResponseData struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewCallbackService(uri string) *CallbackService {
	return &CallbackService{
		BackURL: uri,
	}
}

func (s *CallbackService) RequestTakeMoney(user_id string, amount int) error {

	payload := map[string]any{
		"action": "Debit",
		"userid": user_id,
		"amount": amount,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return errors.New("Error encoding JSON")
	}

	r, err := s.httpRequest(body)
	if err != nil {
		return err
	}

	if r.Status != 200 {
		return errors.New(r.Message)
	}
	return nil
}

func (s *CallbackService) RequestGiveMoney(user_id string, amount int) error {
	return nil
}

func (s *CallbackService) RequestRollback() {

}

func (s *CallbackService) httpRequest(body []byte) (*ResponseData, error) {

	// Data to be sent in the request body
	buff := bytes.NewBuffer(body)
	req, err := http.NewRequest("POST", s.BackURL, buff)
	if err != nil {
		return nil, errors.New("Error creating request")
	}

	// Set the Content-Type header since we're sending JSON data
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Error sending request")
	}

	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Error from server")
	}

	// Read the response body
	var responseData ResponseData
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return nil, errors.New("Error decoding JSON response:")
	}
	return &responseData, nil
}
