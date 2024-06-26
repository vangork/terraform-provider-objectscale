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

release:
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64

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
	rm -rf ./objectscale-client
	rm -f ${BINARY}
