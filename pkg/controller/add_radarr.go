package controller

import (
	"github.com/parflesh/radarr-operator/pkg/controller/radarr"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, radarr.Add)
}
