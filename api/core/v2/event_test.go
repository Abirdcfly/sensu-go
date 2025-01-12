package v2

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFixtureEventIsValid(t *testing.T) {
	e := FixtureEvent("entity", "check")
	assert.NotNil(t, e)
	assert.NotNil(t, e.Entity)
	assert.NotNil(t, e.Check)
}

func TestEventValidate(t *testing.T) {
	event := FixtureEvent("entity", "check")

	event.Check.Name = ""
	assert.Error(t, event.Validate())
	event.Check.Name = "check"

	event.Entity.Name = ""
	assert.Error(t, event.Validate())
	event.Entity.Name = "entity"

	assert.NoError(t, event.Validate())
}

func TestEventValidateNoTimestamp(t *testing.T) {
	// Events without a timestamp are valid
	event := FixtureEvent("entity", "check")
	event.Timestamp = 0
	if err := event.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestMarshalJSON(t *testing.T) {
	event := FixtureEvent("entity", "check")
	_, err := json.Marshal(event)
	require.NoError(t, err)
}

func TestEventHasMetrics(t *testing.T) {
	testCases := []struct {
		name     string
		metrics  *Metrics
		expected bool
	}{
		{
			name:     "No Metrics",
			metrics:  nil,
			expected: false,
		},
		{
			name:     "Metrics",
			metrics:  &Metrics{},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event := &Event{
				Metrics: tc.metrics,
			}
			metrics := event.HasMetrics()
			assert.Equal(t, tc.expected, metrics)
		})
	}
}

func TestEventIsIncident(t *testing.T) {
	testCases := []struct {
		name     string
		status   uint32
		expected bool
	}{
		{
			name:     "OK Status",
			status:   0,
			expected: false,
		},
		{
			name:     "Non-zero Status",
			status:   1,
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event := &Event{
				Check: &Check{
					Status: tc.status,
				},
			}
			incident := event.IsIncident()
			assert.Equal(t, tc.expected, incident)
		})
	}
}

func TestEventIsResolution(t *testing.T) {
	testCases := []struct {
		name     string
		history  []CheckHistory
		status   uint32
		expected bool
	}{
		{
			name:     "check has no history",
			history:  []CheckHistory{{}},
			status:   0,
			expected: false,
		},
		{
			name: "check has not transitioned",
			history: []CheckHistory{
				{Status: 1},
				{Status: 0},
			},
			status:   0,
			expected: true,
		},
		{
			name: "check has just transitioned",
			history: []CheckHistory{
				{Status: 0},
				{Status: 1},
			},
			status:   0,
			expected: false,
		},
		{
			name: "check has transitioned but still an incident",
			history: []CheckHistory{
				{Status: 2},
				{Status: 1},
			},
			status:   1,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event := &Event{
				Check: &Check{
					History: tc.history,
					Status:  tc.status,
				},
			}
			resolution := event.IsResolution()
			assert.Equal(t, tc.expected, resolution)
		})
	}
}

func TestEventIsSilenced(t *testing.T) {
	testCases := []struct {
		name     string
		event    *Event
		silenced []string
		expected bool
	}{
		{
			name:     "No silenced entries",
			event:    FixtureEvent("entity1", "check1"),
			silenced: []string{},
			expected: false,
		},
		{
			name:     "Silenced entry",
			event:    FixtureEvent("entity1", "check1"),
			silenced: []string{"entity1"},
			expected: true,
		},
		{
			name:     "Metric without a check",
			event:    &Event{},
			silenced: []string{"entity1"},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.event.Check != nil {
				tc.event.Check.Silenced = tc.silenced
			}
			silenced := tc.event.IsSilenced()
			assert.Equal(t, tc.expected, silenced)
		})
	}
}

func TestEventIsFlappingStart(t *testing.T) {
	testCases := []struct {
		name     string
		history  []CheckHistory
		state    string
		expected bool
	}{
		{
			name:     "check has no history",
			history:  []CheckHistory{{}},
			state:    EventPassingState,
			expected: false,
		},
		{
			name: "check was not flapping previously, nor is now",
			history: []CheckHistory{
				{Flapping: false},
				{Flapping: false},
			},
			state:    EventPassingState,
			expected: false,
		},
		{
			name: "check is already flapping",
			history: []CheckHistory{
				{Flapping: true},
				{Flapping: true},
			},
			state:    EventFlappingState,
			expected: false,
		},
		{
			name: "check was not previously flapping",
			history: []CheckHistory{
				{Flapping: false},
				{Flapping: true},
			},
			state:    EventFlappingState,
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event := &Event{
				Check: &Check{
					History: tc.history,
					State:   tc.state,
				},
			}
			assert.Equal(t, tc.expected, event.IsFlappingStart())
		})
	}
}

func TestEventIsFlappingEnd(t *testing.T) {
	testCases := []struct {
		name     string
		history  []CheckHistory
		state    string
		expected bool
	}{
		{
			name:     "check has no history",
			history:  []CheckHistory{{}},
			state:    EventPassingState,
			expected: false,
		},
		{
			name: "check was not flapping previously, nor is now",
			history: []CheckHistory{
				{Flapping: false},
				{Flapping: false},
			},
			state:    EventPassingState,
			expected: false,
		},
		{
			name: "check is already flapping",
			history: []CheckHistory{
				{Flapping: true},
				{Flapping: true},
			},
			state:    EventFlappingState,
			expected: false,
		},
		{
			name: "check was previously flapping but now is in OK state",
			history: []CheckHistory{
				{Flapping: true},
				{Flapping: false},
			},
			state:    EventPassingState,
			expected: true,
		},
		{
			name: "check was previously flapping but now is in failing state",
			history: []CheckHistory{
				{Flapping: true},
				{Flapping: false},
			},
			state:    EventFailingState,
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event := &Event{
				Check: &Check{
					History: tc.history,
					State:   tc.state,
				},
			}
			assert.Equal(t, tc.expected, event.IsFlappingEnd())
		})
	}
}

func TestEventIsSilencedBy(t *testing.T) {
	testCases := []struct {
		name     string
		event    *Event
		silenced *Silenced
		expected bool
	}{
		{
			name:     "nil silenced",
			event:    FixtureEvent("entity1", "check1"),
			silenced: nil,
			expected: false,
		},
		{
			name:     "Metric without a check",
			event:    &Event{},
			silenced: nil,
			expected: false,
		},
		{
			name:     "Check provided and doesn't match, no subscription",
			event:    FixtureEvent("entity1", "check1"),
			silenced: FixtureSilenced("*:check2"),
			expected: false,
		},
		{
			name:  "Check provided and matches, no subscription",
			event: FixtureEvent("entity1", "check1"),
			silenced: &Silenced{
				Subscription: "",
				Check:        "check1",
			},
			expected: true,
		},
		{
			name:  "Check provided and matches, no subscription; begins in future",
			event: FixtureEvent("entity1", "check1"),
			silenced: &Silenced{
				Subscription: "",
				Check:        "check1",
				Begin:        time.Now().Unix() + 300,
			},
			expected: false,
		},
		{
			name:     "Subscription provided and doesn't match, no check",
			event:    FixtureEvent("entity1", "check1"),
			silenced: FixtureSilenced("entity:entity2:*"),
			expected: false,
		},
		{
			name:     "Subscription provided and matches, no check",
			event:    FixtureEvent("entity1", "check1"),
			silenced: FixtureSilenced("entity:entity1:*"),
			expected: true,
		},
		{
			name:     "Check provided and doesn't match, subscription matches",
			event:    FixtureEvent("entity1", "check1"),
			silenced: FixtureSilenced("entity:entity1:check2"),
			expected: false,
		},
		{
			name:     "Check provided and matches, subscription doesn't match",
			event:    FixtureEvent("entity1", "check1"),
			silenced: FixtureSilenced("entity:entity2:check1"),
			expected: false,
		},
		{
			name:     "Check and subscription both provided and match",
			event:    FixtureEvent("entity1", "check1"),
			silenced: FixtureSilenced("entity:entity1:check1"),
			expected: true,
		},
		{
			name:     "Subscription provided and doesn't match, check matches",
			event:    FixtureEvent("entity1", "check1"),
			silenced: FixtureSilenced("entity:entity2:check1"),
			expected: false,
		},
		{
			name:     "Subscription provided and matches, check doesn't match",
			event:    FixtureEvent("entity1", "check1"),
			silenced: FixtureSilenced("entity:entity1:check2"),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.event.IsSilencedBy(tc.silenced))
		})
	}
}

func TestSilencedBy(t *testing.T) {
	testCases := []struct {
		name            string
		event           *Event
		entries         []*Silenced
		expectedEntries []*Silenced
	}{
		{
			name:            "no entries",
			event:           FixtureEvent("foo", "check_cpu"),
			entries:         []*Silenced{},
			expectedEntries: []*Silenced{},
		},
		{
			name:  "not silenced",
			event: FixtureEvent("foo", "check_cpu"),
			entries: []*Silenced{
				FixtureSilenced("entity:foo:check_mem"),
				FixtureSilenced("entity:bar:*"),
				FixtureSilenced("foo:check_cpu"),
				FixtureSilenced("foo:*"),
				FixtureSilenced("*:check_mem"),
			},
			expectedEntries: []*Silenced{},
		},
		{
			name:  "silenced by check",
			event: FixtureEvent("foo", "check_cpu"),
			entries: []*Silenced{
				FixtureSilenced("*:check_cpu"),
			},
			expectedEntries: []*Silenced{
				FixtureSilenced("*:check_cpu"),
			},
		},
		{
			name:  "silenced by entity subscription",
			event: FixtureEvent("foo", "check_cpu"),
			entries: []*Silenced{
				FixtureSilenced("entity:foo:*"),
			},
			expectedEntries: []*Silenced{
				FixtureSilenced("entity:foo:*"),
			},
		},
		{
			name:  "silenced by entity's check subscription",
			event: FixtureEvent("foo", "check_cpu"),
			entries: []*Silenced{
				FixtureSilenced("entity:foo:check_cpu"),
			},
			expectedEntries: []*Silenced{
				FixtureSilenced("entity:foo:check_cpu"),
			},
		},
		{
			name:  "silenced by check subscription",
			event: FixtureEvent("foo", "check_cpu"), // has a linux subscription
			entries: []*Silenced{
				FixtureSilenced("linux:*"),
			},
			expectedEntries: []*Silenced{
				FixtureSilenced("linux:*"),
			},
		},
		{
			name:  "silenced by subscription with check",
			event: FixtureEvent("foo", "check_cpu"), // has a linux subscription
			entries: []*Silenced{
				FixtureSilenced("linux:check_cpu"),
			},
			expectedEntries: []*Silenced{
				FixtureSilenced("linux:check_cpu"),
			},
		},
		{
			name:  "silenced by multiple entries",
			event: FixtureEvent("foo", "check_cpu"), // has a linux subscription
			entries: []*Silenced{
				FixtureSilenced("entity:foo:*"),
				FixtureSilenced("linux:check_cpu"),
			},
			expectedEntries: []*Silenced{
				FixtureSilenced("entity:foo:*"),
				FixtureSilenced("linux:check_cpu"),
			},
		},
		{
			name: "not silenced, silenced & client don't have a common subscription",
			event: &Event{
				Check: &Check{
					ObjectMeta: ObjectMeta{
						Name: "check_cpu",
					},
					Subscriptions: []string{"linux", "windows"},
				},
				Entity: &Entity{
					ObjectMeta: ObjectMeta{
						Name: "foo",
					},
					Subscriptions: []string{"linux"},
				},
			},
			entries: []*Silenced{
				FixtureSilenced("windows:check_cpu"),
			},
			expectedEntries: []*Silenced{},
		},
		{
			name: "silenced, silenced & client do have a common subscription",
			event: &Event{
				Check: &Check{
					ObjectMeta: ObjectMeta{
						Name: "check_cpu",
					},
					Subscriptions: []string{"linux", "windows"},
				},
				Entity: &Entity{
					ObjectMeta: ObjectMeta{
						Name: "foo",
					},
					Subscriptions: []string{"linux"},
				},
			},
			entries: []*Silenced{
				FixtureSilenced("linux:check_cpu"),
			},
			expectedEntries: []*Silenced{
				FixtureSilenced("linux:check_cpu"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.event.SilencedBy(tc.entries)
			assert.EqualValues(t, tc.expectedEntries, result)
		})
	}
}

func TestIsSilencedBy(t *testing.T) {
	testCases := []struct {
		name           string
		event          *Event
		silence        *Silenced
		expectedResult bool
	}{
		{
			name:  "silence has not started",
			event: FixtureEvent("foo", "check_cpu"),
			silence: &Silenced{
				ObjectMeta: ObjectMeta{
					Name: "*:check_cpu",
				},
				Begin: time.Now().Add(1 * time.Hour).Unix(),
			},
			expectedResult: false,
		},
		{
			name:           "check matches w/ wildcard subscription",
			event:          FixtureEvent("foo", "check_cpu"),
			silence:        FixtureSilenced("*:check_cpu"),
			expectedResult: true,
		},
		{
			name:           "entity subscription matches w/ wildcard check",
			event:          FixtureEvent("foo", "check_cpu"),
			silence:        FixtureSilenced("entity:foo:*"),
			expectedResult: true,
		},
		{
			name:           "entity subscription and check match",
			event:          FixtureEvent("foo", "check_cpu"),
			silence:        FixtureSilenced("entity:foo:check_cpu"),
			expectedResult: true,
		},
		{
			name: "subscription matches",
			event: &Event{
				Check: &Check{
					ObjectMeta: ObjectMeta{
						Name: "check_cpu",
					},
					Subscriptions: []string{"unix"},
				},
				Entity: &Entity{
					ObjectMeta: ObjectMeta{
						Name: "foo",
					},
					Subscriptions: []string{"unix"},
				},
			},
			silence:        FixtureSilenced("unix:check_cpu"),
			expectedResult: true,
		},
		{
			name: "subscription does not match",
			event: &Event{
				Check: &Check{
					ObjectMeta: ObjectMeta{
						Name: "check_cpu",
					},
					Subscriptions: []string{"unix"},
				},
				Entity: &Entity{
					ObjectMeta: ObjectMeta{
						Name: "foo",
					},
					Subscriptions: []string{"unix"},
				},
			},
			silence:        FixtureSilenced("windows:check_cpu"),
			expectedResult: false,
		},
		{
			name: "entity subscription doesn't match",
			event: &Event{
				Check: &Check{
					ObjectMeta: ObjectMeta{
						Name: "check_cpu",
					},
					Subscriptions: []string{"unix"},
				},
				Entity: &Entity{
					ObjectMeta: ObjectMeta{
						Name: "foo",
					},
					Subscriptions: []string{"windows"},
				},
			},
			silence: &Silenced{
				Subscription: "check",
				Check:        "check_cpu",
			},
			expectedResult: false,
		},
		{
			name:           "check does not match",
			event:          FixtureEvent("foo", "check_mem"),
			silence:        FixtureSilenced("*:check_cpu"),
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.event.IsSilencedBy(tc.silence)
			assert.EqualValues(t, tc.expectedResult, result)
		})
	}
}

func fixtureNoID() *Event {
	e := FixtureEvent("foo", "bar")
	e.ID = nil
	return e
}

func fixtureBadID() *Event {
	e := FixtureEvent("foo", "bar")
	e.ID = []byte("not a uuid")
	return e
}

func TestMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		Name         string
		Event        *Event
		MarshalError bool
	}{
		{
			Name:  "event with no ID",
			Event: fixtureNoID(),
		},
		{
			Name:  "event with ID",
			Event: FixtureEvent("foo", "bar"),
		},
		{
			Name:         "event with invalid ID",
			Event:        fixtureBadID(),
			MarshalError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			b, err := json.Marshal(test.Event)
			if test.MarshalError {
				if err == nil {
					t.Fatal("expected non-nil error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			var e Event
			if err := json.Unmarshal(b, &e); err != nil {
				t.Fatal(err)
			}
			if err := e.Validate(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUnmarshalID(t *testing.T) {
	tests := []struct {
		Name           string
		Data           string
		UnmarshalError bool
	}{
		{
			Name: "no id",
			Data: `{}`,
		},
		{
			Name: "has id",
			Data: fmt.Sprintf(`{"id": %q}`, uuid.NameSpaceDNS),
		},
		{
			Name:           "invalid id",
			Data:           `{"id": "not a uuid"}`,
			UnmarshalError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var e Event
			err := json.Unmarshal([]byte(test.Data), &e)
			if test.UnmarshalError {
				if err == nil {
					t.Fatal("expected non-nil error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestEventFields(t *testing.T) {
	tests := []struct {
		name    string
		args    Resource
		wantKey string
		want    string
	}{
		{
			name:    "exposes entity.name",
			args:    FixtureEvent("frank", "reynolds"),
			wantKey: "event.entity.name",
			want:    "frank",
		},
		{
			name:    "exposes check.name",
			args:    FixtureEvent("frank", "reynolds"),
			wantKey: "event.check.name",
			want:    "reynolds",
		},
		{
			name:    "exposes check.state",
			args:    FixtureEvent("frank", "reynolds"),
			wantKey: "event.check.state",
			want:    "passing",
		},
		{
			name: "exposes check labels",
			args: &Event{
				Check:  &Check{ObjectMeta: ObjectMeta{Labels: map[string]string{"src": "bonsai"}}},
				Entity: &Entity{},
			},
			wantKey: "event.labels.src",
			want:    "bonsai",
		},
		{
			name: "exposes entity labels",
			args: &Event{
				Check:  &Check{},
				Entity: &Entity{ObjectMeta: ObjectMeta{Labels: map[string]string{"region": "philadelphia"}}},
			},
			wantKey: "event.labels.region",
			want:    "philadelphia",
		},
		{
			name: "check labels take precendence",
			args: &Event{
				Check:  &Check{ObjectMeta: ObjectMeta{Labels: map[string]string{"dupe": "check"}}},
				Entity: &Entity{ObjectMeta: ObjectMeta{Labels: map[string]string{"dupe": "entity"}}},
			},
			wantKey: "event.labels.dupe",
			want:    "check",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EventFields(tt.args)
			if !reflect.DeepEqual(got[tt.wantKey], tt.want) {
				t.Errorf("EventFields() = got[%s] %v, want[%s] %v", tt.wantKey, got[tt.wantKey], tt.wantKey, tt.want)
			}
		})
	}
}

func TestSetNamespace(t *testing.T) {
	event := new(Event)
	event.SetNamespace("foobar")
	if event.Entity == nil {
		t.Fatal("nil entity")
	}
	if got, want := event.Namespace, "foobar"; got != want {
		t.Errorf("bad namespace: got %q, want %q", got, want)
	}
	if got, want := event.Entity.Namespace, "foobar"; got != want {
		t.Errorf("bad namespace: got %q, want %q", got, want)
	}
	if event.Check != nil {
		t.Fatal("check should have been nil")
	}
	event.Check = new(Check)
	event.SetNamespace("foobar")
	if got, want := event.Check.Namespace, "foobar"; got != want {
		t.Errorf("bad namespace: got %q, want %q", got, want)
	}
}

func TestEventURIPath(t *testing.T) {
	e := new(Event)
	if got, want := e.URIPath(), "/api/core/v2/events"; got != want {
		t.Errorf("bad URIPath; got %q, want %q", got, want)
	}
	e.SetNamespace("foobar")
	if got, want := e.URIPath(), "/api/core/v2/namespaces/foobar/events"; got != want {
		t.Errorf("bad URIPath; got %q, want %q", got, want)
	}
	e.Entity.Name = "baz"
	if got, want := e.URIPath(), "/api/core/v2/namespaces/foobar/events/baz"; got != want {
		t.Errorf("bad URIPath; got %q, want %q", got, want)
	}
	e.Check = new(Check)
	e.Check.Name = "bep"
	if got, want := e.URIPath(), "/api/core/v2/namespaces/foobar/events/baz/bep"; got != want {
		t.Errorf("bad URIPath; got %q, want %q", got, want)
	}
}
