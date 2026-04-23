package services

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Filters struct {
	Gender    string
	AgeGroup  string
	CountryID string
	MinAge    *int
	MaxAge    *int
}

// countryMap maps country names to ISO codes
var countryMap = map[string]string{
	"nigeria":                  "NG",
	"kenya":                    "KE",
	"angola":                   "AO",
	"benin":                    "BJ",
	"ghana":                    "GH",
	"cameroon":                 "CM",
	"south africa":             "ZA",
	"ethiopia":                 "ET",
	"tanzania":                 "TZ",
	"uganda":                   "UG",
	"mozambique":               "MZ",
	"madagascar":               "MG",
	"malawi":                   "MW",
	"zambia":                   "ZM",
	"zimbabwe":                 "ZW",
	"botswana":                 "BW",
	"namibia":                  "NA",
	"senegal":                  "SN",
	"mali":                     "ML",
	"burkina faso":             "BF",
	"niger":                    "NE",
	"chad":                     "TD",
	"somalia":                  "SO",
	"rwanda":                   "RW",
	"burundi":                  "BI",
	"togo":                     "TG",
	"sierra leone":             "SL",
	"liberia":                  "LR",
	"mauritania":               "MR",
	"eritrea":                  "ER",
	"gambia":                   "GM",
	"guinea":                   "GN",
	"guinea-bissau":            "GW",
	"equatorial guinea":        "GQ",
	"gabon":                    "GA",
	"congo":                    "CG",
	"drc":                      "CD",
	"central african republic": "CF",
	"lesotho":                  "LS",
	"swaziland":                "SZ",
	"eswatini":                 "SZ",
	"djibouti":                 "DJ",
	"comoros":                  "KM",
	"cape verde":               "CV",
	"seychelles":               "SC",
	"mauritius":                "MU",
	"sao tome":                 "ST",
}

func ParseQuery(q string) (*Filters, error) {
	q = strings.ToLower(strings.TrimSpace(q))
	if q == "" {
		return nil, fmt.Errorf("empty query")
	}

	filters := &Filters{}

	// Gender - if both are mentioned, don't filter by gender
	hasMale := strings.Contains(q, "male")
	hasFemale := strings.Contains(q, "female")

	if hasFemale && hasMale {
		// both genders requested — no gender filter
		filters.Gender = ""
	} else if hasFemale {
		filters.Gender = "female"
	} else if hasMale {
		filters.Gender = "male"
	}

	// Age groups
	if strings.Contains(q, "child") || strings.Contains(q, "children") {
		filters.AgeGroup = "child"
	} else if strings.Contains(q, "teenager") || strings.Contains(q, "teen") {
		filters.AgeGroup = "teenager"
	} else if strings.Contains(q, "adult") {
		filters.AgeGroup = "adult"
	} else if strings.Contains(q, "senior") || strings.Contains(q, "elderly") {
		filters.AgeGroup = "senior"
	}

	// "young" maps to ages 16-24
	if strings.Contains(q, "young") {
		min := 16
		max := 24
		filters.MinAge = &min
		filters.MaxAge = &max
	}

	// "above X" or "over X"
	aboveRegex := regexp.MustCompile(`(?:above|over)\s+(\d+)`)
	if match := aboveRegex.FindStringSubmatch(q); len(match) > 1 {
		age, err := strconv.Atoi(match[1])
		if err == nil {
			filters.MinAge = &age
		}
	}

	// "below X" or "under X"
	belowRegex := regexp.MustCompile(`(?:below|under)\s+(\d+)`)
	if match := belowRegex.FindStringSubmatch(q); len(match) > 1 {
		age, err := strconv.Atoi(match[1])
		if err == nil {
			filters.MaxAge = &age
		}
	}

	// "between X and Y"
	betweenRegex := regexp.MustCompile(`between\s+(\d+)\s+and\s+(\d+)`)
	if match := betweenRegex.FindStringSubmatch(q); len(match) > 2 {
		min, err1 := strconv.Atoi(match[1])
		max, err2 := strconv.Atoi(match[2])
		if err1 == nil && err2 == nil {
			filters.MinAge = &min
			filters.MaxAge = &max
		}
	}

	// Country - check longest names first to avoid partial matches
	for name, code := range countryMap {
		if strings.Contains(q, name) {
			filters.CountryID = code
			break
		}
	}

	// Validate that at least one filter was set
	if filters.Gender == "" && filters.AgeGroup == "" &&
		filters.CountryID == "" && filters.MinAge == nil && filters.MaxAge == nil {
		return nil, fmt.Errorf("unable to interpret query")
	}

	return filters, nil
}
