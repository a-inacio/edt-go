package mermaid

import (
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		graph string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "with header",
			args: args{
				graph: ` stateDiagram-v2
				[*] --> B
				B --> [*]
				B --> C
				B --> B
				C --> B
				C --> D
				D --> [*]`,
			},
			want:    7,
			wantErr: false,
		},
		{
			name: "without header",
			args: args{
				graph: `
				[*] --> B
				B --> [*]
				B --> C
				B --> B
				C --> B
				C --> D
				D --> [*]`,
			},
			want:    7,
			wantErr: false,
		},
		{
			name: "with error, missing state name",
			args: args{
				graph: `
				[*] --> B
				--> [*]
				B --> C
				B --> B
				C --> B
				C --> D
				D --> [*]`,
			},
			wantErr: true,
		},
		{
			name: "with error, syntax, wrong arrow",
			args: args{
				graph: `
				[*] --> B
				B -> [*]
				B --> C
				B --> B
				C --> B
				C --> D
				D --> [*]`,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.graph)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("unexpected number of nodes, goto %d, expected %d", len(got), tt.want)
			}
		})
	}
}
