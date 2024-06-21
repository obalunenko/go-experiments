package lists

import (
	"encoding/json"
	"fmt"

	"github.com/zyedidia/generic/list"
)

type Segment struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Track struct {
	StreamID string              `json:"stream_id"`
	Segments *list.List[Segment] `json:"-"`
}

func NewTrack(sID string) *Track {
	return &Track{
		StreamID: sID,
		Segments: list.New[Segment](),
	}
}

func (t *Track) AddSegment(u Segment) {
	t.Segments.PushFront(u)
}

func (t *Track) UnmarshalJSON(bytes []byte) error {
	type Alias Track

	a := struct {
		*Alias
		Segments []Segment `json:"segments"`
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(bytes, &a); err != nil {
		return err
	}

	for i := range a.Segments {
		t.Segments.PushFront(a.Segments[i])
	}

	return nil
}

func (t *Track) MarshalJSON() ([]byte, error) {
	type Alias Track

	segments := convertListToSlice(t.Segments)

	fmt.Println("segments: ", segments)

	return json.Marshal(struct {
		*Alias
		Segments []Segment `json:"segments"`
	}{
		Alias:    (*Alias)(&Track{StreamID: t.StreamID, Segments: list.New[Segment]()}),
		Segments: segments,
	})
}

func convertListToSlice(l *list.List[Segment]) []Segment {
	var n int

	l.Front.Each(func(seg Segment) {
		n++
	})

	s := make([]Segment, 0, n)

	l.Front.Each(func(seg Segment) {
		s = append(s, seg)
	})

	return s
}
