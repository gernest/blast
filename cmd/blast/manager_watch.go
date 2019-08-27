// Copyright (c) 2019 Minoru Osuka
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mosuka/blast/manager"
	"github.com/mosuka/blast/protobuf/management"
	"github.com/urfave/cli"
)

func managerWatch(c *cli.Context) error {
	grpcAddr := c.String("grpc-address")

	key := c.Args().Get(0)

	client, err := manager.NewGRPCClient(grpcAddr)
	if err != nil {
		return err
	}
	defer func() {
		err := client.Close()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
	}()

	req := &management.WatchRequest{
		Key: key,
	}
	watchClient, err := client.Watch(req)
	if err != nil {
		return err
	}

	marshaler := manager.JsonMarshaler{}

	for {
		resp, err := watchClient.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err.Error())
			break
		}

		respBytes, err := marshaler.Marshal(resp)
		if err != nil {
			log.Println(err.Error())
			break
		}

		_, _ = fmt.Fprintln(os.Stdout, fmt.Sprintf("%v", string(respBytes)))
	}

	return nil
}
