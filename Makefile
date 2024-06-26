# Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

# http://mozilla.org/MPL/2.0/

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=registry.terraform.io
NAMESPACE=dell
NAME=objectscale
BINARY=terraform-provider-${NAME}
VERSION=1.0.0
OS_ARCH=linux_amd64

default: install

build:
	go mod download
	go build -o ${BINARY}

install: uninstall build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

uninstall:
	rm -rfv ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	find examples -type d -name ".terraform" -exec rm -rfv "{}" +;
	find examples -type f -name "trace.*" -delete
	find examples -type f -name "*.tfstate" -delete
	find examples -type f -name "*.hcl" -delete
	find examples -type f -name "*.backup" -delete
	rm -rf trace.*

no-extract-build: 
	go mod download
	go build -o ${BINARY}

client-build: 
	git clone -b main https://github.com/vangork/objectscale-client.git
	cd ./objectscale-client/golang && cargo build --release

clean:
	rm -f ${BINARY}
	rm terraform-provider-${NAME}_*
	rm -rf ./objectscale-client

release: clean client-build build
	cp terraform-provider-objectscale terraform-provider-${NAME}_v${VERSION}
	zip -j terraform-provider-${NAME}_${VERSION}_${OS_ARCH}.zip terraform-provider-${NAME}_v${VERSION} ./objectscale-client/target/release/libobjectscale_client.so
	cp terraform-registry-manifest.json terraform-provider-${NAME}_${VERSION}_manifest.json
	shasum -a 256 *.zip terraform-provider-${NAME}_${VERSION}_manifest.json > terraform-provider-${NAME}_${VERSION}_SHA256SUMS
