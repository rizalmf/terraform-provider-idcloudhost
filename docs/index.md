# IDCloudHost Provider
This is my personal idcloudhost terraform provider. It allows managing resources within compute and storage resources. 

## Example Usage
### 1. Load providers
```hcl
terraform {
  required_providers {
    idcloudhost = {
      source = "rizalmf/idcloudhost"
      version = "~> 1.0.0"
    }
  }
}
```

### 2. Configure the idcloudhost provider
```hcl
provider "idcloudhost" {
  # you can obtain by create new access as API Token on idcloudhost dashboard
  apikey="XXXXXXXXXXXXXXXXX"

  ## optional. if unset will use your user default location 
  default_location="jkt01" # jkt01(SouthJKT-a), jkt02(NorthJKT-a), jkt03(WestJKT-a), sgp01(Singapore)
}
```

### 3. Create s3 bucket storage
```hcl
# changable field:
# id = STORAGE NAME
# - billing_account_id
resource "idcloudhost_s3" "mybucket" {
  name = "mybucket"
  billing_account_id = 000000
}
```

### 4. Create a VPC network
```hcl
# changable field:
# id = PRIVATE NETWORK UUID
# - name
resource "idcloudhost_private_network" "myprivatenetwork" {
  # network_uuid = <Computed>
  name = "mynetwork"

  # (optional). if unset will use your user default location 
  # this field overwrite "default_location"
  # you can not change location on update
  location = "jkt01" # jkt01(SouthJKT-a), jkt02(NorthJKT-a), jkt03(WestJKT-a), sgp01(Singapore)
  
  lifecycle {
    ignore_changes = [ location ]
  }

}
```

### 5. Create Floating IP
```hcl
# changable field:
# id = FLOAT IP ADDRESS
# - name
# - billing_account_id
resource "idcloudhost_float_ip" "myfloatip" {
  # address = <Computed>
  name = "myfloatnetwork"
  billing_account_id = 000000

  # (optional). if unset will use your user default location 
  # this field overwrite "default_location"
  # you can not change location on update
  location = "jkt01" # jkt01(SouthJKT-a), jkt02(NorthJKT-a), jkt03(WestJKT-a), sgp01(Singapore)
  
  lifecycle {
    ignore_changes = [ location ]
  }

}
```

### 6. Create vm
```hcl
# changable field:
# id = VM UUID
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
  password = "Mypassword1"
  os_name = "ubuntu"
  os_version = "22.04-lts"
  vcpu = 2
  ram = 2048 #mb
  disks = 20 #gb

  # bind existing private network
  private_network_uuid = idcloudhost_private_network.myprivatenetwork.network_uuid 
  
  # (optional) assign to floating ip network. if not, vm doesnt have public ip
  float_ip_address = idcloudhost_float_ip.myfloatip.address

  # add plan will ignored. changing vcpu & ram require desired_status = "stopped"
  desired_status = "stopped" # "stopped", "running"
  
  # (optional). if unset will use your user default location 
  # this field overwrite "default_location"
  # you can not change location on update
  location = "jkt01" # jkt01(SouthJKT-a), jkt02(NorthJKT-a), jkt03(WestJKT-a), sgp01(Singapore)
  
  lifecycle {
    ignore_changes = [ location, os_name, os_version, username, password, billing_account_id ]
  }

}
```