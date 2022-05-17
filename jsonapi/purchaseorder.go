package jsonapi

import (
	"github.com/MaLowBar/moysklad-app-template/utils"
)

// PurchaseOrder содержит информацию о заказе поставщику. Подробнее смотри https://dev.moysklad.ru/doc/api/remap/1.2/documents/#dokumenty-zakaz-postawschiku
type PurchaseOrder struct {
	AccountID string `json:"accountId"`
	Agent     struct {
		Meta Meta `json:"meta"`
	} `json:"agent"`
	AgentAccount struct {
		Meta Meta `json:"meta"`
	} `json:"agentAccount,omitempty"`
	Applicable bool   `json:"applicable"`
	Code       string `json:"code,omitempty"`
	Contract   struct {
		Meta Meta `json:"meta"`
	} `json:"contract,omitempty"`
	Created               utils.MSJsonTime `json:"created"`
	Deleted               utils.MSJsonTime `json:"deleted,omitempty"`
	DeliveryPlannedMoment utils.MSJsonTime `json:"deliveryPlannedMoment,omitempty"`
	Description           string           `json:"description,omitempty"`
	ExternalCode          string           `json:"externalCode"`
	//Files               TODO: jsonapi.Files `json:"files"`
	Group struct {
		Meta Meta `json:"meta"`
	} `json:"group"`
	ID           string           `json:"id"`
	InvoicedSum  float64          `json:"invoicedSum,omitempty"`
	Meta         Meta             `json:"meta"`
	Moment       utils.MSJsonTime `json:"moment"`
	Name         string           `json:"name"`
	Organization struct {
		Meta Meta `json:"meta"`
	} `json:"organization"`
	OrganizationAccount struct {
		Meta Meta `json:"meta"`
	} `json:"organizationAccount,omitempty"`
	Owner struct {
		Meta Meta `json:"meta"`
	} `json:"owner"`
	PayedSum  float64 `json:"payedSum,omitempty"`
	Positions struct {
		Meta Meta `json:"meta"`
	} `json:"positions"`
	Printed bool `json:"printed"`
	Project struct {
		Meta Meta `json:"meta"`
	} `json:"project,omitempty"`
	Published bool `json:"published"`
	Rate      struct {
		Currency Meta    `json:"currency"`
		Value    float64 `json:"value,omitempty"`
	} `json:"rate,omitempty"`
	Shared     bool    `json:"shared"`
	ShippedSum float64 `json:"shippedSum,omitempty"`
	State      struct {
		Meta Meta `json:"meta"`
	} `json:"state,omitempty"`
	Store struct {
		Meta Meta `json:"meta"`
	} `json:"store,omitempty"`
	Sum         float64          `json:"sum,omitempty"`
	SyncID      string           `json:"syncId,omitempty"`
	Updated     utils.MSJsonTime `json:"updated"`
	VatEnabled  bool             `json:"vatEnabled"`
	VatIncluded bool             `json:"vatIncluded,omitempty"`
	VatSum      float64          `json:"vatSum,omitempty"`
	WaitSum     float64          `json:"waitSum,omitempty"`
	// TODO: Разобраться с полем Attributes

	CustomerOrders []struct {
		Meta Meta `json:"meta"`
	} `json:"customerOrders,omitempty"`
	InvoiceIns []struct {
		Meta Meta `json:"meta"`
	} `json:"invoiceIns,omitempty"`
	Payments []struct {
		Meta Meta `json:"meta"`
	} `json:"payments,omitempty"`
	Supplies []struct {
		Meta Meta `json:"meta"`
	} `json:"supplies,omitempty"`
	InternalOrder struct {
		Meta Meta `json:"meta"`
	} `json:"internalOrder,omitempty"`
}
