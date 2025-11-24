package utils

import (
    "fmt"
    "reflect"
    "regexp"
    "strings"

    "github.com/go-playground/validator/v10"
)

// Validator 验证器实例
var Validator *validator.Validate

// InitValidator 初始化验证器
func InitValidator() {
    Validator = validator.New()
    
    // 注册自定义验证函数
    Validator.RegisterValidation("phone", validatePhone)
    Validator.RegisterValidation("username", validateUsername)
    Validator.RegisterValidation("password", validatePassword)
    Validator.RegisterValidation("id_card", validateIDCard)
    
    // 注册自定义字段名转换函数
    Validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
        name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
        if name == "-" {
            return ""
        }
        return name
    })
}

// ValidateStruct 验证结构体
func ValidateStruct(s interface{}) error {
    if Validator == nil {
        InitValidator()
    }
    
    err := Validator.Struct(s)
    if err == nil {
        return nil
    }
    
    // 转换验证错误为更友好的格式
    validationErrors := err.(validator.ValidationErrors)
    errorMap := make(map[string]interface{})
    
    for _, e := range validationErrors {
        field := e.Field()
        tag := e.Tag()
        
        var message string
        switch tag {
        case "required":
            message = "字段不能为空"
        case "min":
            message = fmt.Sprintf("长度不能小于 %s", e.Param())
        case "max":
            message = fmt.Sprintf("长度不能大于 %s", e.Param())
        case "email":
            message = "邮箱格式不正确"
        case "len":
            message = fmt.Sprintf("长度必须为 %s", e.Param())
        case "numeric":
            message = "必须为数字"
        case "alpha":
            message = "只能包含字母"
        case "alphanum":
            message = "只能包含字母和数字"
        case "oneof":
            message = fmt.Sprintf("必须为以下值之一: %s", e.Param())
        case "phone":
            message = "手机号格式不正确"
        case "username":
            message = "用户名格式不正确（3-20位字母、数字、下划线）"
        case "password":
            message = "密码格式不正确（至少6位，包含字母和数字）"
        case "id_card":
            message = "身份证号格式不正确"
        default:
            message = fmt.Sprintf("字段验证失败: %s", tag)
        }
        
        errorMap[field] = message
    }
    
    return fmt.Errorf("%v", errorMap)
}

// ValidateVar 验证单个变量
func ValidateVar(field interface{}, tag string) error {
    if Validator == nil {
        InitValidator()
    }
    
    err := Validator.Var(field, tag)
    if err == nil {
        return nil
    }
    
    validationErrors := err.(validator.ValidationErrors)
    if len(validationErrors) > 0 {
        e := validationErrors[0]
        switch e.Tag() {
        case "required":
            return fmt.Errorf("字段不能为空")
        case "min":
            return fmt.Errorf("长度不能小于 %s", e.Param())
        case "max":
            return fmt.Errorf("长度不能大于 %s", e.Param())
        case "email":
            return fmt.Errorf("邮箱格式不正确")
        case "len":
            return fmt.Errorf("长度必须为 %s", e.Param())
        case "numeric":
            return fmt.Errorf("必须为数字")
        case "alpha":
            return fmt.Errorf("只能包含字母")
        case "alphanum":
            return fmt.Errorf("只能包含字母和数字")
        case "oneof":
            return fmt.Errorf("必须为以下值之一: %s", e.Param())
        case "phone":
            return fmt.Errorf("手机号格式不正确")
        case "username":
            return fmt.Errorf("用户名格式不正确（3-20位字母、数字、下划线）")
        case "password":
            return fmt.Errorf("密码格式不正确（至少6位，包含字母和数字）")
        case "id_card":
            return fmt.Errorf("身份证号格式不正确")
        default:
            return fmt.Errorf("字段验证失败: %s", e.Tag())
        }
    }
    
    return err
}

// 自定义验证函数

// validatePhone 验证手机号
func validatePhone(fl validator.FieldLevel) bool {
    phone := fl.Field().String()
    // 简化的手机号验证（实际项目中应该更严格）
    if len(phone) != 11 {
        return false
    }
    return phone[0] == '1'
}

// validateUsername 验证用户名
func validateUsername(fl validator.FieldLevel) bool {
    username := fl.Field().String()
    if len(username) < 3 || len(username) > 20 {
        return false
    }
    
    // 只允许字母、数字、下划线
    for _, char := range username {
        if !((char >= 'a' && char <= 'z') || 
             (char >= 'A' && char <= 'Z') || 
             (char >= '0' && char <= '9') || 
             char == '_') {
            return false
        }
    }
    
    return true
}

// validatePassword 验证密码
func validatePassword(fl validator.FieldLevel) bool {
    password := fl.Field().String()
    if len(password) < 6 {
        return false
    }
    
    hasLetter := false
    hasNumber := false
    
    for _, char := range password {
        if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
            hasLetter = true
        }
        if char >= '0' && char <= '9' {
            hasNumber = true
        }
    }
    
    return hasLetter && hasNumber
}

// validateIDCard 验证身份证号
func validateIDCard(fl validator.FieldLevel) bool {
    idCard := fl.Field().String()
    // 简化的身份证号验证（15位或18位）
    if len(idCard) != 15 && len(idCard) != 18 {
        return false
    }
    
    // 检查是否都是数字或最后一位是X
    for i, char := range idCard {
        if len(idCard) == 18 && i == 17 {
            if !((char >= '0' && char <= '9') || char == 'X' || char == 'x') {
                return false
            }
        } else {
            if char < '0' || char > '9' {
                return false
            }
        }
    }
    
    return true
}

// ValidateEmail 验证邮箱
func ValidateEmail(email string) bool {
    err := ValidateVar(email, "email")
    return err == nil
}

// ValidateRequired 验证必填字段
func ValidateRequired(value interface{}) bool {
    if value == nil {
        return false
    }
    
    if str, ok := value.(string); ok {
        return strings.TrimSpace(str) != ""
    }
    
    return true
}

// ValidateLength 验证长度
func ValidateLength(value string, min, max int) bool {
    length := len(value)
    return length >= min && length <= max
}

// ValidateNumeric 验证数字
func ValidateNumeric(value string) bool {
    err := ValidateVar(value, "numeric")
    return err == nil
}

// ValidateAlpha 验证字母
func ValidateAlpha(value string) bool {
    err := ValidateVar(value, "alpha")
    return err == nil
}

// ValidateAlphanumeric 验证字母数字
func ValidateAlphanumeric(value string) bool {
    err := ValidateVar(value, "alphanum")
    return err == nil
}

// ValidateOneOf 验证枚举值
func ValidateOneOf(value interface{}, allowedValues ...interface{}) bool {
    for _, allowed := range allowedValues {
        if value == allowed {
            return true
        }
    }
    return false
}

// ValidateRange 验证数值范围
func ValidateRange(value interface{}, min, max interface{}) bool {
    switch v := value.(type) {
    case int:
        if minVal, ok := min.(int); ok {
            if v < minVal {
                return false
            }
        }
        if maxVal, ok := max.(int); ok {
            if v > maxVal {
                return false
            }
        }
    case float64:
        if minVal, ok := min.(float64); ok {
            if v < minVal {
                return false
            }
        }
        if maxVal, ok := max.(float64); ok {
            if v > maxVal {
                return false
            }
        }
    default:
        return false
    }
    
    return true
}

// ValidateDate 验证日期格式
func ValidateDate(date string, layout string) bool {
    _, err := time.Parse(layout, date)
    return err == nil
}

// ValidateURL 验证URL格式
func ValidateURL(url string) bool {
    err := ValidateVar(url, "url")
    return err == nil
}

// ValidateUUID 验证UUID格式
func ValidateUUID(uuid string) bool {
    err := ValidateVar(uuid, "uuid")
    return err == nil
}

// ValidateBase64 验证Base64格式
func ValidateBase64(data string) bool {
    err := ValidateVar(data, "base64")
    return err == nil
}

// ValidateIPAddress 验证IP地址
func ValidateIPAddress(ip string) bool {
    err := ValidateVar(ip, "ip")
    return err == nil
}

// ValidateIPv4 验证IPv4地址
func ValidateIPv4(ip string) bool {
    err := ValidateVar(ip, "ipv4")
    return err == nil
}

// ValidateIPv6 验证IPv6地址
func ValidateIPv6(ip string) bool {
    err := ValidateVar(ip, "ipv6")
    return err == nil
}

// ValidateMAC 验证MAC地址
func ValidateMAC(mac string) bool {
    err := ValidateVar(mac, "mac")
    return err == nil
}

// ValidateHostname 验证主机名
func ValidateHostname(hostname string) bool {
    err := ValidateVar(hostname, "hostname")
    return err == nil
}

// ValidateFileExtension 验证文件扩展名
func ValidateFileExtension(filename string, allowedExts ...string) bool {
    if len(allowedExts) == 0 {
        return true
    }
    
    ext := strings.ToLower(filename[strings.LastIndex(filename, "."):])
    for _, allowedExt := range allowedExts {
        if ext == strings.ToLower(allowedExt) {
            return true
        }
    }
    
    return false
}

// ValidateFileSize 验证文件大小
func ValidateFileSize(size int64, maxSize int64) bool {
    return size <= maxSize
}

// ValidateImageType 验证图片类型
func ValidateImageType(contentType string) bool {
    allowedTypes := []string{
        "image/jpeg",
        "image/jpg",
        "image/png",
        "image/gif",
        "image/webp",
    }
    
    for _, allowedType := range allowedTypes {
        if contentType == allowedType {
            return true
        }
    }
    
    return false
}