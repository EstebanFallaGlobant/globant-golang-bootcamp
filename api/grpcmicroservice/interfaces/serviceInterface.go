package interfaces

type ServiceInterface interface {
	IsPal(string) string
	Reverse(string) string
}
