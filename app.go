package moyskladapptemplate

type App interface {
	GetBase() *BaseApp
}

type BaseApp struct {
	ID          string
	UID         string
	AccountID   string
	Status      AppStatus
	AccessToken string
	SecretKey   string
}

type AppStatus string

const (
	APP_STATUS_ACTIVATED         = "Activated"
	APP_STATUS_SETTINGS_REQUIRED = "Settings required"
	APP_STATUS_ACTIVATING        = "Activating"
	APP_STATUS_INACTIVE          = "Inactive"
)

func NewBaseApp(id, uid, accountId, accessToken, secretKey string) *BaseApp {
	return &BaseApp{
		ID:          id,
		UID:         uid,
		AccountID:   accountId,
		Status:      APP_STATUS_INACTIVE,
		AccessToken: accessToken,
		SecretKey:   secretKey,
	}
}
