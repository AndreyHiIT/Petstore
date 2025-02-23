package errors

const (
	NoError = iota
	InternalError
	GeneralError
)

const (
	HashPasswordError = iota + 1000
)

const (
	AuthServiceGeneralErr = iota + 2000
	AuthServiceWrongPasswordErr
	AuthServiceAccessTokenGenerationErr
	AuthServiceRefreshTokenGenerationErr
	AuthServiceUserNotVerified
	AuthServiceVerifyErr
	AuthGenerateHashErr
	AuthUrlParseErr
	NotifyEmailSendErr
	UserServiceWrongPhoneCodeErr
	UserServiceCreateUserErr
	UserServiceUserAlreadyExists
	UserServiceRetrieveUserErr
	UserServiceUpdateErr
	AddPetErr
	PetServiceUpdateErr
	PetServiceFindPetbyStatus
	FindPetbyIDErrorDuringConversion
	PetServiceFindPetbyID
	PetServiceErrPetNotFound
	UpdatePetFormErrorDuringConversion
	PetServiceUpdatePetFormBadReuest
	UpdatePetFormError
	OrderServiceCreateOrderErr
	FindOrderByIDErrorDuringConversion
	FindOrderByIDErrorIDLessZero
	OrderServiceFindByIDNotFoundID
	OrderServiceFindByIDIternalErr
	DeleteOrderByIDErrorDuringConversion
	DeleteOrderByIDErrorIDLessZero
	OrderServiceDeleteByIDNotFoundID
	OrderServiceDeleteByIDIternalErr
)
