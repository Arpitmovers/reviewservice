package models

type Provider struct {
	ProviderID   int    `gorm:"primaryKey;column:provider_id"`
	ProviderName string `gorm:"size:100;not null"`
}

func (Provider) TableName() string {
	return "Providers"
}
