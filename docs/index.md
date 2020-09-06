
# awx Provider

Ansible Tower Provider for handle Tower Projects with [rest](https://docs.ansible.com/ansible-tower/latest/html/towerapi/api_ref.html)

## Example Usage

```hcl
provider "awx" {
  hostname = "http://localhost:8078"
  username = "test"
  password = "changeme"
}
```

## Argument Reference

* List any arguments for the provider **block.**
