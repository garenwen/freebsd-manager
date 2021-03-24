package v1

type status string

const (
	StatusSuccess status = "success"
	StatusError   status = "error"
)

type errorType string

const (
	ErrorNone        errorType = ""
	ErrorTimeout     errorType = "timeout"
	ErrorCanceled    errorType = "canceled"
	ErrorBadData     errorType = "bad_data"
	ErrorExec        errorType = "execution"
	ErrorInternal    errorType = "internal"
	ErrorUnavailable errorType = "unavailable"
	ErrorNotFound    errorType = "not_found"
)

type BaseResult struct {
	Status   status      `json:"status"`
	Data     interface{} `json:"data,omitempty"`
	ApiError *ApiError   `json:"api_error"`
}

type ApiError struct {
	Typ errorType `json:"typ"`
	Msg string    `json:"msg"`
}

type CreateVolumeRequest struct {
	Name          string            `json:"name,omitempty"`
	CapacityRange *CapacityRange    `json:"capacity_range,omitempty"`
	Parameters    map[string]string `json:"parameters,omitempty"`
}

type CapacityRange struct {
	RequiredBytes int64 `json:"required_bytes,omitempty"`
	LimitBytes    int64 `json:"limit_bytes,omitempty"`
}

type CreateVolumeResponse struct {
	Name          string `json:"name"`
	Origin        string `json:"origin"`
	Used          uint64 `json:"used"`
	Avail         uint64 `json:"avail"`
	Mountpoint    string `json:"mountpoint"`
	Compression   string `json:"compression"`
	Type          string `json:"type"`
	Written       uint64 `json:"written"`
	Volsize       uint64 `json:"volsize"`
	Logicalused   uint64 `json:"logicalused"`
	Usedbydataset uint64 `json:"usedbydataset"`
	Quota         uint64 `json:"quota"`
	Referenced    uint64 `json:"referenced"`
}
