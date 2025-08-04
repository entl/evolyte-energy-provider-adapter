package inverters

import "time"

// Response structure for solar inverter API calls
type SolarInverterResponse struct {
	Data       []SolarInverter `json:"data"`
	Pagination Pagination      `json:"pagination"`
}

type SolarInverter struct {
	ID              string          `json:"id" validate:"required"`
	UserID          string          `json:"userId" validate:"required"`
	Vendor          string          `json:"vendor" validate:"required"`
	LastSeen        time.Time       `json:"lastSeen" validate:"required"`
	IsReachable     bool            `json:"isReachable" validate:"required"`
	ProductionState ProductionState `json:"productionState" validate:"required"`
	Timezone        string          `json:"timezone"`
	Capabilities    Capabilities    `json:"capabilities" validate:"required"`
	Scopes          []string        `json:"scopes" validate:"required"`
	Information     Information     `json:"information" validate:"required"`
	Location        Location        `json:"location" validate:"required"`
}

// Current production state of the inverter
type ProductionState struct {
	ProductionRate          float64   `json:"productionRate"`
	IsProducing             bool      `json:"isProducing"`
	TotalLifetimeProduction float64   `json:"totalLifetimeProduction"`
	LastUpdated             time.Time `json:"lastUpdated"`
}

type Capabilities struct {
	ProductionState      Capability `json:"productionState" validate:"required"`
	ProductionStatistics Capability `json:"productionStatistics" validate:"required"`
}

type Capability struct {
	IsCapable       bool     `json:"isCapable" validate:"required"`
	InterventionIDs []string `json:"interventionIds" validate:"required"`
}

type Information struct {
	ID               string    `json:"id" validate:"required"`
	SerialNumber     *string   `json:"sn"`
	Brand            string    `json:"brand" validate:"required"`
	Model            string    `json:"model" validate:"required"`
	SiteName         string    `json:"siteName" validate:"required"`
	InstallationDate time.Time `json:"installationDate" validate:"required"`
}

type Location struct {
	ID          string    `json:"id" validate:"omitempty, uuid"`
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
	Timezone    string                `json:"timezone" validate:"required"`
	Resolutions map[string]Resolution `json:"resolutions" validate:"required"`
	RetryAfter  time.Time             `json:"retryAfter"`
}

// Resolution represents data for a specific time resolution (e.g., QUARTER_HOUR, DAY)
type Resolution struct {
	Unit string      `json:"unit" validate:"required"`
	Data []DataPoint `json:"data"`
}

type DataPoint struct {
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
}

// LinkInverterRequest represents the request body for linking an inverter
type LinkInverterRequest struct {
	Scopes      []string `json:"scopes" validate:"required"`
	Language    string   `json:"language" validate:"required"`
	RedirectUri string   `json:"redirectUri" validate:"required"`
}

type LinkInverterResponse struct {
	LinkURL   string `json:"linkUrl" validate:"required"`
	LinkToken string `json:"linkToken" validate:"required"`
}

type AddInverterRequest struct {
	UserID                  string    `json:"userId" validate:"required"`
	Vendor                  string    `json:"vendor" validate:"required"`
	Model                   string    `json:"model" validate:"required"`
	SerialNumber            string    `json:"serialNumber" validate:"required"`
	TotalLifetimeProduction float64   `json:"totalLifetimeProduction" validate:"required"`
	InstallationDate        time.Time `json:"installationDate" validate:"required"`
}

type AddInverterResponse struct {
	ID                      string    `json:"id" validate:"required"`
	UserID                  string    `json:"userId" validate:"required"`
	Vendor                  string    `json:"vendor" validate:"required"`
	Model                   string    `json:"model" validate:"required"`
	SerialNumber            string    `json:"serialNumber" validate:"required"`
	TotalLifetimeProduction float64   `json:"totalLifetimeProduction" validate:"required"`
	InstallationDate        time.Time `json:"installationDate" validate:"required"`
}

type EnodeErrorResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Details string `json:"details"`
}
