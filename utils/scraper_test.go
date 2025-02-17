package utils

import "testing"

func Test_searchHtmlData(t *testing.T) {
	tests := []struct {
		name        string
		str         string
		startMarker string
		endMarker   string
		want        string
		wantErr     bool
	}{
		{
			name:        "正常系：基本的なケース",
			str:         "<div>Hello World</div>",
			startMarker: "<div>",
			endMarker:   "</div>",
			want:        "Hello World",
			wantErr:     false,
		},
		{
			name:        "正常系：複数マーカーがある場合",
			str:         "<div>First</div><div>Second</div>",
			startMarker: "<div>",
			endMarker:   "</div>",
			want:        "First",
			wantErr:     false,
		},
		{
			name:        "異常系：開始マーカーが存在しない",
			str:         "Hello World</div>",
			startMarker: "<div>",
			endMarker:   "</div>",
			want:        "",
			wantErr:     true,
		},
		{
			name:        "異常系：終了マーカーが存在しない",
			str:         "<div>Hello World",
			startMarker: "<div>",
			endMarker:   "</div>",
			want:        "",
			wantErr:     true,
		},
		{
			name:        "正常系：空文字列の抽出",
			str:         "<div></div>",
			startMarker: "<div>",
			endMarker:   "</div>",
			want:        "",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := searchHtmlData(tt.str, tt.startMarker, tt.endMarker)

			if (err != nil) != tt.wantErr {
				t.Errorf("searchHtmlData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("searchHtmlData() = %v, want %v", got, tt.want)
			}
		})
	}
}
