package models

type Malware struct {
	ID     uint   `gorm:"primaryKey"`
	SHA256 string `gorm:"index;unique"`
	MD5    string
	SHA1   string
}
