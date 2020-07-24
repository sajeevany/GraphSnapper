package grafana

import (
	"reflect"
	"sort"
	"testing"
)

func Test_filterPanels(t *testing.T) {
	type args struct {
		panels  map[int]struct{}
		include []int
		exclude []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "test0: No inclusion or exclude list. Expect all ids to be returned ", 
			args: args{
				panels: map[int]struct{}{
					0 : struct {}{},
					1 : struct {}{},
					2 : struct {}{},
					3 : struct {}{},
					4 : struct {}{},
				},
				include: nil,
				exclude: nil,
			},
			want: []int{0, 1, 2, 3, 4},
		},
		{
			name: "test1: Inclusion list. Expect subset of map",
			args: args{
				panels: map[int]struct{}{
					0 : struct {}{},
					1 : struct {}{},
					2 : struct {}{},
					3 : struct {}{},
					4 : struct {}{},
				},
				include: []int{ 1, 2, 3},
				exclude: nil,
			},
			want: []int{1, 2, 3},
		},
		{
			name: "test2: Exclusion list. Expect subset of map",
			args: args{
				panels: map[int]struct{}{
					0 : struct {}{},
					1 : struct {}{},
					2 : struct {}{},
					3 : struct {}{},
					4 : struct {}{},
				},
				include: nil,
				exclude: []int{ 1, 2, 3},
			},
			want: []int{0, 4},
		},
		{
			name: "test3: Inclusion and exclusion list. Expect inclusion to take priority and return subset of map keys",
			args: args{
				panels: map[int]struct{}{
					0 : struct {}{},
					1 : struct {}{},
					2 : struct {}{},
					3 : struct {}{},
					4 : struct {}{},
				},
				include: []int{ 1, 2, 3},
				exclude: []int{ 1, 2, 3},
			},
			want: []int{1 , 2, 3},
		},
		{
			name: "test4: Nil panels map. Expect nil result",
			args: args{
				panels: nil,
				include: []int{ 1, 2, 3},
				exclude: []int{ 1, 2, 3},
			},
			want: nil,
		},
		{
			name: "test5: Empty panels map. Expect empty result",
			args: args{
				panels: map[int]struct{}{},
				include: []int{ 1, 2, 3},
				exclude: []int{ 1, 2, 3},
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := filterPanels(tt.args.panels, tt.args.include, tt.args.exclude)
			sort.Ints(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterPanels() = %v, want %v", got, tt.want)
			}
		})
	}
}