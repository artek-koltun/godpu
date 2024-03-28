// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022-2023 Intel Corporation, or its subsidiaries.
// Copyright (c) 2022-2023 Dell Inc, or its subsidiaries.
// Copyright (c) 2024 Ericsson AB.

// Package network implements the network related CLI commands
package network

import (
	"context"
	"log"
	"time"

	"github.com/opiproject/godpu/network"
	"github.com/spf13/cobra"
)

// CreateVRF Create vrf on OPI Server
func CreateVRF() *cobra.Command {
	var addr string
	var name string
	var vni uint32
	var loopback string
	var vtep string

	cmd := &cobra.Command{
		Use:   "create-vrf",
		Short: "Create a VRF",
		Run: func(_ *cobra.Command, _ []string) {
			var vniparam *uint32
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			evpnClient, err := network.NewVRF(addr)
			if err != nil {
				log.Fatalf("could not create gRPC client: %v", err)
			}
			defer cancel()
			if vni != 0 {
				vniparam = &vni
			}
			vrf, err := evpnClient.CreateVrf(ctx, name, vniparam, loopback, vtep)
			if err != nil {
				log.Fatalf("failed to create vrf: %v", err)
			}
			log.Println("Created VRF:")
			PrintVrf(vrf)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Descriptive name")
	cmd.Flags().Uint32VarP(&vni, "vni", "v", 0, "Must be unique ")
	cmd.Flags().StringVar(&loopback, "loopback", "", "Loopback IP address")
	cmd.Flags().StringVar(&vtep, "vtep", "", "VTEP IP address")
	cmd.Flags().StringVar(&addr, "addr", "localhost:50151", "address of OPI gRPC server")

	if err := cmd.MarkFlagRequired("loopback"); err != nil {
		log.Fatalf("Error marking flag as required: %v", err)
	}
	return cmd
}

// DeleteVRF update the vrf on OPI server
func DeleteVRF() *cobra.Command {
	var addr string
	var name string
	var allowMissing bool

	cmd := &cobra.Command{
		Use:   "delete-vrf",
		Short: "Delete a VRF",
		Run: func(_ *cobra.Command, _ []string) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			evpnClient, err := network.NewVRF(addr)
			if err != nil {
				log.Fatalf("could not create gRPC client: %v", err)
			}
			defer cancel()

			// grpc call to create the bridge port
			_, err = evpnClient.DeleteVrf(ctx, name, allowMissing)
			if err != nil {
				log.Fatalf("DeleteVRF: Error occurred while creating Bridge Port: %q", err)
			}
			log.Printf("Deleted VRF: %s\n", name)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Specify the name of the BridgePort")
	cmd.Flags().BoolVarP(&allowMissing, "allowMissing", "a", false, "Specify the name of the BridgePort")
	cmd.Flags().StringVar(&addr, "addr", "localhost:50151", "address of OPI gRPC server")

	return cmd
}

// GetVRF get vrf details from OPI server
func GetVRF() *cobra.Command {
	var addr string
	var name string

	cmd := &cobra.Command{
		Use:   "get-vrf",
		Short: "Show details of a VRF",
		Run: func(_ *cobra.Command, _ []string) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			evpnClient, err := network.NewVRF(addr)
			if err != nil {
				log.Fatalf("could not create gRPC client: %v", err)
			}
			defer cancel()

			// grpc call to create the bridge port
			vrf, err := evpnClient.GetVrf(ctx, name)
			if err != nil {
				log.Fatalf("DeleteVRF: Error occurred while creating Bridge Port: %q", err)
			}

			log.Println("Get VRF:")
			PrintVrf(vrf)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Specify the name of the vrf")
	cmd.Flags().StringVar(&addr, "addr", "localhost:50151", "address of OPI gRPC server")

	if err := cmd.MarkFlagRequired("name"); err != nil {
		log.Fatalf("Error marking flag as required: %v", err)
	}
	return cmd
}

// ListVRFs list all vrf's with details from OPI server
func ListVRFs() *cobra.Command {
	var addr string
	var pageSize int32
	var pageToken string

	cmd := &cobra.Command{
		Use:   "list-vrfs",
		Short: "Show details of all Vrfs",
		Run: func(_ *cobra.Command, _ []string) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			evpnClient, err := network.NewVRF(addr)
			if err != nil {
				log.Fatalf("could not create gRPC client: %v", err)
			}
			defer cancel()
			for {
				resp, err := evpnClient.ListVrfs(ctx, pageSize, pageToken)
				if err != nil {
					log.Fatalf("Failed to get items: %v", err)
				}
				// Process the server response
				log.Println("list VRFs:")
				for _, vrf := range resp.Vrfs {
					log.Println("VRF with:")
					PrintVrf(vrf)
				}

				// Check if there are more pages to retrieve
				if resp.NextPageToken == "" {
					// No more pages, break the loop
					break
				}
				// Update the page token for the next request
				pageToken = resp.NextPageToken
			}
		},
	}
	cmd.Flags().Int32VarP(&pageSize, "pagesize", "s", 0, "Specify page size")
	cmd.Flags().StringVarP(&pageToken, "pagetoken", "t", "", "Specify the token")
	cmd.Flags().StringVar(&addr, "addr", "localhost:50151", "address of OPI gRPC server")
	return cmd
}

// UpdateVRF update the vrf on OPI server
func UpdateVRF() *cobra.Command {
	var addr string
	var name string
	var updateMask []string
	var allowMissing bool

	cmd := &cobra.Command{
		Use:   "update-vrf",
		Short: "update the VRF",
		Run: func(_ *cobra.Command, _ []string) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			evpnClient, err := network.NewVRF(addr)
			if err != nil {
				log.Fatalf("could not create gRPC client: %v", err)
			}
			defer cancel()

			// grpc call to create the bridge port
			vrf, err := evpnClient.UpdateVrf(ctx, name, updateMask, allowMissing)
			if err != nil {
				log.Fatalf("GetBridgePort: Error occurred while creating Bridge Port: %q", err)
			}
			log.Println("Updated VRF:")
			PrintVrf(vrf)
		},
	}
	cmd.Flags().StringVar(&addr, "addr", "localhost:50151", "address of OPI gRPC server")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Specify the name of the vrf")
	cmd.Flags().StringSliceVar(&updateMask, "update-mask", nil, "update mask")
	cmd.Flags().BoolVarP(&allowMissing, "allowMissing", "a", false, "allow the missing")

	return cmd
}
