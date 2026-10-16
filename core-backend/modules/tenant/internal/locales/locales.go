package locales

import (
	"github.com/example/core-backend/modules/tenant/internal/errcode"
	"github.com/example/core-backend/pkg/i18n"
)

var TenantMessagesEN = map[string]string{
	errcode.TenantNameRequired.Code():       "Tenant name is required (minimum 3 characters).",
	errcode.TenantAlreadyExists.Code():      "A tenant with this name already exists.",
	errcode.InvalidIsolateMode.Code():       "Invalid tenant isolation mode specified.",
	errcode.MissingDedicatedDBConfig.Code(): "Database DSN is required when dedicated isolation mode is selected.",
	errcode.NotFoundTenant.Code():           "Tenant not found with the specified ID.",
	errcode.FailToCreateTenant.Code():       "Failed to provision tenant workspace in the database.",
	errcode.FailToGetTenant.Code():          "Failed to retrieve tenant details.",
	errcode.FailToUpdateTenant.Code():       "Failed to update tenant configuration.",
	errcode.FailToDeleteTenant.Code():       "Failed to remove tenant workspace.",
}

var TenantMessagesVI = map[string]string{
	errcode.TenantNameRequired.Code():       "Tên Tenant là bắt buộc (tối thiểu 3 ký tự).",
	errcode.TenantAlreadyExists.Code():      "Tên Tenant này đã tồn tại trong hệ thống.",
	errcode.InvalidIsolateMode.Code():       "Chế độ cô lập (isolation mode) của Tenant không hợp lệ.",
	errcode.MissingDedicatedDBConfig.Code(): "Chế độ Dedicated bắt buộc phải có DSN kết nối cơ sở dữ liệu.",
	errcode.NotFoundTenant.Code():           "Không tìm thấy Tenant với mã định danh đã cung cấp.",
	errcode.FailToCreateTenant.Code():       "Khởi tạo cơ sở dữ liệu cho Tenant thất bại.",
	errcode.FailToGetTenant.Code():          "Truy xuất thông tin Tenant thất bại.",
	errcode.FailToUpdateTenant.Code():       "Cập nhật cấu hình Tenant thất bại.",
	errcode.FailToDeleteTenant.Code():       "Xóa không gian làm việc Tenant thất bại.",
}

func init() {
	i18n.RegisterMessages(i18n.LocaleEN, TenantMessagesEN)
	i18n.RegisterMessages(i18n.LocaleVI, TenantMessagesVI)
}
