package feature

import (
	"errors"
	"time"
)

var (
	ErrEventFeatureIDRequired = errors.New("event feature id is required")
	ErrEventTypeRequired      = errors.New("event type is required")
)

type Event struct {
	ID        int64
	FeatureID int32
	UserID    string
	Variant   string
	Type      string
	CreatedAt time.Time
}

func NewEvent(featureID int32, userID, variant, eventType string) (*Event, error) {
	event := &Event{
		FeatureID: featureID,
		UserID:    userID,
		Variant:   variant,
		Type:      eventType,
	}

	return event, event.Validate()
}

func (e *Event) Validate() error {
	var errs []error

	if e.FeatureID <= 0 {
		errs = append(errs, ErrEventFeatureIDRequired)
	}

	if e.Type == "" {
		errs = append(errs, ErrEventTypeRequired)
	}

	return errors.Join(errs...)
}
