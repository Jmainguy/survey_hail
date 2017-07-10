package main

type Messages []struct {
	Direction                string        `json:"direction"`
	From                     string        `json:"from"`
	ID                       string        `json:"id"`
	Media                    []interface{} `json:"media"`
	MessageID                string        `json:"messageId"`
	SkipMMSCarrierValidation bool          `json:"skipMMSCarrierValidation"`
	State                    string        `json:"state"`
	Text                     string        `json:"text"`
	Time                     string        `json:"time"`
	To                       string        `json:"to"`
}
