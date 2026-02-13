package ical

import (
	// "bytes"
	"flag"
	"fmt"

	// "io/ioutil"
	"os"
	"path/filepath"
	"testing"

	// "github.com/bazelbuild/rules_go/go/tools/bazel"

	"github.com/segmentio/ksuid"
	// "google.golang.org/protobuf/encoding/prototext"
)

var writeTests = flag.Bool("write_tests", false, "write the tests outputs with the result of the current run")
var writeDir = flag.String("write_dir", "/tmp/internal/ical/testdata/", "write test outputs to this directory")

type rander struct {
	i int
}

func (r rander) Read(b []byte) (int, error) {
	b[0] = byte(r.i % 10)
	r.i += 1
	return 1, nil
}

func TestMain(m *testing.M) {
	// util.InTest = true
	ksuid.SetRand(rander{})
	os.Exit(m.Run())
}

func writeFile(path string, content string) error {
	fmt.Printf("writing to: %s\n", path)
	if err := os.MkdirAll(filepath.Dir(path), 0770); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}

// // See https://eli.thegreenplace.net/2022/file-driven-testing-in-go/
// func TestCalendarToVCalendar(t *testing.T) {
// 	type test struct {
// 		name      string
// 		useragent string
// 		want      string
// 	}

// 	tests := []test{
// 		{name: "calendar_google.ics", useragent: "Google-Calendar-Importer"},
// 		{name: "calendar_macos.ics", useragent: "macOS/13.2.1 (22D68) dataaccessd/1.0"},
// 		{name: "calendar_ios.ics", useragent: "iOS/16.2 (20C65) dataaccessd/1.0"},
// 	}

// 	return

// prs := &entpb.Results{}
// err := util.ReadProtoFile("internal/retrieve/testdata/iea_nine_points_com.textproto", prs)
// if err != nil {
// 	t.Fatalf("error: %v", err)
// }
// events := entity.SuccessfullyParsedEvents(prs)

// for _, testcase := range tests {
// 	cal := NewCalendar(events)
// 	got, err := CalendarToVCalendar(cal, testcase.useragent)
// 	if err != nil {
// 		t.Fatalf("error: %v", err)
// 	}

// 	gots := got.Serialize()
// 	if *writeTests {
// 		if err := writeFile(*writeDir+"/"+testcase.name, gots); err != nil {
// 			t.Fatalf("error: %v", err)
// 		}
// 	}

// 	wantr, err := os.Open("testdata/" + testcase.name)
// 	if err != nil {
// 		t.Fatalf("error: %v", err)
// 	}
// 	want, err := ics.ParseCalendar(wantr)
// 	if err != nil {
// 		t.Fatalf("error: %v", err)
// 	}
// 	if err = CleanParsedVCalendar(want); err != nil {
// 		t.Fatalf("error: %v", err)
// 	}

// 	if diff := cmp.Diff(want, got); diff != "" {
// 		t.Fatalf("unexpected diff (-want +got):\n%s", diff)
// 	}
// }
// }
