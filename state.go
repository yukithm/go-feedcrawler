package feedcrawler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"
)

// State is feed's crawling state.
type State struct {
	CrawledAt time.Time `json:"crawled_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// States is a list of State.
type States map[FeedID]*State

// UpdateState updates states by the result.
func (s States) UpdateState(result Result) {
	id := result.Subscription.ID()

	st, ok := s[id]
	if !ok {
		st = &State{}
		s[id] = st
	}

	st.CrawledAt = time.Now().Local()
	if result.Err == nil && result.Feed != nil {
		published := latestFeedTime(result.Feed)
		if published != nil {
			st.UpdatedAt = published.Local()
		}
	}
}

// LoadStatesFile loads states JSON file.
func LoadStatesFile(file string) (States, error) {
	if file == "" {
		return nil, errors.New("Invalid state file path")
	}

	states := make(States, 0)
	f, err := os.Open(file)
	if err != nil {
		if os.IsNotExist(err) {
			return states, nil
		}
		return nil, err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&states); err != nil {
		return nil, err
	}

	return states, nil
}

// SaveStatesFile save current states into JSON file.
func SaveStatesFile(states States, file string) error {
	if file == "" {
		return errors.New("Invalid state file path")
	}

	buf, err := json.MarshalIndent(states, "", "  ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(file, buf, 0666); err != nil {
		return err
	}

	return nil
}
