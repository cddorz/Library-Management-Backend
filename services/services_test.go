package services

import (
	"fmt"
	"testing"
)

func TestGetMetaDataByISBN(t *testing.T) {
	type args struct {
		isbn string
	}
	tests := []struct {
		name         string
		args         args
		wantBookInfo Book
		wantErr      bool
	}{
		// TODO: Add test cases.
		{
			name:         "chemistry",
			args:         args{isbn: "9780133363975"},
			wantBookInfo: Book{},
			wantErr:      false,
		},
		{
			name:         "real analysis",
			args:         args{isbn: "9787510040535"},
			wantBookInfo: Book{},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		Jikeapikey = "12444.6076a457ef7282751a39cc00e90ab6ab.fc8577f5db168733c16e3fd2641ab4ce"
		t.Run(tt.name, func(t *testing.T) {
			gotBookInfo, err := GetMetaDataByISBN(tt.args.isbn)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMetaDataByISBN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("%v\n", gotBookInfo.Name)
		})
	}
}
