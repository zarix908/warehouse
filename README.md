# Warehouse

## Setup

1. Copy and fill in secrets:
   ```bash
   cp ansible/secrets.yml.example ansible/secrets.yml
   # edit ansible/secrets.yml with your values
   ```

2. Run the playbook:
   ```bash
   ansible-playbook ansible/setup.yml      
      -i ansible/inventory/hosts.yml 
      -e @ansible/secrets.yml 
      --ask-pass
      --ask-become-pass
   ```

## Roles

| Role | Description |
|------|-------------|
| `configurator` | Deploys the configurator binary and creates its `/etc` directory structure |
| `network` | Installs netplan configs for the configurator and togglable internet-enabled/disabled profiles |
| `aliases` | Installs helper scripts to `/usr/local/bin`: `inet-enable`, `inet-disable`, `main-mount`, `main-umount`, `backup-mount`, `backup-umount` |
| `storage` | Creates LUKS encrypted drive mountpoints and registers them in `/etc/crypttab` and `/etc/fstab` |

### Interactions

The `aliases` role is the consumer of all other roles — its scripts wire together the infrastructure each role sets up.

`inet-enable` and `inet-disable` depend on both `configurator` and `network`: they invoke the configurator binary (deployed by `configurator`) with the netplan config directory (deployed by `network`) to swap between the `internet-enabled.yaml` and `internet-disabled.yaml` profiles (also deployed by `network`), then call `netplan apply` to activate the change.

`main-mount`, `backup-mount`, and their umount counterparts depend on `storage`: they open and close the LUKS crypt devices using names registered in `/etc/crypttab` by `storage`, and mount/unmount the decrypted partitions at the paths registered in `/etc/fstab` by `storage`.
