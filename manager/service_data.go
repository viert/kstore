package manager

import "encoding/json"

// ServiceInfo represents all the data concerned with one service
type ServiceInfo struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Comment   string `json:"comment"`
	UpdatedAt string `json:"updated_at"`
	URL       string `json:"url"`
}

type serviceData map[string]*ServiceInfo

func (sd *serviceData) dump() ([]byte, error) {
	return json.Marshal(sd)
}
