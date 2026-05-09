package models

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin       Role = "admin"
	RoleKoordinator Role = "koordinator_pengawas"
	RolePengawas    Role = "pengawas"
	RoleJagal       Role = "jagal"
	RoleKulit       Role = "kulit"
	RoleCacahDaging Role = "cacah_daging"
	RoleCacahTulang Role = "cacah_tulang"
	RolePacking     Role = "packing"
	RoleDistribusi  Role = "distribusi"
)

// ValidRoles is the authoritative set of all accepted roles.
var ValidRoles = map[Role]bool{
	RoleAdmin: true, RoleKoordinator: true, RolePengawas: true,
	RoleJagal: true, RoleKulit: true, RoleCacahDaging: true,
	RoleCacahTulang: true, RolePacking: true, RoleDistribusi: true,
}

type User struct {
	gorm.Model
	NamaLengkap string `gorm:"not null" json:"nama_lengkap"`
	Username    string `gorm:"size:100;uniqueIndex;not null" json:"username"`
	Password    string `gorm:"not null" json:"-"`
	Role        Role   `gorm:"type:enum('admin','koordinator_pengawas','pengawas','jagal','kulit','cacah_daging','cacah_tulang','packing','distribusi');default:'pengawas'" json:"role"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

// BeforeUpdate hashes the password only if it's a new plaintext value.
// Bcrypt hashes always start with "$2a$" or "$2b$", so we use that as the check.
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if u.Password != "" && !strings.HasPrefix(u.Password, "$2") {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashed)
	}
	return nil
}

func (u *User) CheckPassword(plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain)) == nil
}
