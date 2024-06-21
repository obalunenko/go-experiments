package lists

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zyedidia/generic/list"
)

func TestUsersList(t *testing.T) {
	vt := NewTrack("test")

	vt.AddSegment(Segment{ID: 1, Name: "A"})
	vt.AddSegment(Segment{ID: 2, Name: "B"})
	vt.AddSegment(Segment{ID: 3, Name: "C"})
	vt.AddSegment(Segment{ID: 4, Name: "D"})
	vt.AddSegment(Segment{ID: 5, Name: "E"})

	vt.Segments.Front.Each(func(seg Segment) {
		t.Logf("id: %d, name: %s", seg.ID, seg.Name)
	})

	b, err := json.Marshal(vt)
	require.NoError(t, err)

	t.Log(string(b))

	vt2 := NewTrack("")

	err = json.Unmarshal(b, &vt2)
	require.NoError(t, err)

	equalTrack(t, vt, vt2)
}

func equalTrack(t *testing.T, want, got *Track) {
	require.Equal(t, want.StreamID, got.StreamID)

	var gotList, wantList []Segment

	want.Segments.Front.Each(func(wantSeg Segment) {
		wantList = append(wantList, wantSeg)
	})

	got.Segments.Front.Each(func(gotSeg Segment) {
		gotList = append(gotList, gotSeg)
	})

	require.ElementsMatch(t, wantList, gotList)
}

func Test_convertListToSlice(t *testing.T) {
	type args struct {
		input []Segment
	}
	tests := []struct {
		name string
		args args
		want []Segment
	}{
		{
			name: "",
			args: args{
				input: []Segment{
					Segment{ID: 1, Name: "A"},
					Segment{ID: 2, Name: "B"},
					Segment{ID: 3, Name: "C"},
				},
			},
			want: []Segment{
				Segment{ID: 1, Name: "A"},
				Segment{ID: 2, Name: "B"},
				Segment{ID: 3, Name: "C"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := list.New[Segment]()

			for _, seg := range tt.args.input {
				l.PushFront(seg)
			}

			got := convertListToSlice(l)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}
