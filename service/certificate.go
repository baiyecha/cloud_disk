package service

import (
	"github.com/baiyecha/cloud_disk/model"
)

type certificateService struct {
	model.CertificateStore
}

func NewCertificateService(cs model.CertificateStore) model.CertificateService {
	return &certificateService{cs}
}
