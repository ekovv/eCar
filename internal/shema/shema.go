package shema

type RegNum struct {
	RegNums []string `json:"regNums"`
}

type Car struct {
	RegNum string `json:"regNum"`
	Mark   string `json:"mark"`
	Model  string `json:"model"`
	Year   int    `json:"year"`
	Owner  People `json:"owner"`
}

type People struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

type Filter struct {
	RegNum          string `json:"regNum"`
	Mark            string `json:"mark"`
	Model           string `json:"model"`
	Year            int    `json:"year"`
	OwnerName       string `json:"owner_name"`
	OwnerSurname    string `json:"owner_surname"`
	OwnerPatronymic string `json:"owner_patronymic"`
	Page            int    `json:"page"`
	Limit           int    `json:"limit"`
}
