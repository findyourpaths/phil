package ical

import (
	"fmt"

	ics "github.com/arran4/golang-ical"
	// "github.com/segmentio/ksuid"
)

// func NewCalendar(events *entpb.Events) *entpb.Calendar {
// 	log.Info().Msgf("Creating calendar with %d events.", len(events.Items))
// 	return &entpb.Calendar{
// 		Id:     ksuid.New().String(),
// 		Events: events,
// 	}
// }

// func CalendarToVCalendar(c *entpb.Calendar, useragent string) (*ics.Calendar, error) {
// 	r := ics.NewCalendar()

// 	r.SetXWRCalName("Paths")
// 	r.SetRefreshInterval("P1D")
// 	setCalendarProperty(r, ics.PropertyComment, "USER-AGENT: "+useragent)
// 	// UserAgent(c.UserAgent)

// 	htmlmode := HTMLInDescription
// 	if useragent == "Google-Calendar-Importer" {
// 		htmlmode = HTMLInDescription
// 	} else if strings.HasPrefix(useragent, "iOS") || strings.HasPrefix(useragent, "macOS") {
// 		htmlmode = OnlyURLsInDescription
// 	}

// 	// r.SetDescription(e.Description)
// 	if c.Events == nil {
// 		return r, nil
// 	}

// 	setCalendarProperty(r, ics.PropertyComment, fmt.Sprintf("EVENT-DATA: %d events", len(c.Events.Items)))
// 	for _, event := range c.Events.Items {
// 		vevent, err := EventToVEvent(event, htmlmode)
// 		if err != nil {
// 			return nil, err
// 		}
// 		r.AddVEvent(vevent)
// 	}

// 	return r, nil
// }

// CleanParsedVCalendar fixes fields to ensure than vcal == ics.ParseCalendar(vcal.Serialize()).
func CleanParsedVCalendar(vcal *ics.Calendar) error {
	if vcal.Components == nil {
		vcal.Components = []ics.Component{}
	}
	calprop := getCalendarProperty(vcal, ics.Property("REFRESH-INTERVAL"))
	if calprop == nil {
		return fmt.Errorf("expected to find calendar property")
	}
	calprop.IANAToken = "REFRESH-INTERVAL;VALUE=DURATION"
	calprop.ICalParameters = map[string][]string{} //"VALUE": {"DURATION"}}
	return nil
}

// copied from github.com/arran4/golang-ical/calendar.go
func setCalendarProperty(calendar *ics.Calendar, property ics.Property, value string, props ...ics.PropertyParameter) {
	for i := range calendar.CalendarProperties {
		if calendar.CalendarProperties[i].IANAToken == string(property) {
			calendar.CalendarProperties[i].Value = value
			calendar.CalendarProperties[i].ICalParameters = map[string][]string{}
			for _, p := range props {
				k, v := p.KeyValue()
				calendar.CalendarProperties[i].ICalParameters[k] = v
			}
			return
		}
	}
	r := ics.CalendarProperty{
		ics.BaseProperty{
			IANAToken:      string(property),
			Value:          value,
			ICalParameters: map[string][]string{},
		},
	}
	for _, p := range props {
		k, v := p.KeyValue()
		r.ICalParameters[k] = v
	}
	calendar.CalendarProperties = append(calendar.CalendarProperties, r)
}

func getCalendarProperty(cal *ics.Calendar, prop ics.Property) *ics.CalendarProperty {
	for i := range cal.CalendarProperties {
		calprop := &(cal.CalendarProperties[i])
		if calprop.IANAToken == string(prop) {
			return calprop
		}
	}
	return nil
}

/*
      public static Calendar getCalendarWithEvents(String userAgent, List<Event> events) {
       // Create calendar
       Calendar calendar = new Calendar();
       calendar.replace(new DtStamp(Instant.EPOCH));
       calendar.add(new ProdId("-//Ben Fortuna//iCal4j 1.0//EN"));
       calendar.add(Version.VERSION_2_0);
       String name = "Connection Central";
       calendar.add(new XProperty("X-WR-CALNAME", name));
       calendar.add(new Name(name));
       calendar.add(CalScale.GREGORIAN);
       // calendar.add(new XProperty("X-WR-CALNAME", "Connection Central"));
       ParameterList refresh = new ParameterList();
       refresh.add(new Value("DURATION"));
       calendar.add(new RefreshInterval(refresh, java.time.Duration.ofMinutes(7 * 24 * 60)));

       //calendar.add(new RefreshInterval(new ParameterList(List.of(new Value("DURATION"))), java.time.Duration.ofMinutes(7 * 24 * 60)));

       // Add metadata in comments
       //
       calendar.add(new Comment("USER-AGENT: " + userAgent));

       // Do a little browser sniffing to return calendar descriptions that can be rendered as best as possible.
       Event.HTMLMode htmlMode = Event.HTMLMode.HTML_IN_DESCRIPTION;
       if (userAgent.equals("Google-Calendar-Importer")) {
           htmlMode = Event.HTMLMode.HTML_IN_DESCRIPTION;
       } else if (userAgent.startsWith("iOS") ||
           userAgent.startsWith("macOS")) {
           // "macOS/13.2.1 (22D68) dataaccessd/1.0"
           // iOS/16.2 (20C65) dataaccessd/1.0"
           htmlMode = Event.HTMLMode.ONLY_URLS_IN_DESCRIPTION;
       }

       calendar.add(new Comment("EVENT-DATA: " + events.size() + " events and " + Person.personsBySlug.size() + " persons"));
       calendar.add(new Comment("HTML-MODE: " + htmlMode));
       for (Event event : events) {
           if (event != null) {
               calendar.add(event.toVEvent(htmlMode));
           }
       }
       return calendar;
   }
*/
