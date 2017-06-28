package error

type ServiceError interface {
	Error() string
}
