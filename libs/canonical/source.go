package canonical

type Source struct {
	Integration string `json:"integration"`
	Protocol    string `json:"protocol"`
	Network     string `json:"network"`
	DeviceID    string `json:"device_id"` // The device ID from the source system, if applicable (The source named it)
}
