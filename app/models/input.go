package models

type Input struct {
	Name          string `json:"name"`
	MelliNumber   string `json:"mellinumber" gorm:"unique"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	KartInfo      string `json:"kartinfo"`
	TarikhTavalod string `json:"tarikh_tavalod"`
	Code          string `json:"code"`
	Amount        string `json:"amount"`
	Authority     string `json:"authority"`
	AdminToken    string `json:"admin_token"`
}
