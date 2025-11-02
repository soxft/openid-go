package passkey

import (
	"errors"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/soxft/openid-go/app/model"
	"github.com/soxft/openid-go/process/dbutil"
	"gorm.io/gorm"
)

func loadUserPasskeys(userID int) ([]model.PassKey, error) {
	var passkeys []model.PassKey
	if err := dbutil.D.Where("user_id = ?", userID).Find(&passkeys).Error; err != nil {
		return nil, err
	}
	return passkeys, nil
}

func saveCredential(userID int, credential *webauthn.Credential) (*model.PassKey, error) {
	if credential == nil {
		return nil, errors.New("credential is nil")
	}

	encodedID := encodeKey(credential.ID)
	if encodedID == "" {
		return nil, errors.New("credential id is empty")
	}

	encodedPK := encodeKey(credential.PublicKey)
	encodedAAGUID := encodeKey(credential.Authenticator.AAGUID)
	transports := transportsToString(credential.Transport)

	now := time.Now().Unix()
	passkey := model.PassKey{
		UserID:       userID,
		CredentialID: encodedID,
		PublicKey:    encodedPK,
		Attestation:  credential.AttestationType,
		AAGUID:       encodedAAGUID,
		SignCount:    credential.Authenticator.SignCount,
		Transport:    transports,
		CloneWarning: credential.Authenticator.CloneWarning,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	var existing model.PassKey
	err := dbutil.D.Where("user_id = ? AND credential_id = ?", userID, encodedID).Take(&existing).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		if err := dbutil.D.Create(&passkey).Error; err != nil {
			return nil, err
		}
		return &passkey, nil
	case err != nil:
		return nil, err
	default:
		updates := map[string]interface{}{
			"public_key":    encodedPK,
			"attestation":   credential.AttestationType,
			"aaguid":        encodedAAGUID,
			"sign_count":    credential.Authenticator.SignCount,
			"transport":     transports,
			"clone_warning": credential.Authenticator.CloneWarning,
			"updated_at":    time.Now().Unix(),
		}
		if err := dbutil.D.Model(&existing).Updates(updates).Error; err != nil {
			return nil, err
		}
		if err := dbutil.D.Where("id = ?", existing.ID).Take(&existing).Error; err != nil {
			return nil, err
		}
		return &existing, nil
	}
}

func updateCredentialAfterLogin(userID int, credential *webauthn.Credential) (*model.PassKey, error) {
	encodedID := encodeKey(credential.ID)
	if encodedID == "" {
		return nil, errors.New("credential id is empty")
	}

	now := time.Now().Unix()
	updates := map[string]interface{}{
		"sign_count":    credential.Authenticator.SignCount,
		"clone_warning": credential.Authenticator.CloneWarning,
		"transport":     transportsToString(credential.Transport),
		"last_used_at":  now,
		"updated_at":    now,
	}

	if err := dbutil.D.Model(&model.PassKey{}).
		Where("user_id = ? AND credential_id = ?", userID, encodedID).
		Updates(updates).Error; err != nil {
		return nil, err
	}

	var result model.PassKey
	if err := dbutil.D.Where("user_id = ? AND credential_id = ?", userID, encodedID).Take(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func removeCredential(userID, passkeyID int) error {
	res := dbutil.D.Where("user_id = ? AND id = ?", userID, passkeyID).Delete(&model.PassKey{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
