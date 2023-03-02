// specifies the default hour (UTC) at which tests should be triggered, if enabled
var defaultStartHour = 0

// specifies the default level of parallelism per-service-package
var defaultParallelism = 20

// specifies the default version of Terraform Core which should be used for testing
var defaultTerraformCoreVersion = "1.1.5"

// This represents a cron view of days of the week, Monday - Friday.
const val defaultDaysOfWeek = "2,3,4,5,6"

// Cron value for any day of month
const val defaultDaysOfMonth = "*"

var locations = mapOf(
        "public" to LocationConfiguration("westeurope", "eastus2", "westus2", true)
)

// specifies the list of Azure Environments where tests should be run nightly
var runNightly = mapOf(
        "public" to true
)

// specifies a list of services which should be run with a custom test configuration
var serviceTestConfigurationOverrides = mapOf(

        // Server is only available in certain locations
        "analysisservices" to testConfiguration(locationOverride = LocationConfiguration("westus", "northeurope", "southcentralus", true), useDevTestSubscription = true),

        // App Service Plans for Linux are currently unavailable in WestUS2
        "appservice" to testConfiguration(startHour = 3, daysOfWeek = "2,4,6", locationOverride = LocationConfiguration("westeurope", "westus", "eastus2", true)),

        // these tests all conflict with one another
        "authorization" to testConfiguration(parallelism = 1),

        // HCICluster is only available in certain locations
        "azurestackhci" to testConfiguration(locationOverride = LocationConfiguration("australiaeast", "eastus", "westeurope", true), useDevTestSubscription = true),

        //Blueprints are constrained on the number of targets available - these execute quickly, so can be serialised
        "blueprints" to testConfiguration(parallelism = 1),

        // CDN is only available in certain locations
        "cdn" to testConfiguration(locationOverride = LocationConfiguration("centralus", "eastus2", "westeurope", true), useDevTestSubscription = true),

        // "cognitive" is expensive - Monday, Wednesday, Friday
        // cognitive is only available in certain locations
        "cognitive" to testConfiguration(daysOfWeek = "2,4,6", locationOverride = LocationConfiguration("westeurope", "eastus", "southcentralus", true), useDevTestSubscription = true),

        // Cosmos is only available in certain locations
        "cosmos" to testConfiguration(locationOverride = LocationConfiguration("westus", "northeurope", "southcentralus", true), useDevTestSubscription = true),

        //Confidential Ledger
        "confidentialledger" to testConfiguration(locationOverride = LocationConfiguration("eastus","southcentralus","westeurope", false)),

        // Container App Managed Environments are limited to 20 per location, using 10 as they can take some time to clear
        "containerapps" to testConfiguration(parallelism = 10, locationOverride = LocationConfiguration("westeurope","eastus","canadacentral", false)),

        // The AKS API has a low rate limit
        "containers" to testConfiguration(parallelism = 5, locationOverride = LocationConfiguration("eastus","westeurope","eastus2", false), useDevTestSubscription = true),

        // Custom Providers is only available in certain locations
        "customproviders" to testConfiguration(locationOverride = LocationConfiguration("eastus", "westus2", "westeurope", true)),

        // Dashboard is only available in certain locations
        "dashboard" to testConfiguration(locationOverride = LocationConfiguration("westeurope", "westus2", "eastus2", false)),

        // Datadog is available only in WestUS2 region
        "datadog" to testConfiguration(locationOverride = LocationConfiguration("westus2", "westus2", "centraluseuap", false)),

        // data factory uses NC class VMs which are not available in eastus2
        "datafactory" to testConfiguration(daysOfWeek = "2,4,6", locationOverride = LocationConfiguration("westeurope", "southeastasia", "westus2", false), useDevTestSubscription = true),

        // Data Lake has a low quota
        "datalake" to testConfiguration(parallelism = 2, useDevTestSubscription = true),

        // "hdinsight" is super expensive - G class VM's are not available in westus2, quota only available in westeurope currently
        "hdinsight" to testConfiguration(daysOfWeek = "2,4,6", locationOverride = LocationConfiguration("westeurope", "southeastasia", "eastus2", false), useDevTestSubscription = true),

        // Elastic can't provision many in parallel
        "elastic" to testConfiguration(parallelism = 1, useDevTestSubscription = true),

        // HPC Cache has a 4 instance per subscription quota as of early 2021
        "hpccache" to testConfiguration(parallelism = 3, daysOfWeek = "2,4,6", useDevTestSubscription = true),

        // HSM has low quota and potentially slow recycle time, Only run on Mondays
        "hsm" to testConfiguration(parallelism = 1, daysOfWeek = "1", useDevTestSubscription = true),

        // IoT Central is only available in certain locations
        "iotcentral" to testConfiguration(locationOverride = LocationConfiguration("westeurope", "southeastasia", "eastus2", false), useDevTestSubscription = true),

        // IoT Hub Device Update is only available in certain locations
        "iothub" to testConfiguration(locationOverride = LocationConfiguration("northeurope", "eastus2", "westus2", false), useDevTestSubscription = true),

        // Lab Service is only available in certain locations
        "labservice" to testConfiguration(locationOverride = LocationConfiguration("westeurope", "eastus", "westus", false)),

        // Log Analytics Clusters have a max deployments of 2 - parallelism set to 1 or `importTest` fails
        "loganalytics" to testConfiguration(parallelism = 1, useDevTestSubscription = true),

        // Logic uses app service which is only available in certain locations
        "logic" to testConfiguration(locationOverride = LocationConfiguration("westeurope", "francecentral", "eastus2", false)),

        // Logz is only available in certain locations
        "logz" to testConfiguration(locationOverride = LocationConfiguration("westeurope", "westus2", "eastus2", false)),

        // MSSQl uses app service which is only available in certain locations
        "mssql" to testConfiguration(locationOverride = LocationConfiguration("westeurope", "francecentral", "eastus2", false), useDevTestSubscription = true),

        // MySQL has quota available in certain locations
        "mysql" to testConfiguration(locationOverride = LocationConfiguration("westeurope", "francecentral", "eastus2", false), useDevTestSubscription = true),

        // netapp has a max of 10 accounts and the max capacity of pool is 25 TiB per subscription so lets limit it to 1 to account for broken ones, run Monday, Wednesday, Friday
        "netapp" to testConfiguration(parallelism = 1, daysOfWeek = "2,4,6", locationOverride = LocationConfiguration("westeurope", "eastus2", "westus2", false), useDevTestSubscription = true),

        // Orbital is only available in certain locations
        "orbital" to testConfiguration(locationOverride = LocationConfiguration("eastus", "southcentralus", "westus2", false)),

        "policy" to testConfiguration(useAltSubscription = true),

        // Private DNS Resolver is only available in certain locations
        "privatednsresolver" to testConfiguration(locationOverride = LocationConfiguration("eastus", "westus3", "westeurope", true)),

        // redisenterprise is costly - Monday, Wednesday, Friday
        "redisenterprise" to testConfiguration(daysOfWeek = "2,4,6", useDevTestSubscription = true),

        // servicebus quotas are limited and we experience failures if tests
        // execute too quickly as we run out of namespaces in the sub
        "servicebus" to testConfiguration(parallelism = 10, useDevTestSubscription = true),

        // SignalR only allows provisioning one "Free" instance at a time,
        // which is used in multiple tests
        "signalr" to testConfiguration(parallelism = 1),

        // Spring Cloud only allows a max of 10 provisioned
        "springcloud" to testConfiguration(parallelism = 5, useDevTestSubscription = true),

        // SQL has quota available in certain locations
        "sql" to testConfiguration(locationOverride = LocationConfiguration("westeurope", "francecentral", "eastus2", false), useDevTestSubscription = true),

        // StreamAnalytics has quota available in certain locations
        "streamanalytics" to testConfiguration(locationOverride = LocationConfiguration("westeurope", "francecentral", "eastus2", false), useDevTestSubscription = true),

        // Synapse is only available in certain locations
        "synapse" to testConfiguration(locationOverride = LocationConfiguration("westeurope", "francecentral", "eastus2", false), useDevTestSubscription = true),

        // Currently, we have insufficient quota to actually run these, but there are a few nodes in West Europe, so we'll pin it there for now
        "vmware" to testConfiguration(parallelism = 3, locationOverride = LocationConfiguration("westeurope", "westus2", "eastus2", false), useDevTestSubscription = true),

        // Offset start hour to avoid collision with new App Service, reduce frequency of testing days
        "web" to testConfiguration(startHour = 3, daysOfWeek = "2,4,6", locationOverride = LocationConfiguration("westeurope", "francecentral", "eastus2", false)),

        // moved to alt subsription
        "appconfiguration" to testConfiguration(useDevTestSubscription = true),
        "dns" to testConfiguration(useDevTestSubscription = true),
        "eventgrid" to testConfiguration(useDevTestSubscription = true),
        "eventhub" to testConfiguration(useDevTestSubscription = true),
        "firewall" to testConfiguration(useDevTestSubscription = true),
        "keyvault" to testConfiguration(useDevTestSubscription = true),
        "postgres" to testConfiguration(locationOverride = LocationConfiguration("northeurope", "centralus", "eastus", false), useDevTestSubscription = true),
        "privatedns" to testConfiguration(useDevTestSubscription = true),
        "redis" to testConfiguration(useDevTestSubscription = true)
)
