package timefmt

const (
	Date             = "2006-01-02" // 2026-10-20
	DateCompact      = "20060102"   // 20261020
	DateDot          = "02.01.2006" // 20.10.2026
	DateSlash        = "02/01/2006" // 20/10/2026
	DateUS           = "01/02/2006" // 10/20/2026
	DateShort        = "02.01"      // 20.10
	DateShortSlash   = "02/01"      // 20/10
	DateShortUS      = "01/02"      // 10/20
	YearMonth        = "2006-01"    // 2026-10
	YearMonthCompact = "200601"     // 202610

	Time          = "15:04"              // 18:30
	TimeSeconds   = "15:04:05"           // 18:30:45
	TimeMillis    = "15:04:05.000"       // 18:30:45.123
	TimeMicros    = "15:04:05.000000"    // 18:30:45.123456
	TimeNanos     = "15:04:05.000000000" // 18:30:45.123456789
	Time12        = "3:04 PM"            // 6:30 PM
	Time12Seconds = "3:04:05 PM"         // 6:30:45 PM

	DateTimeShort        = "02.01 15:04"              // 20.10 18:30
	DateTimeShortSeconds = "02.01 15:04:05"           // 20.10 18:30:45
	DateTimeShortMillis  = "02.01 15:04:05.000"       // 20.10 18:30:45.123
	DateTimeShortMicros  = "02.01 15:04:05.000000"    // 20.10 18:30:45.123456
	DateTimeShortNanos   = "02.01 15:04:05.000000000" // 20.10 18:30:45.123456789
	DateTimeShortSlash   = "02/01 15:04"              // 20/10 18:30
	DateTimeShortUS      = "01/02 15:04"              // 10/20 18:30

	DateTime        = "2006-01-02 15:04"           // 2026-10-20 18:30
	DateTimeSeconds = "2006-01-02 15:04:05"        // 2026-10-20 18:30:45
	DateTimeMillis  = "2006-01-02 15:04:05.000"    // 2026-10-20 18:30:45.123
	DateTime12      = "2006-01-02 3:04 PM"         // 2026-10-20 6:30 PM
	DateTimeOffset  = "2006-01-02 15:04:05 -07:00" // 2026-10-20 18:30:45 +02:00
	DateTimeDot     = "02.01.2006 15:04"           // 20.10.2026 18:30
	DateTimeSlash   = "02/01/2006 15:04"           // 20/10/2026 18:30
	DateTimeUS      = "01/02/2006 15:04"           // 10/20/2026 18:30

	HumanDate          = "2 January 2006"         // 20 October 2026
	HumanDateShort     = "2 Jan 2006"             // 20 Oct 2026
	HumanDateWeekday   = "Monday, 2 January 2006" // Tuesday, 20 October 2026
	HumanDateCompact   = "Mon, 2 Jan"             // Tue, 20 Oct
	HumanDay           = "2 January"              // 20 October
	HumanDayShort      = "2 Jan"                  // 20 Oct
	HumanMonth         = "January 2006"           // October 2026
	HumanMonthShort    = "Jan 2006"               // Oct 2026
	HumanDateTime      = "2 January 2006, 15:04"  // 20 October 2026, 18:30
	HumanDateTimeShort = "2 Jan 2006, 15:04"      // 20 Oct 2026, 18:30

	DateTimeURL     = "2006-01-02-15-04-05" // 2026-10-20-18-30-45
	DateTimeFile    = "2006-01-02_15-04-05" // 2026-10-20_18-30-45
	DateTimeCompact = "20060102T150405"     // 20261020T183045
)
