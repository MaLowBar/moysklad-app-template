package jsonapi

import (
	"github.com/MaLowBar/moysklad-app-template/utils"
)

// PurchaseOrder содержит информацию о заказе поставщику. Подробнее смотри https://dev.moysklad.ru/doc/api/remap/1.2/documents/#dokumenty-zakaz-postawschiku
type PurchaseOrder struct {
	AccountID             string           `json:"accountId"`
	Agent                 Meta             `json:"agent"`
	AgentAccount          Meta             `json:"agentAccount,omitempty"`
	Applicable            bool             `json:"applicable"`
	Code                  string           `json:"code,omitempty"`
	Contract              Meta             `json:"contract,omitempty"`
	Created               utils.MSJsonTime `json:"created"`
	Deleted               utils.MSJsonTime `json:"deleted,omitempty"`
	DeliveryPlannedMoment utils.MSJsonTime `json:"deliveryPlannedMoment,omitempty"`
	Description           string           `json:"description,omitempty"`
	ExternalCode          string           `json:"externalCode"`
	//Files                 TODO: jsonapi.Files           `json:"files"`
	Group               Meta             `json:"group"`
	ID                  string           `json:"id"`
	InvoicedSum         float64          `json:"invoicedSum,omitempty"`
	Meta                Meta             `json:"meta"`
	Moment              utils.MSJsonTime `json:"moment"`
	Name                string           `json:"name"`
	Organization        Meta             `json:"organization"`
	OrganizationAccount Meta             `json:"organizationAccount,omitempty"`
	Owner               Meta             `json:"owner"`
	PayedSum            float64          `json:"payedSum,omitempty"`
	Positions           []Meta           `json:"positions"`
	Printed             bool             `json:"printed"`
	Project             Meta             `json:"project,omitempty"`
	Published           bool             `json:"published"`
	Rate                struct {
		Currency Meta    `json:"currency"`
		Value    float64 `json:"value,omitempty"`
	} `json:"rate,omitempty"`
	Shared      bool             `json:"shared"`
	ShippedSum  float64          `json:"shippedSum,omitempty"`
	State       Meta             `json:"state,omitempty"`
	Store       Meta             `json:"store,omitempty"`
	Sum         float64          `json:"sum,omitempty"`
	SyncID      string           `json:"syncId,omitempty"`
	Updated     utils.MSJsonTime `json:"updated"`
	VatEnabled  bool             `json:"vatEnabled"`
	VatIncluded bool             `json:"vatIncluded,omitempty"`
	VatSum      float64          `json:"vatSum,omitempty"`
	WaitSum     float64          `json:"waitSum,omitempty"`
	// TODO: Разобраться с полем Attributes

	CustomerOrders []Meta `json:"customerOrders,omitempty"`
	InvoiceIns     []Meta `json:"invoiceIns,omitempty"`
	Payments       []Meta `json:"payments,omitempty"`
	Supplies       []Meta `json:"supplies,omitempty"`
	InternalOrder  Meta   `json:"internalOrder,omitempty"`
}

type PurchaseOrders struct {
	Meta    Meta `json:"meta"`
	Context struct {
		Employee struct {
			Meta Meta `json:"meta"`
		} `json:"employee"`
	} `json:"context"`
	Rows []PurchaseOrder `json:"rows"`
}
