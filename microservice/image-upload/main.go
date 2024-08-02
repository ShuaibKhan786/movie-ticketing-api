package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"io"
	"net/url"
	"os"
	"strings"
	"syscall"
	"net"

	_ "image/jpeg"
	_ "image/png"

	pb "github.com/ShuaibKhan786/service/image/upload/pb"
	webp "github.com/chai2010/webp"
	uuid "github.com/google/uuid"
	resize "github.com/nfnt/resize"
	grpc "google.golang.org/grpc"
)
const (
	PROCESSED_IMAGE_DIR = "/home/images"
	IMAGE_EXTENSION 	= "webp"
)


type ImageServiceServer struct {
	pb.UnimplementedImageServiceServer
}

func (s *ImageServiceServer) UploadImage(
	stream pb.ImageService_UploadImageServer,
) error {
	imageBuffer, imageMetaData, err := readAllStreamPayload(stream)
	if err != nil {
		return err
	}

	if !verifyImageFormat(imageMetaData.ImageFormat) {
		return errors.New("unsupported image format")
	}
	

	imageReader := bytes.NewReader(imageBuffer.Bytes())

	imageConfig, _, err := image.DecodeConfig(imageReader)
	if err != nil {
		return err
	}

	if err := resetImageReaderPointer(imageReader); err != nil {
		return errors.New("internal service error")
	}

	decodeImageBuffer, _, err := image.Decode(imageReader)
	if err != nil {
		return err
	}

	var processImage image.Image
	if imageConfig.Width != int(imageMetaData.ImageWidth) ||
	   imageConfig.Height != int(imageMetaData.ImageHeight) {
		processImage = resize.Resize(
			uint(imageMetaData.ImageWidth),
			uint(imageMetaData.ImageHeight),
			decodeImageBuffer,
			resize.Lanczos3,
		)
	}else {
		processImage = decodeImageBuffer
	}

	processImageName := generateUniqueNameUsingUUID()

	processImageUrl := constructImageUrl(processImageName)

	processImageFullPath := fmt.Sprintf("%s/%s.%s", PROCESSED_IMAGE_DIR, processImageName, IMAGE_EXTENSION)

	processImageFile, err := os.OpenFile(
		processImageFullPath,
		syscall.O_WRONLY | syscall.O_CREAT,
		0666);		
	if err != nil {
		return err
	}	
	defer processImageFile.Close()

	if err := webp.Encode(
		processImageFile,
		processImage, 
		&webp.Options{ 
			Lossless: true,
		}); err != nil {
			return err
		}

	stream.SendAndClose(&pb.UploadImageResponse{
		Url: processImageUrl,
	})

	return nil
}

func readAllStreamPayload(stream pb.ImageService_UploadImageServer) (bytes.Buffer, *pb.ImageMetaData, error) {
	var imageBuffer bytes.Buffer
	var imageMetaData *pb.ImageMetaData

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return imageBuffer, nil, err
		}

		switch payload := req.GetPayload().(type) {
		case *pb.UploadImageRequest_MetaData:
			imageMetaData = payload.MetaData
		case *pb.UploadImageRequest_Image:
			if _, err := imageBuffer.Write(payload.Image); err != nil {
				return imageBuffer, nil, err
			}
		default:
			return imageBuffer, nil, errors.New("unsupported types in the payload")
		}
	}

	if imageMetaData == nil {
		return imageBuffer, nil, errors.New("metadata is missing in the payload")
	}

	return imageBuffer, imageMetaData, nil
}


func verifyImageFormat(format string) bool {
	var validFormats = []string{
		"jpeg",
		"jpg",
		"png",
		"webp",
	}

	for _, validFormat := range validFormats {
		if validFormat == format {
			return true
		}
	}

	return false
}

func resetImageReaderPointer(imageReader *bytes.Reader) error {
	_, err := imageReader.Seek(0, io.SeekStart)
	return err
}

func generateUniqueNameUsingUUID() string {
	uuidString := uuid.New().String()
	uuidSliceString := strings.Split(uuidString, "-")
	uuidString = strings.Join(uuidSliceString, "")
	return uuidString
}

func constructImageUrl(imageName string) string {
	urlPath := fmt.Sprintf("/public/static/images/%s.%s", imageName, IMAGE_EXTENSION)
	baseUrl := &url.URL{
		Scheme: "http",
		Host: "localhost:8080",
		Path: urlPath,
	}

	return baseUrl.String()
}


func main() {
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterImageServiceServer(grpcServer, &ImageServiceServer{})

	fmt.Println("gRPC server is running on port :50051")
	if err := grpcServer.Serve(listener); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
		os.Exit(1)
	}
}