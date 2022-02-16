package framework

import (
	"net/http"
	"reflect"
	"testing"
)

func TestCore_FindRouteByRequest(t *testing.T) {
	type args struct {
		request *http.Request
	}
	tests := []struct {
		name string
		c    *Core
		args args
		want ControllerHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.FindRouteByRequest(tt.args.request); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Core.FindRouteByRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
