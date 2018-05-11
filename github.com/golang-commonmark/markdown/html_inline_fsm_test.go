package markdown

import "testing"

func TestMatchHTML(t *testing.T) {
	type testCase struct {
		in   string
		want string
	}
	testCases := []testCase{
		{"<!-->", ""},
		{"<!-- ->", ""},
		{"<!-- -- >", ""},
		{"<!-- -*- -->", "<!-- -*- -->"},
		{"<!#-- -->", ""},
		{"<?...??>", "<?...??>"},
		{"<?...?#>", ""},
		{"</#>", ""},
		{"", ""},
		{"</a # >", ""},
		{"</a#>", ""},
		{"<#a>", ""},
		{"<a # >", ""},
		{"<a", ""},
		{"<a#>", ""},
		{"a>", ""},
		{"a", ""},
		{"</a >", "</a >"},
		{"</a  >", "</a  >"},
		{"<a>", "<a>"},
		{"<a  >", "<a  >"},
		{"<a h#ref>", ""},
		{"<a href #>", ""},
		{"<a href= ''>", "<a href= ''>"},
		{"<a href=/ >", "<a href=/ >"},
		{"<a href = ''>", "<a href = ''>"},
		{"<a href  >", "<a href  >"},
		{"<a href=/blog>", "<a href=/blog>"},
		{"<a href='/blog'>", "<a href='/blog'>"},
		{`<a href="/blog" title="Blog">`, `<a href="/blog" title="Blog">`},
		{`<a href="http://google.com">google.com</a>`, `<a href="http://google.com">`},
		{"<a href=\x00>", ""},
		{"<a l : href>", "<a l : href>"},
		{"<a l:href>", "<a l:href>"},
		{"<a X>", "<a X>"},
		{"<br / >", ""},
		{"<br />", "<br />"},
		{"<br/>", "<br/>"},
		{"<![CDATA[...]]  >", ""},
		{"<![CDATA[...]#  >", ""},
		{"<![CDATA[ xxx xxx xxx ]]>", "<![CDATA[ xxx xxx xxx ]]>"},
		{"<!-- comment -->", "<!-- comment -->"},
		{"<!doctype html>", ""},
		{"<!Doctype html>", ""},
		{"<!DOCTYPE html>", "<!DOCTYPE html>"},
		{"<em><b", "<em>"},
		{"<em><b>", "<em><b>"},
		{"</em>", "</em>"},
		{"<em><", "<em>"},
		{`<img src="http://google.com"#/>`, ""},
		{`<img src="http://google.com"/>`, `<img src="http://google.com"/>`},
		{"<img src=image.jpeg/>", "<img src=image.jpeg/>"},
		{`<img src="image.jpg" />`, `<img src="image.jpg" />`},
		{"<img src=image\x00.jpeg/>", ""},
		{"<img src/>", "<img src/>"},
		{"<?processing-instruction?>", "<?processing-instruction?>"},
		{"<![XDATA[...]]>", ""},
		{"<!-xxx-->", ""},
	}
	for _, tc := range testCases {
		got := matchHTML(tc.in)
		if got != tc.want {
			t.Errorf("matchHTML(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
