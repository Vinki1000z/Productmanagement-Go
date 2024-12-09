package imageprocessor

import (
	"context"
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryService(apiKey, apiSecret, cloudName string) (*CloudinaryService, error) {
	url := fmt.Sprintf("cloudinary://%s:%s@%s", apiKey, apiSecret, cloudName)
	cld, err := cloudinary.NewFromURL(url)
	if err != nil {
		return nil, err
	}
	return &CloudinaryService{cld: cld}, nil
}

func (s *CloudinaryService) ProcessImage(ctx context.Context, imageURL string) (string, error) {
	resp, err := s.cld.Upload.Upload(ctx, imageURL, uploader.UploadParams{
		Transformation: "w_800,h_800,c_limit,q_auto:good",
	})
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}
