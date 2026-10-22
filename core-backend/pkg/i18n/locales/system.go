package locales

import (
	"github.com/example/core-backend/pkg/constants/response_code"
	"github.com/example/core-backend/pkg/i18n"
)

var SystemMessagesEN = map[string]string{
	response_code.InternalServerGeneral.Code(): "An internal server error occurred. Please try again later.",
	response_code.DatabaseError.Code():         "A database error occurred.",
	response_code.ThirdPartyServiceError.Code():"External service integration error.",
	response_code.BadRequestGeneral.Code():     "Invalid request parameters or payload.",
	response_code.BadRequestValidation.Code():  "Invalid request payload structure or syntax.",
	response_code.NotFoundGeneral.Code():       "The requested resource was not found.",
	response_code.NotFoundResource.Code():      "Resource not found.",
	response_code.CommonError.Code():           "Validation failed on one or more fields.",
	response_code.RequiredField.Code():         "This field is required.",
	response_code.NotEmail.Code():              "Invalid email address format.",
	response_code.NotUrl.Code():                "Invalid URL format.",
	response_code.NotAlphanum.Code():           "Field must contain only letters and numbers.",
	response_code.NotAlpha.Code():              "Field must contain only letters.",
	response_code.NotNumeric.Code():            "Field must contain only numeric characters.",
	response_code.NotNumber.Code():             "Field must be a valid number.",
	response_code.NotUuid.Code():               "Field must be a valid UUID.",
	response_code.NotIP.Code():                 "Field must be a valid IP address.",
	response_code.NotHexColor.Code():           "Field must be a valid HEX color code.",
	response_code.NotOneOf.Code():              "Value is not within the permitted options.",
	response_code.MinLength.Code():             "Input does not meet the minimum length constraint.",
	response_code.MaxLength.Code():             "Input exceeds the maximum length constraint.",
	response_code.LengthExact.Code():           "Input length does not match required length.",
	response_code.MustGreatThan.Code():         "Value must be greater than required minimum.",
	response_code.MustGreatThanEqual.Code():    "Value must be greater than or equal to required minimum.",
	response_code.MustLessThan.Code():          "Value must be less than required maximum.",
	response_code.MustLessThanEqual.Code():     "Value must be less than or equal to required maximum.",
}

var SystemMessagesVI = map[string]string{
	response_code.InternalServerGeneral.Code(): "Đã xảy ra lỗi máy chủ nội bộ. Vui lòng thử lại sau.",
	response_code.DatabaseError.Code():         "Đã xảy ra lỗi cơ sở dữ liệu.",
	response_code.ThirdPartyServiceError.Code():"Lỗi kết nối dịch vụ bên thứ ba.",
	response_code.BadRequestGeneral.Code():     "Tham số hoặc dữ liệu yêu cầu không hợp lệ.",
	response_code.BadRequestValidation.Code():  "Định dạng hoặc cấu trúc dữ liệu không hợp lệ.",
	response_code.NotFoundGeneral.Code():       "Không tìm thấy tài nguyên được yêu cầu.",
	response_code.NotFoundResource.Code():      "Tài nguyên không tồn tại.",
	response_code.CommonError.Code():           "Dữ liệu nhập vào chưa hợp lệ.",
	response_code.RequiredField.Code():         "Trường này là bắt buộc.",
	response_code.NotEmail.Code():              "Định dạng email không hợp lệ.",
	response_code.NotUrl.Code():                "Định dạng đường dẫn URL không hợp lệ.",
	response_code.NotAlphanum.Code():           "Trường này chỉ được chứa chữ cái và chữ số.",
	response_code.NotAlpha.Code():              "Trường này chỉ được chứa chữ cái.",
	response_code.NotNumeric.Code():            "Trường này chỉ được chứa số.",
	response_code.NotNumber.Code():             "Giá trị phải là một số hợp lệ.",
	response_code.NotUuid.Code():               "Giá trị phải là định dạng UUID hợp lệ.",
	response_code.NotIP.Code():                 "Giá trị phải là địa chỉ IP hợp lệ.",
	response_code.NotHexColor.Code():           "Mã màu HEX không hợp lệ.",
	response_code.NotOneOf.Code():              "Giá trị không nằm trong danh sách cho phép.",
	response_code.MinLength.Code():             "Độ dài ký tự chưa đạt mức tối thiểu.",
	response_code.MaxLength.Code():             "Độ dài ký tự vượt quá mức cho phép.",
	response_code.LengthExact.Code():           "Độ dài ký tự không đúng với yêu cầu.",
	response_code.MustGreatThan.Code():         "Giá trị phải lớn hơn mức yêu cầu.",
	response_code.MustGreatThanEqual.Code():    "Giá trị phải lớn hơn hoặc bằng mức yêu cầu.",
	response_code.MustLessThan.Code():          "Giá trị phải nhỏ hơn mức yêu cầu.",
	response_code.MustLessThanEqual.Code():     "Giá trị phải nhỏ hơn hoặc bằng mức yêu cầu.",
}

func init() {
	i18n.RegisterMessages(i18n.LocaleEN, SystemMessagesEN)
	i18n.RegisterMessages(i18n.LocaleVI, SystemMessagesVI)
}
