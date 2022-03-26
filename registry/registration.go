package registry

type Registration struct {
	ServiceName      ServiceName   `json:"service_name"`
	ServiceURL       string        `json:"service_url"`
	RequireServices  []ServiceName `json:"require_services"`
	ServiceUpdateURL string        `json:"service_update_url"`
}

type ServiceName string

const (
	LogService     = ServiceName("LogService")
	GradingService = ServiceName("GradingService")
	Portal         = ServiceName("Portal")
)

type patchEntry struct {
	Name ServiceName `json:"name"`
	URL  string      `json:"url"`
}

type patch struct {
	Added   []patchEntry `json:"added"`
	Removed []patchEntry `json:"removed"`
}
