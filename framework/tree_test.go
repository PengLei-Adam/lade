package framework

import (
	"reflect"
	"testing"
)

func TestTree_AddRouter(t *testing.T) {
	type args struct {
		uri     string
		handler ControllerHandler
	}
	tree := newTree()
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				uri: "foo",
				handler: func(c *Context) error {
					return nil
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tree.AddRouter(tt.args.uri, tt.args.handler); (err != nil) != tt.wantErr {
				t.Errorf("Tree.AddRouter() error = %v, wantErr %v", err, tt.wantErr)
			}
			for i, child := range tree.root.childs {
				t.Logf("child %v segment: %v\n", i, child.segment)
			}
		})
	}
}

func TestTree_FindHandler(t *testing.T) {

	hdl1 := func(ctx *Context) error {
		t.Log("hdl1\n")
		return nil
	}
	hdl1 = ControllerHandler(hdl1)
	tests := []struct {
		name string
		uri  string
		want ControllerHandler
	}{
		// TODO: Add test cases.
		{
			name: "1",
			uri:  "/foo",
			want: hdl1,
		},
	}

	tree := newTree()
	tree.AddRouter("/foo", hdl1)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tree.FindHandler(tt.uri)

			if !reflect.DeepEqual(got, tt.want) {

				t.Errorf("Tree.FindHandler() = %v, want %v", got, tt.want)
			}

		})
	}
}
