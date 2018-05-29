/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package volume

import (
	"github.com/golang/glog"
	"k8s-dev/lib/controller"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/exec"
	"path/filepath"
)

const (
	// are we allowed to set this? else make up our own
	annCreatedBy = "kubernetes.io/createdby"
	createdBy    = "flex-dynamic-provisioner"

	// A PV annotation for the identity of the flexProvisioner that provisioned it
	annProvisionerID = "Provisioner_Id"
	fatherKey        = "share"
)

// NewFlexProvisioner creates a new flex provisioner
func NewFlexProvisioner(client kubernetes.Interface, execCommand, driver string) controller.Provisioner {
	return newFlexProvisionerInternal(client, execCommand, driver)
}

func newFlexProvisionerInternal(client kubernetes.Interface, execCommand, driver string) *flexProvisioner {
	var identity types.UID
	glog.Infof("Driver name is: %s\n", driver)
	provisioner := &flexProvisioner{
		client:      client,
		execCommand: execCommand,
		identity:    identity,
		runner:      exec.New(),
		driver: driver,
	}
	return provisioner
}

type flexProvisioner struct {
	client        kubernetes.Interface
	execCommand   string
	identity      types.UID
	runner        exec.Interface
	driver        string
}

var _ controller.Provisioner = &flexProvisioner{}

// Provision creates a volume i.e. the storage asset and returns a PV object for
// the volume.
func (p *flexProvisioner) Provision(options controller.VolumeOptions) (*v1.PersistentVolume, error) {
	err := p.createVolume(options)
	if err != nil {
		return nil, err
	}

	annotations := make(map[string]string)
	annotations[annCreatedBy] = createdBy

	annotations[annProvisionerID] = string(p.identity)

	flexOptions := map[string]string{}
	flexOptions[optionPVorVolumeName] = options.PVName

	for key, value := range options.Parameters {
		flexOptions[key] = value
	}

	if v, ok := options.Parameters[fatherKey]; ok {
		flexOptions[fatherKey] = filepath.Join(v, options.PVName)
	}

	pv := &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:        options.PVName,
			Labels:      map[string]string{},
			Annotations: annotations,
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: options.PersistentVolumeReclaimPolicy,
			AccessModes:                   options.PVC.Spec.AccessModes,
			Capacity: v1.ResourceList{
				v1.ResourceName(v1.ResourceStorage): options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)],
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				FlexVolume: &v1.FlexPersistentVolumeSource{
					Driver:   p.driver,
					Options:  flexOptions,
					ReadOnly: false,
				},
			},
		},
	}

	return pv, nil
}

func (p *flexProvisioner) createVolume(volumeOptions controller.VolumeOptions) error {
	extraOptions := map[string]string{}
	extraOptions[optionPVorVolumeName] = volumeOptions.PVName

	call := p.NewDriverCall(p.execCommand, provisionCmd)
	call.AppendSpec(volumeOptions.Parameters, extraOptions)
	output, err := call.Run()
	if err != nil {
		glog.Errorf("Failed to create volume %s, output: %s, error: %s", volumeOptions, output.Message, err.Error())
		return err
	}
	return nil
}
