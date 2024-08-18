# Unofficial Terraform Provider for IDCloudHost

This is my personal idcloudhost terraform provider. It allows managing resources within compute and storage resources.


## Example Usage
### 1. Load providers
```hcl
terraform {
  required_providers {
    idcloudhost = {
      source = "rizalmf/idcloudhost"
      version = "~> 1.2.0"
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
# id = STORAGE NAME
# changable field:
# - billing_account_id
resource "idcloudhost_s3" "mybucket" {
  name = "mybucket"
  billing_account_id = 000000
}
```

### 4. Create a VPC network
```hcl
# id = PRIVATE NETWORK UUID
# changable field:
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
# id = FLOAT IP ADDRESS
# changable field:
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

### 6. Create Vm
```hcl
# id = VM UUID
# changable field:
# - name
# - ram
# - vcpu
# - disks
# - float_ip_address
# - desired_status
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

### 7. Create Load Balancer (Vms target only)
```hcl
# prepare a local values
locals {
  # forwarding rules
  lb_rules = [
    { source_port = 8080, target_port=80 }
  ]

  # target servers
  lb_targets = [
    { uuid="XXXX-XXXXXXX-XXXXXXX-XXXXXX" }, # uuid of vm1
    { uuid="XXXX-XXXXXXX-XXXXXXX-XXXXXX" }, # uuid of vm2
    { uuid="XXXX-XXXXXXX-XXXXXXX-XXXXXX" }, # uuid of vm3
  ]
}

# id = UUID
# changable field:
# - name
# - billing_account_id
# - targets
# - rules
resource "idcloudhost_loadbalancer" "mylb" {
  # uuid = <Computed>
  name = "myloadbalancer"
  billing_account_id = 000000
  # bind existing private network
  private_network_uuid = idcloudhost_private_network.myprivatenetwork.network_uuid

  # (optional) assign to floating ip network. if not, lb doesnt have public ip
  float_ip_address = idcloudhost_float_ip.anotherfloatip.address

  # (optional). if unset will use your user default location 
  # this field overwrite "default_location"
  # you can not change location on update
  location = "jkt01" # jkt01(SouthJKT-a), jkt02(NorthJKT-a), jkt03(WestJKT-a), sgp01(Singapore)

  dynamic "rules" {
    for_each = locals.lb_rules
    content {
      source_port = rules.value.source_port
      target_port = rules.value.target_port
    }
  }

  dynamic "targets" {
    for_each = locals.lb_targets
    content {
      target_uuid = targets.value.target_uuid
      target_type = targets.value.target_type
    }
  }
  
  lifecycle {
    ignore_changes = [ location, private_network_uuid ]
  }

}
```

## Next Development
- Resource LB Network(Load Balancer)
- ✅ Resource VM add desired_status (v1.2.0)
- ✅ Specific resource location (v1.1.0)
- ✅ Support terraform import (v1.1.0)
- ✅ Docs (v1.0.6)
- ✅ Github Action workflow (1.0.5)
- ✅ Resource Float IP (v1.0.4)
- ✅ Resource VM (v1.0.3)
- ✅ Resource Private Network (v1.0.3)
- ✅ Resource S3 (v1.0.3)