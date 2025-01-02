package parse

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"cloud.google.com/go/civil"
	"github.com/findyourpaths/phil/glr"
	"github.com/kr/pretty"
	// "github.com/ijt/go-anytime"
	// anytimev2 "github.com/ijt/go-anytime/v2"
	// "github.com/findyourpaths/paths/backend/internal/ident"
)

// A DateTimeTZs represents a sequence of date and time ranges.
//
// This type DOES include location information.
type DateTimeTZRanges struct {
	Items []*DateTimeTZRange
}

// A DateTimeTZRange represents a range of dates and times with time zones.
//
// This type DOES include location information.
type DateTimeTZRange struct {
	Start *DateTimeTZ
	End   *DateTimeTZ
}

// A DateTimeTZ represents a date and time with a time zone.
//
// This type DOES include location information.
type DateTimeTZ struct {
	civil.DateTime
	TimeZone string
}

// // Message for a sequence of datetime ranges.
// type DatetimeRanges struct {
// 	// state         protoimpl.MessageState
// 	// sizeCache     protoimpl.SizeCache
// 	// unknownFields protoimpl.UnknownFields

// 	Items []*DatetimeRange `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
// }

// // DatetimeRange is the model entity for the DatetimeRange schema.
// type DatetimeRange struct {
// 	// config `json:"-"`
// 	// ID of the ent.
// 	Id string `json:"id,omitempty"`
// 	// StartRfc3339 holds the value of the "start_rfc3339" field.
// 	StartRfc3339 string `json:"start_rfc3339,omitempty"`
// 	// Start holds the value of the "start" field.
// 	Start *timestamppb.Timestamp // time.Time `json:"start,omitempty"`
// 	// Duration holds the value of the "duration" field.
// 	Duration int64 `json:"duration,omitempty"`
// 	// EndRfc3339 holds the value of the "end_rfc3339" field.
// 	EndRfc3339 string `json:"end_rfc3339,omitempty"`
// 	// End holds the value of the "end" field.
// 	End *timestamppb.Timestamp //                   time.Time `json:"end,omitempty"`
// 	// event_datetime_ranges *string
// 	// selectValues          sql.SelectValues
// }

// func NewDatetimeRange(times ...time.Time) *DatetimeRange {
// 	start := times[0]
// 	r := &DatetimeRange{
// 		// Id:    ident.NewInt64("datetime_ranges"), // int64(rand.Int()), //id.String(),
// 		// Ksuid: ident.New().String(),
// 		Start:        timestamppb.New(start),
// 		StartRfc3339: start.Format(time.RFC3339),
// 	}
// 	r.Id = r.StartRfc3339

// 	if len(times) > 1 {
// 		end := times[1]
// 		r.End = timestamppb.New(end)
// 		r.EndRfc3339 = end.Format(time.RFC3339)
// 		r.Id += "-to-" + r.EndRfc3339
// 	}

// 	return r
// }

// daysre is a regexp to match day names, either long or short, regardless of case.
var daysRE = regexp.MustCompile(`(?i:\b` + strings.Join(
	[]string{
		"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday",
		"sun", "mon", "tue", "wed", "thu", "fri", "sat",
	}, `\b|\b`) + `\b,?)`)

// var daysRE = regexp.MustCompile(`(?i:\bthu|sun\b)`)

// Match a digit on one side and a letter on another. Used to separate `12pm`.
var boundaryRE1 = regexp.MustCompile(`([[:alpha:]])([[:^alpha:]])`)
var boundaryRE2 = regexp.MustCompile(`([[:^alpha:]])([[:alpha:]])`)

// func datetimeToTime(dt *civil.DateTime) time.Time {
// 	if dt.Date.Year == 0 {
// 		dt.Date.Year = 2025
// 	}
// 	return dt.In(time.UTC)
// }

func ExtractDateTimeTZRanges(mode, in string) (*DateTimeTZRanges, error) {
	defer func() {
		if err := recover(); err != nil {
			slog.Warn("in ExtractDatetimeRanges(), got a panic trying to extract", "in", in, "err", err)
		}
	}()

	ambiguousDateMode = mode
	// yyDebug = 3
	if yyDebug == 3 {
		fmt.Printf("in before: %q\n", in)
	}
	in = daysRE.ReplaceAllString(in, ``)
	in = boundaryRE1.ReplaceAllString(in, `$1 $2`)
	in = boundaryRE2.ReplaceAllString(in, `$1 $2`)
	in = strings.TrimPrefix(in, "When ")
	in = strings.TrimPrefix(in, "when ")
	in = CleanTextLine(in)
	if yyDebug == 3 {
		fmt.Printf("in after: %q\n", in)
	}

	forest, err := glr.Parse(parseRules, parseStates, NewDatetimeLexer(in))
	if yyDebug == 3 {
		fmt.Printf("tree:\n%# v\n", pretty.Formatter(forest))
		fmt.Printf("err: %v\n", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse datetime ranges from %q", in)
	}
	if len(forest) == 0 || len(forest[0].Children) == 0 {
		return nil, fmt.Errorf("no datetime ranges found in %q", in)
	}
	return NewDateTimeTZRangesFromParse(forest[0].Children[0])
}

func NewDateTimeTZRangesFromParse(root *glr.ParseNode) (*DateTimeTZRanges, error) {
	rs := &DateTimeTZRanges{}
	for _, dtr := range root.Children {
		for _, dt := range dtr.Children {
			for _, elt := range dt.Children {
				fmt.Printf("elt: %#v\n", elt)

				// package main

				// import (
				// 	"fmt"
				// 	"reflect"
				// )

				// type Time struct {
				// 	Hour        int `json:"hour"`        // The hour of the day in 24-hour format; range [0-23]
				// 	Minute      int `json:"minute"`      // The minute of the hour; range [0-59]
				// 	Second      int `json:"second"`      // The second of the minute; range [0-59]
				// 	Nanosecond int `json:"nanosecond"` // The nanosecond of the second; range [0-999999999]
				// }

				// func main() {
				// 	data := map[string]int{"Hour": 1, "Minute": 2}

				// 	timeStruct := Time{}
				// 	timeVal := reflect.ValueOf(&timeStruct).Elem()
				// 	timeType := timeVal.Type()

				// 	for i := 0; i < timeVal.NumField(); i++ {
				// 		field := timeType.Field(i)
				// 		fieldName := field.Name
				// 		value, ok := data[fieldName]
				// 		if ok {
				// 			timeVal.Field(i).SetInt(int64(value))
				// 		}
				// 	}

				// 	fmt.Println(timeStruct) // Output: {1 2 0 0}
				// }

			}
		}
	}
	// 	start := datetimeToTime(dtr.start)
	// 	var r *DatetimeRange
	// 	if dtr.end == nil {
	// 		r = NewDatetimeRange(start)
	// 	} else {
	// 		r = NewDatetimeRange(start, datetimeToTime(dtr.end))
	// 	}
	// 	rs.Items = append(rs.Items, r)
	// }
	// fmt.Printf("datetime: %q\n", t.Format(time.RFC3339))
	return rs, nil
	// return nil, nil
}

// // yyParse(

// rng2, _, err := anytimev2.ParseRange(in, time.Now(), anytimev2.Future)
// if err == nil {
// 	// log.Info().Msgf("in ExtractDateTimes(), anytimev2.ParseRange(%q)", in)
// 	// log.Info().Msgf("in ExtractDateTimes(), rng2.Start() %#v", rng2.Start())
// 	// log.Info().Msgf("in ExtractDateTimes(), rng2.End() %#v", rng2.End())

// 	return &DatetimeRanges{Items: []*DatetimeRange{
// 		NewDatetimeRange(rng2.Start(), rng2.End()),
// 	}}, nil
// }

// rng1, err := anytime.ParseRange(in, time.Now())
// if err == nil {
// 	// log.Info().Msgf("in ExtractDateTimes(), anytime.ParseRange(%q)", in)
// 	// log.Info().Msgf("in ExtractDateTimes(), rng.Start() %#v", rng.Start())
// 	// log.Info().Msgf("in ExtractDateTimes(), rng.End() %#v", rng.End())

// 	return &DatetimeRanges{Items: []*DatetimeRange{
// 		NewDatetimeRange(rng1.Start(), rng1.End()),
// 	}}, nil
// }

// t, err := anytime.Parse(in, time.Now())
// if err == nil {
// 	// log.Info().Msgf("in ExtractDateTimes(), anytime.Parse(%q)", in)
// 	// log.Info().Msgf("in ExtractDateTimes(), t %#v", t)

// 	return &DatetimeRanges{Items: []*DatetimeRange{
// 		NewDatetimeRange(t),
// 	}}, nil
// }
