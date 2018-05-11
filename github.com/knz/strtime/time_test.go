package strtime

import (
	"testing"
	"time"
)

func TestTimeConversion(t *testing.T) {
	tests := []struct {
		start     string
		format    string
		tm        string
		revformat string
		reverse   string
	}{
		// %a %A %b %B (+ %Y)
		{`Wed Oct 05 2016`, `%a %b %d %Y`, `2016-10-05T00:00:00Z`, ``, ``},
		{`Wednesday October 05 2016`, `%A %B %d %Y`, `2016-10-05T00:00:00Z`, ``, ``},
		// %c
		{`Wed Oct 5 01:02:03 2016`, `%c`, `2016-10-05T01:02:03Z`, ``, ``},
		// %C %d (+ %m %y)
		{`20 06 10 12`, `%C %y %m %d`, `2006-10-12T00:00:00Z`, ``, ``},
		// %D
		{`10/12/06`, `%D`, `2006-10-12T00:00:00Z`, ``, ``},
		// %e (+ %Y %m)
		{`2006 10  3`, `%Y %m %e`, `2006-10-03T00:00:00Z`, ``, ``},
		// %f (+ %c)
		{`Wed Oct 5 01:02:03 2016 .123`, `%c .%f`, `2016-10-05T01:02:03.123Z`, `.%f`, `.123000000`},
		{`Wed Oct 5 01:02:03 2016 .123456`, `%c .%f`, `2016-10-05T01:02:03.123456Z`, `.%f`, `.123456000`},
		{`Wed Oct 5 01:02:03 2016 .123456789`, `%c .%f`, `2016-10-05T01:02:03.123456789Z`, `.%f`, `.123456789`},
		{`Wed Oct 5 01:02:03 2016 .999999999`, `%c .%f`, `2016-10-05T01:02:03.999999999Z`, `.%f`, `.999999999`},
		// %F
		{`2006-10-03`, `%F`, `2006-10-03T00:00:00Z`, ``, ``},
		// %h (+ %Y %d)
		{`2006 Oct 03`, `%Y %h %d`, `2006-10-03T00:00:00Z`, ``, ``},
		// %H (+ %S %M)
		{`20061012 01:03:02`, `%Y%m%d %H:%S:%M`, `2006-10-12T01:02:03Z`, ``, ``},
		// %I (+ %Y %m %d)
		{`20161012 11`, `%Y%m%d %I`, `2016-10-12T11:00:00Z`, ``, ``},
		// %j (+ %Y)
		{`2016 286`, `%Y %j`, `2016-10-12T00:00:00Z`, ``, ``},
		// %k (+ %Y %m %d)
		{`20061012 23`, `%Y%m%d %k`, `2006-10-12T23:00:00Z`, ``, ``},
		// %l (+ %Y %m %d %p)
		{`20061012  5 PM`, `%Y%m%d %l %p`, `2006-10-12T17:00:00Z`, ``, ``},
		// %n (+ %Y %m %d)
		{"2006\n10\n03", `%Y%n%m%n%d`, `2006-10-03T00:00:00Z`, ``, ``},
		// %p cannot be parsed before hour specifiers, so be sure that
		// they appear in this order.
		{`20161012 11 PM`, `%Y%m%d %I %p`, `2016-10-12T23:00:00Z`, ``, ``},
		{`20161012 11 AM`, `%Y%m%d %I %p`, `2016-10-12T11:00:00Z`, ``, ``},
		// %r
		{`20161012 11:02:03 PM`, `%Y%m%d %r`, `2016-10-12T23:02:03Z`, ``, ``},
		// %R
		{`20161012 11:02`, `%Y%m%d %R`, `2016-10-12T11:02:00Z`, ``, ``},
		// %s
		{`1491920586`, `%s`, `2017-04-11T14:23:06Z`, ``, ``},
		// %t (+ %Y %m %d)
		{"2006\t10\t03", `%Y%t%m%t%d`, `2006-10-03T00:00:00Z`, ``, ``},
		// %T (+ %Y %m %d)
		{`20061012 01:02:03`, `%Y%m%d %T`, `2006-10-12T01:02:03Z`, ``, ``},
		// %U %u (+ %Y)
		{`2018 10 4`, `%Y %U %u`, `2018-03-15T00:00:00Z`, ``, ``},
		// %W %w (+ %Y)
		{`2018 10 4`, `%Y %W %w`, `2018-03-08T00:00:00Z`, ``, ``},
		// %x
		{`10/12/06`, `%x`, `2006-10-12T00:00:00Z`, ``, ``},
		// %X
		{`20061012 01:02:03`, `%Y%m%d %X`, `2006-10-12T01:02:03Z`, ``, ``},
		// %y (+ %m %d)
		{`000101`, `%y%m%d`, `2000-01-01T00:00:00Z`, ``, ``},
		{`680101`, `%y%m%d`, `2068-01-01T00:00:00Z`, ``, ``},
		{`690101`, `%y%m%d`, `1969-01-01T00:00:00Z`, ``, ``},
		{`990101`, `%y%m%d`, `1999-01-01T00:00:00Z`, ``, ``},
		// %Y
		{`19000101`, `%Y%m%d`, `1900-01-01T00:00:00Z`, ``, ``},
		{`20000101`, `%Y%m%d`, `2000-01-01T00:00:00Z`, ``, ``},
		{`30000101`, `%Y%m%d`, `3000-01-01T00:00:00Z`, ``, ``},
		// %z causes the time zone to adjust the time when parsing, but the time zone information
		// is not retained when printing the timestamp out back.
		{`20160101 13:00 +0655`, `%Y%m%d %H:%M %z`, `2016-01-01T06:05:00Z`, `%Y%m%d %H:%M %z`, `20160101 06:05 +0000`},
	}

	for _, test := range tests {
		tm, err := Strptime(test.start, test.format)
		if err != nil {
			t.Errorf("strptime(%q, %q): %v", test.format, test.start, err)
			continue
		}
		tm = tm.UTC()

		tmS := tm.Format(time.RFC3339Nano)
		if tmS != test.tm {
			t.Errorf("strptime(%q, %q): got %q, expected %q", test.start, test.format, tmS, test.tm)
			continue
		}

		revfmt := test.format
		if test.revformat != "" {
			revfmt = test.revformat
		}

		ref := test.start
		if test.reverse != "" {
			ref = test.reverse
		}

		revS, err := Strftime(tm, revfmt)
		if err != nil {
			t.Errorf("strftime(%q, %q): %v", tm, revfmt, err)
			continue
		}
		if ref != revS {
			t.Errorf("strftime(%q, %q): got %q, expected %q", tm, revfmt, revS, ref)
		}
	}
}
