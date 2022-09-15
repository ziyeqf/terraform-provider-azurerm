---
subcategory: "HDInsight"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_hdinsight_kafka_cluster"
description: |-
  Manages a HDInsight Kafka Cluster.
---

# azurerm_hdinsight_kafka_cluster

Manages a HDInsight Kafka Cluster.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "hdinsightstor"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "hdinsight"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_hdinsight_kafka_cluster" "example" {
  name                = "example-hdicluster"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  cluster_version     = "4.0"
  tier                = "Standard"

  component_version {
    kafka = "2.1"
  }

  gateway {
    username = "acctestusrgw"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_id = azurerm_storage_container.example.id
    storage_account_key  = azurerm_storage_account.example.primary_access_key
    is_default           = true
  }

  roles {
    head_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }

    worker_node {
      vm_size                  = "Standard_D3_V2"
      username                 = "acctestusrvm"
      password                 = "AccTestvdSC4daf986!"
      number_of_disks_per_node = 3
      target_instance_count    = 3
    }

    zookeeper_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name for this HDInsight Kafka Cluster. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group in which this HDInsight Kafka Cluster should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region which this HDInsight Kafka Cluster should exist. Changing this forces a new resource to be created.

* `cluster_version` - (Required) Specifies the Version of HDInsights which should be used for this Cluster. Changing this forces a new resource to be created.

* `component_version` - (Required) A `component_version` block as defined below.

* `gateway` - (Required) A `gateway` block as defined below.

* `roles` - (Required) A `roles` block as defined below.

* `network` - (Optional) A `network` block as defined below.

* `storage_account` - (Required) One or more `storage_account` block as defined below.

* `storage_account_gen2` - (Required) A `storage_account_gen2` block as defined below.

* `tier` - (Required) Specifies the Tier which should be used for this HDInsight Kafka Cluster. Possible values are `Standard` or `Premium`. Changing this forces a new resource to be created.

* `tls_min_version` - (Optional) The minimal supported TLS version. Possible values are `1.0`, `1.1` or `1.2`. Changing this forces a new resource to be created.

* `encryption_in_transit_enabled` - (Optional) Whether encryption in transit is enabled for this HDInsight Kafka Cluster. Changing this forces a new resource to be created.

~> **NOTE:** Starting on June 30, 2020, Azure HDInsight will enforce TLS 1.2 or later versions for all HTTPS connections. For more information, see [Azure HDInsight TLS 1.2 Enforcement](https://azure.microsoft.com/en-us/updates/azure-hdinsight-tls-12-enforcement/).

---

* `tags` - (Optional) A map of Tags which should be assigned to this HDInsight Kafka Cluster.

* `metastores` - (Optional) A `metastores` block as defined below.

* `monitor` - (Optional) A `monitor` block as defined below.

* `extension` - (Optional) An `extension` block as defined below.

* `rest_proxy` - (Optional) A `rest_proxy` block as defined below.

* `security_profile` - (Optional) A `security_profile` block as defined below.

---

A `component_version` block supports the following:

* `kafka` - (Required) The version of Kafka which should be used for this HDInsight Kafka Cluster. Changing this forces a new resource to be created.

---

A `gateway` block supports the following:

* `password` - (Required) The password used for the Ambari Portal.

-> **NOTE:** This password must be different from the one used for the `head_node`, `worker_node` and `zookeeper_node` roles.

* `username` - (Required) The username used for the Ambari Portal. Changing this forces a new resource to be created.

---

A `head_node` block supports the following:

* `username` - (Required) The Username of the local administrator for the Head Nodes. Changing this forces a new resource to be created.

* `vm_size` - (Required) The Size of the Virtual Machine which should be used as the Head Nodes. Changing this forces a new resource to be created.

* `password` - (Optional) The Password associated with the local administrator for the Head Nodes. Changing this forces a new resource to be created.

-> **NOTE:** If specified, this password must be at least 10 characters in length and must contain at least one digit, one uppercase and one lower case letter, one non-alphanumeric character (except characters ' " ` \).

* `ssh_keys` - (Optional) A list of SSH Keys which should be used for the local administrator on the Head Nodes. Changing this forces a new resource to be created.

-> **NOTE:** Either a `password` or one or more `ssh_keys` must be specified - but not both.

* `subnet_id` - (Optional) The ID of the Subnet within the Virtual Network where the Head Nodes should be provisioned within. Changing this forces a new resource to be created.

* `virtual_network_id` - (Optional) The ID of the Virtual Network where the Head Nodes should be provisioned within. Changing this forces a new resource to be created.

---

A `roles` block supports the following:

* `head_node` - (Required) A `head_node` block as defined above.

* `worker_node` - (Required) A `worker_node` block as defined below.

* `zookeeper_node` - (Required) A `zookeeper_node` block as defined below.

* `kafka_management_node` - (Optional) A `kafka_management_node` block as defined below.

---

A `network` block supports the following:

* `connection_direction` - (Optional) The direction of the resource provider connection. Possible values include `Inbound` or `Outbound`. Defaults to `Inbound`. Changing this forces a new resource to be created.

-> **NOTE:** To enabled the private link the `connection_direction` must be set to `Outbound`.

* `private_link_enabled` - (Optional) Is the private link enabled? Possible values include `True` or `False`. Defaults to `False`. Changing this forces a new resource to be created.

---

A `storage_account` block supports the following:

* `is_default` - (Required) Is this the Default Storage Account for the HDInsight Hadoop Cluster? Changing this forces a new resource to be created.

-> **NOTE:** One of the `storage_account` or `storage_account_gen2` blocks must be marked as the default.

* `storage_account_key` - (Required) The Access Key which should be used to connect to the Storage Account. Changing this forces a new resource to be created.

* `storage_container_id` - (Required) The ID of the Storage Container. Changing this forces a new resource to be created.

-> **NOTE:** This can be obtained from the `id` of the `azurerm_storage_container` resource.

---

A `storage_account_gen2` block supports the following:

* `is_default` - (Required) Is this the Default Storage Account for the HDInsight Hadoop Cluster? Changing this forces a new resource to be created.

-> **NOTE:** One of the `storage_account` or `storage_account_gen2` blocks must be marked as the default.

* `storage_resource_id` - (Required) The ID of the Storage Account. Changing this forces a new resource to be created.

* `filesystem_id` - (Required) The ID of the Gen2 Filesystem. Changing this forces a new resource to be created.

* `managed_identity_resource_id` - (Required) The ID of Managed Identity to use for accessing the Gen2 filesystem. Changing this forces a new resource to be created.

-> **NOTE:** This can be obtained from the `id` of the `azurerm_storage_container` resource.

---

A `worker_node` block supports the following:

* `number_of_disks_per_node` - (Required) The number of Data Disks which should be assigned to each Worker Node, which can be between 1 and 8. Changing this forces a new resource to be created.

* `username` - (Required) The Username of the local administrator for the Worker Nodes. Changing this forces a new resource to be created.

* `vm_size` - (Required) The Size of the Virtual Machine which should be used as the Worker Nodes. Changing this forces a new resource to be created.

* `password` - (Optional) The Password associated with the local administrator for the Worker Nodes. Changing this forces a new resource to be created.

-> **NOTE:** If specified, this password must be at least 10 characters in length and must contain at least one digit, one uppercase and one lower case letter, one non-alphanumeric character (except characters ' " ` \).

* `ssh_keys` - (Optional) A list of SSH Keys which should be used for the local administrator on the Worker Nodes. Changing this forces a new resource to be created.

-> **NOTE:** Either a `password` or one or more `ssh_keys` must be specified - but not both.

* `subnet_id` - (Optional) The ID of the Subnet within the Virtual Network where the Worker Nodes should be provisioned within. Changing this forces a new resource to be created.

* `target_instance_count` - (Optional) The number of instances which should be run for the Worker Nodes.

* `virtual_network_id` - (Optional) The ID of the Virtual Network where the Worker Nodes should be provisioned within. Changing this forces a new resource to be created.

---

A `zookeeper_node` block supports the following:

* `username` - (Required) The Username of the local administrator for the Zookeeper Nodes. Changing this forces a new resource to be created.

* `vm_size` - (Required) The Size of the Virtual Machine which should be used as the Zookeeper Nodes. Changing this forces a new resource to be created.

* `password` - (Optional) The Password associated with the local administrator for the Zookeeper Nodes. Changing this forces a new resource to be created.

-> **NOTE:** If specified, this password must be at least 10 characters in length and must contain at least one digit, one uppercase and one lower case letter, one non-alphanumeric character (except characters ' " ` \).

* `ssh_keys` - (Optional) A list of SSH Keys which should be used for the local administrator on the Zookeeper Nodes. Changing this forces a new resource to be created.

-> **NOTE:** Either a `password` or one or more `ssh_keys` must be specified - but not both.

* `subnet_id` - (Optional) The ID of the Subnet within the Virtual Network where the Zookeeper Nodes should be provisioned within. Changing this forces a new resource to be created.

* `virtual_network_id` - (Optional) The ID of the Virtual Network where the Zookeeper Nodes should be provisioned within. Changing this forces a new resource to be created.

---

A `display_encryption_properties` block supports the following:

* `encryption_algorithm` - (Optional) This is an algorithm identifier for encryption. Possible values are `RSA1_5`, `RSA-OAEP`, `RSA-OAEP-256`.

* `encryption_at_host_enabled` - (Optional) This is indicator to show whether resource disk encryption is enabled.

* `key_vault_key_id` - (Optional) The ID of the key vault key.

* `key_vault_managed_identity_id` - (Optional) This is the resource ID of Managed Identity used to access the key vault.

---

A `kafka_management_node` block supports the following:

* `username` - (Required) The Username of the local administrator for the Kafka Management Nodes. Changing this forces a new resource to be created.

* `vm_size` - (Required) The Size of the Virtual Machine which should be used as the Kafka Management Nodes. Changing this forces a new resource to be created.

* `password` - (Optional) The Password associated with the local administrator for the Kafka Management Nodes. Changing this forces a new resource to be created.

-> **NOTE:** If specified, this password must be at least 10 characters in length and must contain at least one digit, one uppercase and one lower case letter, one non-alphanumeric character (except characters ' " ` \).

* `ssh_keys` - (Optional) A list of SSH Keys which should be used for the local administrator on the Kafka Management Nodes. Changing this forces a new resource to be created.

-> **NOTE:** Either a `password` or one or more `ssh_keys` must be specified - but not both.

* `subnet_id` - (Optional) The ID of the Subnet within the Virtual Network where the Kafka Management Nodes should be provisioned within. Changing this forces a new resource to be created.

* `virtual_network_id` - (Optional) The ID of the Virtual Network where the Kafka Management Nodes should be provisioned within. Changing this forces a new resource to be created.

--- 

A `metastores` block supports the following:

* `hive` - (Optional) A `hive` block as defined below.

* `oozie` - (Optional) An `oozie` block as defined below.

* `ambari` - (Optional) An `ambari` block as defined below.

---

A `hive` block supports the following:

* `server` - (Required) The fully-qualified domain name (FQDN) of the SQL server to use for the external Hive metastore.  Changing this forces a new resource to be created.

* `database_name` - (Required) The external Hive metastore's existing SQL database.  Changing this forces a new resource to be created.

* `username` - (Required) The external Hive metastore's existing SQL server admin username.  Changing this forces a new resource to be created.

* `password` - (Required) The external Hive metastore's existing SQL server admin password.  Changing this forces a new resource to be created.


---

An `oozie` block supports the following:

* `server` - (Required) The fully-qualified domain name (FQDN) of the SQL server to use for the external Oozie metastore.  Changing this forces a new resource to be created.

* `database_name` - (Required) The external Oozie metastore's existing SQL database.  Changing this forces a new resource to be created.

* `username` - (Required) The external Oozie metastore's existing SQL server admin username.  Changing this forces a new resource to be created.

* `password` - (Required) The external Oozie metastore's existing SQL server admin password.  Changing this forces a new resource to be created.

---

An `ambari` block supports the following:

* `server` - (Required) The fully-qualified domain name (FQDN) of the SQL server to use for the external Ambari metastore.  Changing this forces a new resource to be created.

* `database_name` - (Required) The external Hive metastore's existing SQL database.  Changing this forces a new resource to be created.

* `username` - (Required) The external Ambari metastore's existing SQL server admin username.  Changing this forces a new resource to be created.

* `password` - (Required) The external Ambari metastore's existing SQL server admin password.  Changing this forces a new resource to be created.

---

A `monitor` block supports the following:

* `log_analytics_workspace_id` - (Required) The Operations Management Suite (OMS) workspace ID.

* `primary_key` - (Required) The Operations Management Suite (OMS) workspace key.

---

A `extension` block supports the following:

* `log_analytics_workspace_id` - (Required) The workspace ID of the log analytics extension.

* `primary_key` - (Required) The workspace key of the log analytics extension.

---

A `rest_proxy` block supports the following:

* `security_group_id` - (Required) The Azure Active Directory Security Group ID. Changing this forces a new resource to be created.

* `security_group_name` - (Optional) The Azure Active Directory Security Group name. Changing this forces a new resource to be created.

-> **Note:** The `security_group_name` property will be Required in version 3.0 of the AzureRM Provider.

---

A `security_profile` block supports the following:

* `aadds_resource_id` - (Required) The resource ID of the Azure Active Directory Domain Service. Changing this forces a new resource to be created.

* `domain_name` - (Required) The name of the Azure Active Directory Domain. Changing this forces a new resource to be created.

* `domain_username` - (Required) The username of the Azure Active Directory Domain. Changing this forces a new resource to be created.

* `domain_user_password` - (Required) The user password of the Azure Active Directory Domain. Changing this forces a new resource to be created.

* `ldaps_urls` - (Required) A list of the LDAPS URLs to communicate with the Azure Active Directory. Changing this forces a new resource to be created.

* `msi_resource_id` - (Required) The User Assigned Identity for the HDInsight Cluster. Changing this forces a new resource to be created.

* `cluster_users_group_dns` - (Optional) A list of the distinguished names for the cluster user groups. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the HDInsight Kafka Cluster.

* `https_endpoint` - The HTTPS Connectivity Endpoint for this HDInsight Kafka Cluster.

* `kafka_rest_proxy_endpoint` - The Kafka Rest Proxy Endpoint for this HDInsight Kafka Cluster.

* `ssh_endpoint` - The SSH Connectivity Endpoint for this HDInsight Kafka Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Kafka HDInsight Cluster.
* `update` - (Defaults to 60 minutes) Used when updating the Kafka HDInsight Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kafka HDInsight Cluster.
* `delete` - (Defaults to 60 minutes) Used when deleting the Kafka HDInsight Cluster.

## Import

HDInsight Kafka Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_hdinsight_kafka_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.HDInsight/clusters/cluster1
```
