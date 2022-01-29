package wechat

import (
	"context"
	"reflect"
	"testing"
)

func TestGetKeyInfoByCode(t *testing.T) {
	type args struct {
		ctx  context.Context
		code string
	}
	tests := []struct {
		name    string
		args    args
		want    CodeSessionResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetKeyInfoByCode(tt.args.ctx, tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKeyInfoByCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetKeyInfoByCode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
