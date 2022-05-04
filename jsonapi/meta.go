package jsonapi

type Meta struct {
	Href         string `json:"href"`
	MetadataHref string `json:"metadataHref"`
	Type         string `json:"type"`
	UUIDHref     string `json:"uuidHref,omitempty"`
	DownloadHref string `json:"downloadHref,omitempty"`
	Size         int    `json:"size,omitempty"`
	Limit        int    `json:"limit,omitempty"`
	Offset       int    `json:"offset,omitempty"`
}
