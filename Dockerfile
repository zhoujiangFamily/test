# Copyright 2019 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Use the offical golang image to create a binary.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.17-buster as builder
ENV SERVICE first-test
# Create and change to the app directory.
# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download
RUN go build
ADD . /$SERVICE
RUN chmod u+x /$SERVICE/run.sh
RUN chmod u+x /$SERVICE/$SERVICE
RUN mkdir -p /var/log/go_log

# Copy local code to the container image.
COPY . ./

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /app/server

# Copy any certificates IF present.
COPY ./certs /app/certs
# Run the web service on container startup.
CMD ["/app/server/run.sh"]