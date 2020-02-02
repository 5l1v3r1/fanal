package image

import (
	"context"
	"io"

	"github.com/containers/image/v5/image"
	"github.com/containers/image/v5/pkg/blobinfocache"
	"github.com/containers/image/v5/pkg/compression"
	"github.com/containers/image/v5/transports/alltransports"
	imageTypes "github.com/containers/image/v5/types"
	"github.com/docker/distribution/reference"
	digest "github.com/opencontainers/go-digest"
	"golang.org/x/xerrors"

	"github.com/aquasecurity/fanal/types"
)

type ImageSource interface {
	GetBlob(ctx context.Context, info imageTypes.BlobInfo, cache imageTypes.BlobInfoCache) (reader io.ReadCloser, n int64, err error)
}

type ImageCloser interface {
	LayerInfos() (layerInfos []imageTypes.BlobInfo)
	ConfigBlob(ctx context.Context) (blob []byte, err error)
	ConfigInfo() imageTypes.BlobInfo
}

type Reference struct {
	Name   string
	IsFile bool
}

type Image struct {
	name          string // image name or tar file name
	blobInfoCache imageTypes.BlobInfoCache
	rawSource     ImageSource
	src           ImageCloser
}

func NewImage(ctx context.Context, image Reference, transports []string, option types.DockerOption) (Image, error) {
	var domain string
	var auth *imageTypes.DockerAuthConfig

	originalName := image.Name
	if !image.IsFile {
		named, err := reference.ParseNormalizedNamed(image.Name)
		if err != nil {
			return Image{}, xerrors.Errorf("invalid image name: %w", err)
		}

		// add 'latest' tag
		named = reference.TagNameOnly(named)
		image.Name = named.String()

		// get a credential for Docker registry
		domain = reference.Domain(named)
		auth = GetToken(ctx, domain, option)
	}

	sys := &imageTypes.SystemContext{
		// TODO: make OSChoice configurable
		OSChoice:                          "linux",
		DockerAuthConfig:                  auth,
		DockerDisableV1Ping:               option.SkipPing,
		DockerInsecureSkipTLSVerify:       imageTypes.NewOptionalBool(option.InsecureSkipTLSVerify),
		OCIInsecureSkipTLSVerify:          option.InsecureSkipTLSVerify,
		DockerDaemonInsecureSkipTLSVerify: option.InsecureSkipTLSVerify,
	}

	rawSource, src, err := newSource(ctx, image.Name, transports, sys)
	if err != nil {
		return Image{}, err
	}

	return Image{
		name:          originalName,
		blobInfoCache: blobinfocache.DefaultCache(sys),
		rawSource:     rawSource,
		src:           src,
	}, nil
}

func newSource(ctx context.Context, imageName string, transports []string, sys *imageTypes.SystemContext) (
	ImageSource, ImageCloser, error) {
	err := xerrors.New("no valid transport")
	for _, transport := range transports {
		imgName := transport + imageName
		var ref imageTypes.ImageReference
		ref, err = alltransports.ParseImageName(imgName)
		if err != nil {
			return nil, nil, xerrors.Errorf("failed to parse an image name: %w", err)
		}

		var rawSource imageTypes.ImageSource
		rawSource, err = ref.NewImageSource(ctx, sys)
		if err != nil {
			// try next transport
			continue
		}

		var src imageTypes.ImageCloser
		src, err = image.FromSource(ctx, sys, rawSource)
		if err != nil {
			return nil, nil, xerrors.Errorf("failed to initialize: %w", err)
		}

		return rawSource, src, nil
	}
	// return only the last error
	return nil, nil, err
}

func (img Image) Name() string {
	return img.name
}

func (img *Image) LayerIDs() []string {
	var layerIDs []string
	for _, l := range img.src.LayerInfos() {
		layerIDs = append(layerIDs, string(l.Digest))
	}
	return layerIDs
}

func (img *Image) ConfigInfo() imageTypes.BlobInfo {
	return img.src.ConfigInfo()
}

func (img *Image) GetBlob(ctx context.Context, dig digest.Digest) (io.ReadCloser, error) {
	rc, _, err := img.rawSource.GetBlob(ctx, imageTypes.BlobInfo{Digest: dig, Size: -1}, img.blobInfoCache)
	if err != nil {
		return nil, xerrors.Errorf("failed to download the layer(%s): %w", dig, err)
	}

	stream, _, err := compression.AutoDecompress(rc)
	if err != nil {
		return nil, xerrors.Errorf("failed to download the layer(%s): %w", dig, err)
	}

	return stream, nil
}
