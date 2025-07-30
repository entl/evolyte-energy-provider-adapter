package inverters

import "time"

// Response structure for solar inverter API calls
type SolarInverterResponse struct {
	Data       []SolarInverter `json:"data"`
	Pagination Pagination      `json:"pagination"`
}

type SolarInverter struct {
	ID              string          `json:"id"`
	UserID          string          `json:"userId"`
	Vendor          string          `json:"vendor"`
	LastSeen        time.Time       `json:"lastSeen"`
	IsReachable     bool            `json:"isReachable"`
	ProductionState ProductionState `json:"productionState"`
	Timezone        string          `json:"timezone"`
	Capabilities    Capabilities    `json:"capabilities"`
	Scopes          []string        `json:"scopes"`
	Information     Information     `json:"information"`
	Location        Location        `json:"location"`
}

// Current production state of the inverter
type ProductionState struct {
	ProductionRate          float64   `json:"productionRate"`
	IsProducing             bool      `json:"isProducing"`
	TotalLifetimeProduction float64   `json:"totalLifetimeProduction"`
	LastUpdated             time.Time `json:"lastUpdated"`
}

type Capabilities struct {
	ProductionState      Capability `json:"productionState"`
	ProductionStatistics Capability `json:"productionStatistics"`
}

type Capability struct {
	IsCapable       bool     `json:"isCapable"`
	InterventionIDs []string `json:"interventionIds"`
}

type Information struct {
	ID               string    `json:"id"`
	SerialNumber     string    `json:"sn"`
	Brand            string    `json:"brand"`
	Model            string    `json:"model"`
	SiteName         string    `json:"siteName"`
	InstallationDate time.Time `json:"installationDate"`
}

type Location struct {
	ID          string    `json:"id"`
	Longitude   float64   `json:"longitude"`
	Latitude    float64   `json:"latitude"`
	LastUpdated time.Time `json:"lastUpdated"`
}

type Pagination struct {
	After  string `json:"after"`
	Before string `json:"before"`
}

// InverterStatistic represents production statistics for an inverter
type InverterStatistic struct {
	Timezone    string                `json:"timezone"`
	Resolutions map[string]Resolution `json:"resolutions"`
	RetryAfter  time.Time             `json:"retryAfter,omitempty"`
}

// Resolution represents data for a specific time resolution (e.g., QUARTER_HOUR, DAY)
type Resolution struct {
	Unit string      `json:"unit"`
	Data []DataPoint `json:"data"`
}

type DataPoint struct {
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
}
