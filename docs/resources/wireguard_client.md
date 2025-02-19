---
page_title: "opnsense_wireguard_client Resource - terraform-provider-opnsense"
subcategory: Wireguard
description: |-
  Client resources can be used to setup Wireguard clients.
---

# opnsense_wireguard_client (Resource)

Client resources can be used to setup Wireguard clients.

## Example Usage

```terraform
// Generate random 256-bit base64 public key
resource "random_id" "pubkey" {
  byte_length = 32
}

// Generate random 256-bit base64 private key
resource "random_id" "privkey" {
  byte_length = 32
}

// Configure a peer
resource "opnsense_wireguard_client" "example0" {
  enabled = false
  name = "example0"

  public_key = "/CPjuEdvHJulOIQ56TNyeNHkDJmRCMor4U9k68vMyac="
  psk        = "CJG05xgaLA8RiisoCAmp2U0v329LsIdK1GW4EMc9fmU="

  tunnel_address = [
    "192.168.1.1/32",
    "192.168.4.1/24",
  ]

  server_address = "10.10.10.10"
  server_port    = "1234"
}

// Configure the server
resource "opnsense_wireguard_server" "example0" {
  name = "example0"

  private_key = random_id.privkey.b64_std
  public_key  = random_id.pubkey.b64_std

  dns = [
    "1.1.1.1",
    "8.8.8.8"
  ]

  tunnel_address = [
    "192.168.1.100/32",
    "10.10.0.0/24"
  ]

  peers = [
    opnsense_wireguard_client.example0.id
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the client config.
- `public_key` (String) Public key of this client config. Must be a 256-bit base64 string.
- `tunnel_address` (Set of String) List of addresses allowed to pass trough the tunnel adapter. Please use CIDR notation like `"10.0.0.1/24"`. Defaults to `[]`.

### Optional

- `enabled` (Boolean) Enable this client config. Defaults to `true`.
- `keep_alive` (Number) The persistent keepalive interval in seconds. Defaults to `-1`.
- `psk` (String) Shared secret (PSK) for this peer. You can generate a key using `wg genpsk` on a client with WireGuard installed. Must be a 256-bit base64 string. Defaults to `""`.
- `server_address` (String) The public IP address the endpoint listens to. Defaults to `""`.
- `server_port` (Number) The port the endpoint listens to. Defaults to `-1`.

### Read-Only

- `id` (String) UUID of the client.

