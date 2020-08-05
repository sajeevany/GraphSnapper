package grafana

import (
	"reflect"
	"sort"
	"testing"
)

func Test_filterPanels(t *testing.T) {
	type args struct {
		panels  map[int]PanelDescriptor
		include []int
		exclude []int
	}
	tests := []struct {
		name string
		args args
		want []PanelDescriptor
	}{
		{
			name: "test0: No inclusion or exclude list. Expect all ids to be returned ",
			args: args{
				panels: map[int]PanelDescriptor{
					0: {
						Title: "Panel0 Title",
						ID:    0,
					},
					1: {
						Title: "Panel1 Title",
						ID:    1,
					},
					2: {
						Title: "Panel2 Title",
						ID:    2,
					},
					3: {
						Title: "Panel3 Title",
						ID:    3,
					},
					4: {
						Title: "Panel4 Title",
						ID:    4,
					},
				},
				include: nil,
				exclude: nil,
			},
			want: []PanelDescriptor{
				0: {
					Title: "Panel0 Title",
					ID:    0,
				},
				1: {
					Title: "Panel1 Title",
					ID:    1,
				},
				2: {
					Title: "Panel2 Title",
					ID:    2,
				},
				3: {
					Title: "Panel3 Title",
					ID:    3,
				},
				4: {
					Title: "Panel4 Title",
					ID:    4,
				},
			},
		},
		{
			name: "test1: Inclusion list. Expect subset of map",
			args: args{
				panels: map[int]PanelDescriptor{
					0: {
						Title: "Panel0 Title",
						ID:    0,
					},
					1: {
						Title: "Panel1 Title",
						ID:    1,
					},
					2: {
						Title: "Panel2 Title",
						ID:    2,
					},
					3: {
						Title: "Panel3 Title",
						ID:    3,
					},
					4: {
						Title: "Panel4 Title",
						ID:    4,
					},
					5: {
						Title: "Panel5 Title",
						ID:    5,
					},
				},
				include: []int{1, 2, 3},
				exclude: nil,
			},
			want: func()[]PanelDescriptor{
				m := make([]PanelDescriptor, 3)
				m[0] = PanelDescriptor{
					Title: "Panel1 Title",
					ID:    1,
				}
				m[1]= PanelDescriptor{
					Title: "Panel2 Title",
					ID:    2,
				}
				m[2]= PanelDescriptor{
					Title: "Panel3 Title",
					ID:    3,
				}
				return m
			}(),
		},
		{
			name: "test2: Exclusion list. Expect subset of map",
			args: args{
				panels: map[int]PanelDescriptor{
					0: {
						Title: "Panel0 Title",
						ID:    0,
					},
					1: {
						Title: "Panel1 Title",
						ID:    1,
					},
					2: {
						Title: "Panel2 Title",
						ID:    2,
					},
					3: {
						Title: "Panel3 Title",
						ID:    3,
					},
					4: {
						Title: "Panel4 Title",
						ID:    4,
					},
				},
				include: nil,
				exclude: []int{1, 2, 3},
			},
			want: func()[]PanelDescriptor{
				m := make([]PanelDescriptor, 2)
				m[0] = PanelDescriptor{
					Title: "Panel0 Title",
					ID:    0,
				}
				m[1]= PanelDescriptor{
					Title: "Panel4 Title",
					ID:    4,
				}
				return m
			}(),
		},
		{
			name: "test3: Inclusion and exclusion list. Expect inclusion to take priority and return subset of map keys",
			args: args{
				panels: map[int]PanelDescriptor{
					0: {
						Title: "Panel0 Title",
						ID:    0,
					},
					1: {
						Title: "Panel1 Title",
						ID:    1,
					},
					2: {
						Title: "Panel2 Title",
						ID:    2,
					},
					3: {
						Title: "Panel3 Title",
						ID:    3,
					},
					4: {
						Title: "Panel4 Title",
						ID:    4,
					},
					5: {
						Title: "Panel5 Title",
						ID:    5,
					},
				},
				include: []int{1, 2, 3},
				exclude: []int{1, 2, 3},
			},
			want: func()[]PanelDescriptor{
				m := make([]PanelDescriptor, 3)
				m[0] = PanelDescriptor{
					Title: "Panel1 Title",
					ID:    1,
				}
				m[1]= PanelDescriptor{
					Title: "Panel2 Title",
					ID:    2,
				}
				m[2]= PanelDescriptor{
					Title: "Panel3 Title",
					ID:    3,
				}
				return m
			}(),
		},
		{
			name: "test4: Nil panels map. Expect nil result",
			args: args{
				panels:  nil,
				include: []int{1, 2, 3},
				exclude: []int{1, 2, 3},
			},
			want: nil,
		},
		{
			name: "test5: Empty panels map. Expect empty result",
			args: args{
				panels:  map[int]PanelDescriptor{},
				include: []int{1, 2, 3},
				exclude: []int{1, 2, 3},
			},
			want: []PanelDescriptor{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := filterPanels(tt.args.panels, tt.args.include, tt.args.exclude)
			sort.Sort(PanelDescriptors(got))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterPanels() = %v, expected panels %v", got, tt.want)
			}
		})
	}
}
