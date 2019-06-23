package goapi

import (
	"fmt"

	"github.com/sisatech/goapi/pkg/objects"
)

// VirtualMachine provides access to APIs specific to a single virtual machine.
type VirtualMachine struct {
	mgr  *MachinesManager
	id   string
	name string
}

// VirtualMachineList ..
type VirtualMachineList struct {
	PageInfo objects.PageInfo
	Items    []VirtualMachineListItem
}

// VirtualMachineListItem ..
type VirtualMachineListItem struct {
	Cursor         string
	VirtualMachine VirtualMachine
}

// ID ..
func (v *VirtualMachine) ID() string {
	return v.id
}

// Name ..
func (v *VirtualMachine) Name() string {
	return v.name
}

// Delete the virtual machine.
func (v *VirtualMachine) Delete() error {
	req := v.mgr.environment.newRequest(fmt.Sprintf(`
		mutation {
			deleteVM(id: "%s") {
				id
			}
		}
	`, v.ID()))

	type responseContainer struct {
		VM objects.VM `json:"vm"`
	}
	resp := new(responseContainer)
	return v.mgr.c.graphql.Run(v.mgr.c.ctx, req, &resp)
}

// Image downloads the virtual machine disk image.
func (v *VirtualMachine) Image() error {
	return nil
}

// Pause the virtual machine.
func (v *VirtualMachine) Pause() error {
	req := v.mgr.environment.newRequest(fmt.Sprintf(`
		mutation {
			pauseVM(id: "%s") {
				id
			}
		}
	`, v.ID()))

	type responseContainer struct {
		VM objects.VM `json:"vm"`
	}
	resp := new(responseContainer)
	return v.mgr.c.graphql.Run(v.mgr.c.ctx, req, &resp)
}

// Stop the virtual machine.
func (v *VirtualMachine) Stop() error {

	req := v.mgr.environment.newRequest(fmt.Sprintf(`
		mutation {
			stopVM(id: "%s") {
				id
			}
		}
	`, v.ID()))

	type responseContainer struct {
		VM objects.VM `json:"vm"`
	}
	resp := new(responseContainer)
	return v.mgr.c.graphql.Run(v.mgr.c.ctx, req, &resp)
}

// Start the virtual machine.
func (v *VirtualMachine) Start() error {

	req := v.mgr.environment.newRequest(fmt.Sprintf(`
		mutation {
			startVM(id: "%s") {
				id
			}
		}
	`, v.ID()))

	type responseContainer struct {
		VM objects.VM `json:"vm"`
	}
	resp := new(responseContainer)
	return v.mgr.c.graphql.Run(v.mgr.c.ctx, req, &resp)
}

// Status returns the state of the virtual machine.
func (v *VirtualMachine) Status() (string, error) {

	req := v.mgr.environment.newRequest(fmt.Sprintf(`
		query {
			vm(id: "%s") {
				status
			}
		}
	`, v.ID()))

	type responseContainer struct {
		VM objects.VM `json:"vm"`
	}

	resp := new(responseContainer)
	err := v.mgr.c.graphql.Run(v.mgr.c.ctx, req, &resp)
	if err != nil {
		return "", err
	}

	return resp.VM.Status, nil
}

// Tail the virtual machine serial output.
func (v *VirtualMachine) Tail() error {
	// TODO implement when log-streaming is available
	return nil
}
