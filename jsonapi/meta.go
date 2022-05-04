package jsonapi

type Meta struct {
	Href         string `json:"href"`
	MetadataHref string `json:"metadataHref"`
	Type         string `json:"type"`
	UUIDHref     string `json:"uuidHref,omitempty"`
	DownloadHref string `json:"downloadHref,omitempty"`
	Size         string `json:"size,omitempty"`
	Limit        string `json:"limit,omitempty"`
	Offset       string `json:"offset,omitempty"`
}
