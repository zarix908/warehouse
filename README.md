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
   ```
