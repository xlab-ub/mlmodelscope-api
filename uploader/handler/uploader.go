package handler

import (
	glog "log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/tus/tusd/pkg/filestore"
	"github.com/tus/tusd/pkg/handler"
	"github.com/tus/tusd/pkg/memorylocker"
	"github.com/unknwon/com"
)

func MakeUploadHandler(uploadPath string, log *logrus.Entry) (http.Handler, error) {
	uploadDir := filepath.Join(uploadPath, "carml_uploads")
	log.Info("using " + uploadDir + " as the upload directory")

	if !com.IsDir(uploadDir) {
		os.MkdirAll(uploadDir, os.FileMode(0755))
	}

	store := filestore.New(uploadDir)
	composer := handler.NewStoreComposer()
	locker := memorylocker.New()
	locker.UseIn(composer)
	store.UseIn(composer)

	handler, err := handler.NewHandler(handler.Config{
		BasePath:                "/",
		StoreComposer:           composer,
		RespectForwardedHeaders: true,
		Logger: glog.New(
			log.Writer(),
			"tusd",
			glog.LstdFlags,
		),
	})
	if err != nil {
		log.WithError(err).Error("Unable to create handler")
		return nil, err
	}
	return handler, nil
}
