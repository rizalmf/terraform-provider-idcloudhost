# Unofficial Terraform Provider for IDCloudHost

This is my personal idcloudhost terraform provider. It allows managing resources within compute and storage resources.


## Usage Example

```hcl
# 1. Load providers
terraform {
  required_providers {
    idcloudhost = {
      source = "rizalmf/idcloudhost"
      version = "=1.0.0"
    }
  }
}

# 2. Configure the idcloudhost provider
provider "idcloudhost" {
  # you can obtain by create new access as API Token on idcloudhost dashboard
  apikey="XXXXXXXXXXXXXXXXX"
}

# 3. Create s3 bucket storage
# changable field:
# - billing_account_id
resource "idcloudhost_s3" "mybucket" {
  name = "mybucket"
  billing_account_id = 000000
}

# 4. Create a VPC network
# changable field:
# - name
resource "idcloudhost_private_network" "myprivatenetwork" {
  # network_uuid = <Computed>
  name = "mynetwork"
}

# 5. Create Floating IP
# changable field:
# - name
# - billing_account_id
resource "idcloudhost_float_ip" "myfloatip" {
  # address = <Computed>
  name = "myfloatnetwork"
  billing_account_id = 000000
}

# 6. Create vm
# changable field:
# - name
# - ram
# - vcpu
# - disks
# - float_ip_address
resource "idcloudhost_vm" "myvm" {
  # uuid = <Computed>
  # disks_uuid = <Computed>
  name = "myvm"
  billing_account_id = 000000
  username = "myusername"
  password = "mypassword"
  os_name = "ubuntu"
  os_version = "22.04"
  vcpu = 2
  ram = 2048 #mb
  disks = 20 #gb

  # bind existing private network
  private_network_uuid = idcloudhost_private_network.myprivatenetwork.network_uuid 
  
  # (optional) assign to floating ip network. if not, vm doesnt have public ip
  float_ip_address = idcloudhost_float_ip.myfloatip.address
}
```


## Next Development
- LB Network
- Specific resource location