package selectel

import (
	"fmt"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v2/volumes"
	"github.com/hashicorp/packer/packer"
	"github.com/mitchellh/multistep"
)

type StepCreateVolume struct {
  ImageId        string
	Name           string
	VolumeType     string
	Size           int
	volume         *volumes.Volume
}

func (s *StepCreateVolume) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	config := state.Get("config").(Config)

	blockStorageClient, err := config.blockStorageV2Client()
	if err != nil {
		err = fmt.Errorf("Error initializing block storage client: %s", err)
		state.Put("error", err)
		return multistep.ActionHalt
	}

	createOpts := &volumes.CreateOpts{
		ImageID:            s.ImageId,
		Name:               s.Name,
		Size:               s.Size,
		VolumeType:         s.VolumeType,
	}

	ui.Say(fmt.Sprintf("[DEBUG] Create Options: %#v", createOpts))
	s.volume, err = volumes.Create(blockStorageClient, createOpts).Extract()
	if err != nil {
		err = fmt.Errorf("Error creating OpenStack volume: %s", err)
		state.Put("error", err)
		return multistep.ActionHalt
	}
	ui.Say(fmt.Sprintf("[INFO] Volume ID: %s", s.volume.ID))

	// Wait for the volume to become available.
	ui.Say(fmt.Sprintf(
		"[DEBUG] Waiting for volume (%s) to become available",
		s.volume.ID))

	stateChange := StateChangeConf{
		Pending:   []string{"downloading", "creating"},
		Target:    []string{"available"},
		Refresh:   VolumeV2StateRefreshFunc(blockStorageClient, s.volume.ID),
		StepState: state,
	}

	_, err = WaitForState(&stateChange)
	if err != nil {
		err := fmt.Errorf("Error waiting for server (%s) to become ready: %s", s.volume.ID, err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("volume_id", s.volume.ID)

	return multistep.ActionContinue
}

func (s *StepCreateVolume) Cleanup(state multistep.StateBag) {
	if s.volume == nil {
		return
	}

	config := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)

	blockStorageClient, err := config.blockStorageV2Client()
	if err != nil {
		err = fmt.Errorf("Error terminating volume, may still be around: %s", err)
		return
	}

	ui.Say(fmt.Sprintf("Terminating the volume: %s ...", s.volume.ID))
	if err := volumes.Delete(blockStorageClient, s.volume.ID).ExtractErr(); err != nil {
		ui.Error(fmt.Sprintf("Error terminating volume, may still be around: %s", err))
		return
	}

	stateChange := StateChangeConf{
		Pending:   []string{"deleting", "downloading", "available"},
		Target:    []string{"deleted"},
		Refresh:   VolumeV2StateRefreshFunc(blockStorageClient, s.volume.ID),
	}

	WaitForState(&stateChange)
}
