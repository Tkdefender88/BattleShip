package solver

import "testing"

func Test_indexFromString(t *testing.T) {
	type args struct {
		target []rune
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Simple",
			args: args{
				[]rune{'A', '0'},
			},
			want: 0,
		},
		{
			name: "last tile",
			args: args{
				[]rune{'J', '9'},
			},
			want: 99,
		},
		{
			name: "middleish",
			args: args{
				[]rune{'E', '5'},
			},
			want: 45,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := indexFromString(tt.args.target); got != tt.want {
				t.Errorf("indexFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}
