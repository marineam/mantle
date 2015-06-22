// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kola

import (
	"github.com/coreos/mantle/harness"
	"github.com/coreos/mantle/kola/tests/coretest"
	"github.com/coreos/mantle/platform"
)

type KolaSingleNode struct {
	harness.BaseGroup
	c platform.Cluster
	m platform.Machine
}

func (k *KolaSingleNode) Prepare(h *harness.H) {
	var err error
	k.c, err = platform.NewQemuCluster(*QemuImage)
	if err != nil {
		h.Fatalf("QEMU cluster failed: %v", err)
	}

	k.m, err = k.c.NewMachine("")
	if err != nil {
		k.c.Destroy()
		h.Fatalf("QEMU machine failed: %v", err)
	}

	err = scpKolet(k.m)
	if err != nil {
		k.c.Destroy()
		h.Fatalf("scp kolet failed: %v", err)
	}

	h.Log("QEMU instance up")
}

func (k *KolaSingleNode) TestKoletSingleNode(h *harness.H) {
	out, err := k.m.SSH("./kolet single-node")
	if err != nil {
		h.Errorf("kolet failed: %v", err)
	}
	if len(out) != 0 {
		h.Log(string(out))
	}
}

func (k *KolaSingleNode) Cleanup(h *harness.H) {
	if err := k.c.Destroy(); err != nil {
		h.Fatalf("QEMU cluster destroy failed: %v", err)
	}
}

type KoletSingleNode struct {
	harness.BaseGroup
}

func (k *KoletSingleNode) TestCloudinitCloudConfig(h *harness.H) {
	if err := coretest.TestCloudinitCloudConfig(); err != nil {
		h.Errorf("%v", err)
	}
}

func (k *KoletSingleNode) TestCloudinitScript(h *harness.H) {
	if err := coretest.TestCloudinitScript(); err != nil {
		h.Errorf("%v", err)
	}
}

func (k *KoletSingleNode) TestPortSsh(h *harness.H) {
	if err := coretest.TestPortSsh(); err != nil {
		h.Errorf("%v", err)
	}
}

func (k *KoletSingleNode) TestDbusPerms(h *harness.H) {
	if err := coretest.TestDbusPerms(); err != nil {
		h.Errorf("%v", err)
	}
}

func (k *KoletSingleNode) TestSymlinkResolvConf(h *harness.H) {
	if err := coretest.TestSymlinkResolvConf(); err != nil {
		h.Errorf("%v", err)
	}
}

func (k *KoletSingleNode) TestInstalledUpdateEngineRsaKeys(h *harness.H) {
	if err := coretest.TestInstalledUpdateEngineRsaKeys(); err != nil {
		h.Errorf("%v", err)
	}
}

func (k *KoletSingleNode) TestServicesActive(h *harness.H) {
	if err := coretest.TestServicesActive(); err != nil {
		h.Errorf("%v", err)
	}
}

func (k *KoletSingleNode) TestReadOnlyFs(h *harness.H) {
	if err := coretest.TestReadOnlyFs(); err != nil {
		h.Errorf("%v", err)
	}
}
