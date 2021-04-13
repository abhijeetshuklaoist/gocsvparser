package parser

type Person struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Wage           string `json:"wage"`
	EmployeeNumber string `json:"number"`
}

func (p Person) value() []string {
	return [] string {p.FirstName, p.LastName, p.Email, p.Wage, p.EmployeeNumber}
}

func (p Person) fields() []string {
	return [] string {"FirstName", "LastName", "Email", "Wage", "EmployeeNumber"}
}
