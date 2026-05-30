# Warehouse

## Setup

### Prerequisites

Install [age](https://github.com/FiloSottile/age) and [sops](https://github.com/getsops/sops), then install the Ansible collection:

```bash
ansible-galaxy collection install -r ansible/requirements.yml
```

### One-time key setup

1. Generate an age key (keep it outside the repo):
   ```bash
   mkdir -p ~/.config/sops/age
   age-keygen -o ~/.config/sops/age/warehouse.age
   ```

2. Copy the public key printed by `age-keygen` and set it in [.sops.yaml](.sops.yaml):
   ```yaml
   age: age1REPLACE_WITH_YOUR_PUBLIC_KEY
   ```

### Create and encrypt secrets

3. Copy the example, fill in your values, then encrypt in-place:
   ```bash
   cp ansible/secrets.yml.example ansible/secrets.sops.yml
   # edit ansible/secrets.sops.yml with your values
   SOPS_AGE_KEY_FILE=~/.config/sops/age/warehouse.age \
     sops encrypt --in-place ansible/secrets.sops.yml
   ```

   To edit the encrypted file later:
   ```bash
   SOPS_AGE_KEY_FILE=~/.config/sops/age/warehouse.age \
     sops ansible/secrets.sops.yml
   ```

   The encrypted file is safe to commit.

### Run the playbook

4. The `community.sops` vars plugin decrypts secrets automatically at runtime:
   ```bash
   SOPS_AGE_KEY_FILE=~/.config/sops/age/warehouse.age \
   ansible-playbook ansible/setup.yml \
     -i ansible/inventory/hosts.yml \
     --ask-pass \
     --ask-become-pass
   ```

## Roles

| Role | Description |
|------|-------------|
| `configurator` | Deploys the configurator binary and creates its `/etc` directory structure |
| `network` | Installs netplan configs for the configurator and togglable internet-enabled/disabled profiles |
| `aliases` | Installs helper scripts to `/usr/local/bin`: `inet-enable`, `inet-disable`, `main-mount`, `main-umount`, `backup-mount`, `backup-umount` |
| `storage` | Creates LUKS encrypted drive mountpoints and registers them in `/etc/crypttab` and `/etc/fstab` |
| `nfs` | Installs `nfs-kernel-server` and shares `/mnt/main/nfs`; service is disabled and must be started manually |

### Interactions

The `aliases` role is the consumer of all other roles — its scripts wire together the infrastructure each role sets up.

`inet-enable` and `inet-disable` depend on both `configurator` and `network`: they invoke the configurator binary (deployed by `configurator`) with the netplan config directory (deployed by `network`) to swap between the `internet-enabled.yaml` and `internet-disabled.yaml` profiles (also deployed by `network`), then call `netplan apply` to activate the change.

`main-mount`, `backup-mount`, and their umount counterparts depend on `storage`: they open and close the LUKS crypt devices using names registered in `/etc/crypttab` by `storage`, and mount/unmount the decrypted partitions at the paths registered in `/etc/fstab` by `storage`.

## NFS

### Starting the service

The `nfs-server` service is disabled and must be started manually after the storage is mounted:

```bash
sudo systemctl start nfs-server
```

If `/mnt/main/nfs` does not exist the service will refuse to start.

### Mounting on a client

Install the NFS client utilities first (required once per client machine):

```bash
sudo apt install nfs-common
```

**One-time:**

```bash
sudo mount -t nfs <server-ip>:/mnt/main/nfs /mnt/warehouse
```

**Persistent via `/etc/fstab`:**

```
<server-ip>:/mnt/main/nfs  /mnt/warehouse  nfs  nofail,noauto,x-systemd.automount  0  0
```

- `nofail` — boot succeeds even if the server is unreachable
- `noauto,x-systemd.automount` — mounts on first access, not at boot

Apply without rebooting:

```bash
sudo systemctl daemon-reload
sudo mount /mnt/warehouse
```

**Unmount:**

```bash
sudo umount /mnt/warehouse
```
