package services

import (
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

type MerchantResponse struct {
	StatusCode int       `json:"statuscode"`
	Message    string    `json:"message"`
	Timestamp  time.Time `json:"timestamp"`
}
type SuccessResponse struct {
	MerchantResponse
	Data Data
}

type Data struct {
	FaceId       string `json:"face_id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Pin          string `json:"pin"`
	RoleName     string `json:"role_name"`
	URLFrontFace string `json:"url_front_face"`
	LoginToken   string `json:"login_token"`
	RegDaftar    bool   `json:"reg_daftar_usaha"`
	VerifyKYC    bool   `json:"verif_kyc"`
}

// ThirdPartyService -> struct
type ThirdPartyService struct {
	logger infrastructure.Logger
	env    infrastructure.Env
}

// NewThirdPartyService -> creates a new ThirdPartyService
func NewThirdPartyService(
	logger infrastructure.Logger,
	env infrastructure.Env,
) ThirdPartyService {
	return ThirdPartyService{
		logger: logger,
		env:    env,
	}
}

func (m ThirdPartyService) MerchantRegister(data models.MerchanntRegisterInput) (*SuccessResponse, error) {
	// url := fmt.Sprintf("%s/Accounts/%s/Messages.json", m.MerchantRegisterUrl)
	m.logger.Zap.Info(m.env.MerchantRegisterUrl)

	method := "POST"

	// dataBody, err := json.Marshal(data)

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("name", data.Name)
	_ = writer.WriteField("gender", data.Gender)
	_ = writer.WriteField("birth_date", data.BirthDate)
	_ = writer.WriteField("pin", data.Pin)
	_ = writer.WriteField("email", data.Email)
	_ = writer.WriteField("password", data.Password)
	_ = writer.WriteField("device_token", data.DeviceToken)
	m.logger.Zap.Info("birth date", data.BirthDate)
	m.logger.Zap.Info("full Name", data.Name)
	err := writer.Close()
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, m.env.MerchantRegisterUrl, payload)

	if err != nil {
		return nil, err
	}
	m.logger.Zap.Info("request body", req.Body)

	// token := fmt.Sprintf("Basic %s", m.getBasicToken())
	// m.logger.Zap.Info(token)

	// req.Header.Add("Authorization", token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusCreated {
		result := MerchantResponse{}
		if err := json.Unmarshal(body, &result); err != nil {
			m.logger.Zap.Error("error", err)
			return nil, err
		}
		m.logger.Zap.Info("result in http status code", result)
		return nil, err
	}

	result := SuccessResponse{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	m.logger.Zap.Info("utlimate repsonse ", result)

	return &result, nil

}
