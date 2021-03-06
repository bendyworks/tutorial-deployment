package apis

import (
	"fmt"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/ubuntu/tutorial-deployment/paths"
	"github.com/ubuntu/tutorial-deployment/testtools"

	"os"

	"github.com/ubuntu/tutorial-deployment/consts"
)

func TestNewEvents(t *testing.T) {
	testCases := []struct {
		eventsDir string

		wantEvents Events
		wantErr    bool
	}{
		{"testdata/events/valid",
			Events{"event-1": event{Name: "Event 1", Logo: "img/event1.jpg", Description: "This workshop is taking place at Event 1."},
				"event-2": event{Name: "Event 2", Logo: "event2.jpg", Description: "This workshop is taking place at Event 2."},
			},
			false},
		{"doesnt/exist", nil, true},
		{"testdata/events/valid-missing-image", // we still load correctly, we don't touch images at this stage
			Events{"event-1": event{Name: "Event 1", Logo: "img/event1.jpg", Description: "This workshop is taking place at Event 1."},
				"event-2": event{Name: "Event 2", Logo: "event2.jpg", Description: "This workshop is taking place at Event 2."},
			},
			false},
		{"testdata/events/no-events", Events{}, false},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("create events for: %+v", tc.eventsDir), func(t *testing.T) {
			// Setup/Teardown
			p, teardown := paths.MockPath()
			defer teardown()
			p.MetaData = tc.eventsDir

			// Test
			e, err := NewEvents()

			if (err != nil) != tc.wantErr {
				t.Errorf("NewEvents() error = %v, wantErr %v", err, tc.wantErr)
			}
			if err != nil {
				return
			}

			if !reflect.DeepEqual(*e, tc.wantEvents) {
				t.Errorf("Generated events: got %+v; want %+v", *e, tc.wantEvents)
			}
		})
	}
}

func TestSaveImages(t *testing.T) {
	testCases := []struct {
		eventsDir string
		eventsObj Events

		wantEvents Events
		wantErr    bool
	}{
		{"testdata/events/valid",
			Events{"event-1": event{Name: "Event 1", Logo: "img/event1.jpg", Description: "This workshop is taking place at Event 1."},
				"event-2": event{Name: "Event 2", Logo: "event2.jpg", Description: "This workshop is taking place at Event 2."},
			},
			Events{"event-1": event{Name: "Event 1", Logo: fmt.Sprintf("%sevent1.jpg", consts.ImagesURL), Description: "This workshop is taking place at Event 1."},
				"event-2": event{Name: "Event 2", Logo: fmt.Sprintf("%sevent2.jpg", consts.ImagesURL), Description: "This workshop is taking place at Event 2."},
			},
			false},
		{"testdata/events/valid-missing-image",
			Events{"event-1": event{Name: "Event 1", Logo: "img/event1.jpg", Description: "This workshop is taking place at Event 1."},
				"event-2": event{Name: "Event 2", Logo: "event2.jpg", Description: "This workshop is taking place at Event 2."},
			},
			nil,
			true},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("save events: %+v", tc.eventsDir), func(t *testing.T) {
			// Setup/Teardown
			apiout, teardown := testtools.TempDir(t)
			defer teardown()
			imagesout, teardown := testtools.TempDir(t)
			defer teardown()
			p, teardown := paths.MockPath()
			defer teardown()
			p.MetaData = tc.eventsDir
			p.API = apiout
			p.Images = imagesout

			// Test
			err := tc.eventsObj.SaveImages()

			if (err != nil) != tc.wantErr {
				t.Errorf("SaveImages() error = %v, wantErr %v", err, tc.wantErr)
			}
			if err != nil {
				return
			}

			if !reflect.DeepEqual(tc.eventsObj, tc.wantEvents) {
				t.Errorf("Image paths not correctly changed in event: got %+v; want %+v", tc.eventsObj, tc.wantEvents)
			}
			for _, e := range tc.wantEvents {
				imgP := path.Join(p.Images, strings.TrimPrefix(e.Logo, consts.ImagesURL))
				if _, err := os.Stat(imgP); os.IsNotExist(err) {
					t.Errorf("%s doesn't exist when we wanted it", imgP)
				}
			}
		})
	}
}
