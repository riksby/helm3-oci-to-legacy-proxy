package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/containers/image/docker"
	"github.com/containers/image/pkg/blobinfocache/memory"
)

const (
	imageConfigMediaType = "application/vnd.cncf.helm.config.v1+json"
	imageLayerMediaType  = "application/tar+gzip"
)

func proxy(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	ref, err := docker.ParseReference(fmt.Sprintf("/%s", r.RequestURI))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	img, err := ref.NewImage(ctx, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer img.Close()

	imageConfigInfo := img.ConfigInfo()
	if imageConfigInfo.MediaType != imageConfigMediaType {
		msg := fmt.Sprintf("Expected \"%s\" config media type, got \"%s\" instead", imageConfigMediaType, imageConfigInfo.MediaType)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	layerInfos := img.LayerInfos()
	if layerInfos == nil {
		http.Error(w, "No layers found", http.StatusBadRequest)
		return
	}

	if len(layerInfos) != 1 {
		http.Error(w, "Expected only one layer", http.StatusBadRequest)
		return
	}

	layerInfo := layerInfos[0]

	if layerInfo.MediaType != imageLayerMediaType {
		msg := fmt.Sprintf("Expected \"%s\" layer media type, got \"%s\" instead", imageLayerMediaType, layerInfo.MediaType)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	src, err := ref.NewImageSource(ctx, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if src == nil {
		http.Error(w, "Failed to create new image source", http.StatusBadRequest)
		return
	}

	cache := memory.New()
	layer, _, err := src.GetBlob(ctx, layerInfo, cache)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer layer.Close()

	w.Header().Set("Content-Type", imageLayerMediaType)
	io.Copy(w, layer)
}

func main() {
	log.Fatal(http.ListenAndServe(":80", http.HandlerFunc(proxy)))
}
