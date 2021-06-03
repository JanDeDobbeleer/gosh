package gitops

import (
	"errors"
	"gosh/log"
	"os"
)

type Resource interface {
	Exists() bool
	GetFilePath() string
	mapFromKapitanFile(f *kapitanFile)
	mapToKapitanFile() *kapitanFile
	isValid() bool
	getResourceType() string
	getResourceName() string
}

var (
	ResourceDoesNotExistErr  = errors.New("resource does not exist")
	ResourceAlreadyExistsErr = errors.New("resource already exist")
)

func Exists(resource Resource) bool {
	if f, err := os.Stat(resource.GetFilePath()); err == nil && !f.IsDir() {
		return true
	}
	return false
}

func Read(resource Resource) error {
	log.Tracef("Read %s with input: %i", resource)
	if resource == nil || !resource.isValid() {
		return log.Err(ValidationErr, "Invalid %s struct, use NewAppGroup() to create one", resource.getResourceType())
	}
	if !resource.Exists() {
		return log.Errf(ResourceDoesNotExistErr, "The %s '%s' does not exist", resource.getResourceType(), resource.getResourceName())
	}
	if f, err := ReadKapitanFile(resource.GetFilePath()); err == nil {
		resource.mapFromKapitanFile(f)
		log.Tracef("Read %s, result: %i", resource.getResourceType(), resource)
		log.Infof("Read %s '%s'", resource.getResourceType(), resource.getResourceName())
		return nil
	} else {
		return log.Errf(err, "Could not read %s '%s' file", resource.getResourceType(), resource.getResourceName())
	}
}
