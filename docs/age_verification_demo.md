# Age Verification Demo (18+ without revealing exact age/DOB)

## Tổng quan

Demo này minh họa cách sử dụng BBS+ signatures để thực hiện xác minh tuổi một cách bảo vệ quyền riêng tư. Người dùng có thể chứng minh họ trên 18 tuổi mà không cần tiết lộ:
- Tuổi chính xác
- Ngày sinh cụ thể  
- Năm sinh
- Thông tin cá nhân khác

## Kịch bản sử dụng

### Các bên tham gia:
1. **Government Authority (Issuer)**: Cơ quan chính phủ cấp căn cước công dân điện tử
2. **Citizen (Holder)**: Công dân sở hữu căn cước điện tử
3. **Gaming Platform (Verifier)**: Nền tảng game trực tuyến cần xác minh độ tuổi

### Yêu cầu kinh doanh:
- Gaming platform cần xác minh người dùng >= 18 tuổi
- Tuân thủ quy định về nội dung giới hạn độ tuổi
- Bảo vệ quyền riêng tư của người dùng

## Cách thức hoạt động

### 1. Cấp credential với các claims về độ tuổi

Government Authority tạo digital ID credential với nhiều claims:

```go
claims := []vc.Claim{
    // Thông tin cá nhân (sẽ được ẩn)
    {Key: "firstName", Value: "Minh"},
    {Key: "lastName", Value: "Tran Duc"},
    {Key: "dateOfBirth", Value: "1995-03-15"},
    {Key: "address", Value: "456 Le Loi St, District 1, Ho Chi Minh City"},
    
    // Claims xác minh độ tuổi (có thể selective disclosure)
    {Key: "ageOver13", Value: true},
    {Key: "ageOver16", Value: true},
    {Key: "ageOver18", Value: true},
    {Key: "ageOver21", Value: true},
    {Key: "ageOver25", Value: true},
    
    // Thông tin xác thực
    {Key: "nationality", Value: "Vietnamese"},
    {Key: "documentType", Value: "national_id"},
}
```

### 2. Selective Disclosure chỉ tiết lộ thông tin cần thiết

Khi gaming platform yêu cầu xác minh tuổi, citizen chỉ tiết lộ:

```go
selectiveDisclosure := []vc.SelectiveDisclosureRequest{
    {
        CredentialID: credential.ID,
        RevealedAttributes: []string{
            "ageOver18",    // Boolean: true/false cho 18+
            "nationality",  // Cho regional content
            "documentType", // Chứng minh từ chính phủ
        },
    },
}
```

### 3. Gaming platform nhận được thông tin tối thiểu

Gaming platform chỉ biết:
- ✅ User >= 18 tuổi (boolean: true)
- ✅ Nationality: Vietnamese
- ✅ Document từ government authority

Gaming platform KHÔNG biết:
- 🔒 Tuổi chính xác (có thể là 18, 25, 35, 50...)
- 🔒 Ngày sinh
- 🔒 Năm sinh  
- 🔒 Tên thật
- 🔒 Địa chỉ
- 🔒 Số căn cước

## Các tính năng bảo mật

### 1. Zero-Knowledge Proof
- Chứng minh tuổi >= 18 mà không tiết lộ tuổi chính xác
- Sử dụng boolean claims thay vì giá trị số

### 2. Selective Disclosure
- Chỉ tiết lộ attributes cần thiết
- Ẩn các thông tin cá nhân nhạy cảm

### 3. Unlinkable Presentations
- Mỗi presentation có thể được tạo với nonce khác nhau
- Không thể link các presentation với nhau

### 4. Tamper-Evident
- Credential được ký bằng BBS+ signature
- Không thể chỉnh sửa mà không phát hiện

## Ứng dụng thực tế

### Các ngưỡng tuổi khác nhau:
- **Social Media (13+)**: Sử dụng `ageOver13`
- **Movie Theater R-rated (17+)**: Sử dụng `ageOver16` 
- **Alcohol Purchase (21+)**: Sử dụng `ageOver21`
- **Senior Discount (65+)**: Cần thêm `ageOver65` claim

### Các ngành nghề ứng dụng:
1. **Gaming & Entertainment**: Xác minh tuổi cho nội dung giới hạn
2. **E-commerce**: Xác minh mua sản phẩm có giới hạn tuổi
3. **Financial Services**: KYC/AML compliance
4. **Healthcare**: Xác minh độ tuổi cho dịch vụ y tế
5. **Education**: Xác minh độ tuổi cho khóa học

## Chạy demo

```bash
# Build demo
make build-age-demo

# Chạy demo
make age-demo

# Hoặc chạy trực tiếp
go run ./cmd/age_verification_demo
```

## Kết quả mong đợi

Demo sẽ hiển thị:
1. Setup các bên tham gia (Government, Citizen, Gaming Platform)
2. Cấp enhanced credential với age verification claims
3. Tạo privacy-preserving presentation
4. Xác minh thành công tuổi >= 18
5. Demonstration về privacy protection
6. Multiple age verification scenarios

## Lợi ích

### Cho người dùng:
- ✅ Bảo vệ quyền riêng tư
- ✅ Kiểm soát thông tin được chia sẻ
- ✅ Không lo bị tracking qua tuổi chính xác

### Cho doanh nghiệp:
- ✅ Tuân thủ quy định về xác minh tuổi
- ✅ Giảm rủi ro lưu trữ thông tin cá nhân
- ✅ Tăng trust từ người dùng

### Cho xã hội:
- ✅ Privacy-by-design
- ✅ Giảm data breach risks
- ✅ Chuẩn hóa age verification

## Technical Details

- **BBS+ Signatures**: Cho selective disclosure
- **W3C Verifiable Credentials**: Chuẩn credential format
- **DID (Decentralized Identifiers)**: Định danh phi tập trung
- **Go Implementation**: High-performance, enterprise-ready
