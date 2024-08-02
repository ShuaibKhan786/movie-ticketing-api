package handlers

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
	pb "github.com/ShuaibKhan786/service/image/upload/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


const (
    PORTRAIT_IMAGE_WIDTH   = 452
    PORTRAIT_IMAGE_HEIGHT  = 678

    LANDSCAPE_IMAGE_WIDTH  = 800
    LANDSCAPE_IMAGE_HEIGHT = 450

    PASSPORT_IMAGE_WIDTH   = 312
    PASSPORT_IMAGE_HEIGHT  = 312
	MAX_IMAGE_SIZE 		   = 10 << 20
	IMAGE_FORM_KEY		   = "image"
	IMAGE_LAYOUT_FORM_KEY  = "layout"
	IMAGE_FORMAT_FORM_KEY  = "format"
	CHUNK_SIZE			   = 1024
)


type ImageURL struct {
	ImageUrl string `json:"image-url"`
}


func UploadImage(w http.ResponseWriter, r *http.Request) {
	imageMetaData, imageFile, err := parseAndReadFormData(r)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}
	defer imageFile.Close()

	conn, err := grpc.NewClient(config.Env.GRPC_IMAGE_UPLOAD_SERVER_HOST, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		utils.JSONResponse(&w, "connection"+err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := pb.NewImageServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
	defer cancel()
	stream, err := client.UploadImage(ctx)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := stream.Send(&pb.UploadImageRequest{
		Payload: &pb.UploadImageRequest_MetaData{
			MetaData: imageMetaData,
		},
	}); err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}

	buffer := make([]byte, CHUNK_SIZE)

	for {
		n, err := imageFile.Read(buffer)
		if err == io.EOF {
			break
		}

		if err != nil {
			utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = stream.Send(&pb.UploadImageRequest{
			Payload: &pb.UploadImageRequest_Image{
				Image: buffer[:n],
			},
		})

		if err != nil {
			utils.JSONResponse(&w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	url, err := stream.CloseAndRecv()
	if err != nil {
		utils.JSONResponse(&w, err.Error()+"internal server error", http.StatusInternalServerError)
		return
	}

	if url.Url == "" {
		utils.JSONResponse(&w, "no url internal server error", http.StatusInternalServerError)
		return
	}

	var imageUrl ImageURL
	imageUrl.ImageUrl = url.Url

	jsonImageUrl, err := utils.EncodeJson(&imageUrl)
	if err != nil {
		utils.JSONResponse(&w, err.Error()+"internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonImageUrl)
}


func parseAndReadFormData(r *http.Request) (*pb.ImageMetaData, multipart.File, error){
	var imageMetaData = &pb.ImageMetaData{}
	var file multipart.File 
	var err error

	if err := r.ParseMultipartForm(MAX_IMAGE_SIZE); err != nil {
		return imageMetaData, file, errors.New("unable to parse form data")
	}

	file, _, err = r.FormFile(IMAGE_FORM_KEY)
	if err != nil {
		return imageMetaData, file, errors.New("unable to retrive file")
	}

	imageWidth, imageHeight, err := imageLayout(r.FormValue(IMAGE_LAYOUT_FORM_KEY))
	if err != nil {
		return imageMetaData, file, err
	}

	imageFormat := r.FormValue(IMAGE_FORMAT_FORM_KEY)
	if imageFormat == "" {
		return imageMetaData, file, errors.New("unable to retire image format")
	}
	imageMetaData.ImageFormat = imageFormat
	imageMetaData.ImageWidth = uint32(imageWidth)
	imageMetaData.ImageHeight = uint32(imageHeight)

	return imageMetaData, file, nil
}

func imageLayout(layout string) (int, int, error) {
	switch layout {
	case "portrait":
		return PORTRAIT_IMAGE_WIDTH, PORTRAIT_IMAGE_HEIGHT, nil
	case "landscape":
		return LANDSCAPE_IMAGE_WIDTH, LANDSCAPE_IMAGE_HEIGHT, nil
	case "passport":
		return PASSPORT_IMAGE_WIDTH, PASSPORT_IMAGE_HEIGHT, nil
	default:
		return 0, 0, errors.New("unsupported image layout")
	}
}