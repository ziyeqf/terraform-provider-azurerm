
# Configure Hyper-V to store VMs and disk images in a suitable location:
Set-VMHost -VirtualHardDiskPath C:\Disks -VirtualMachinePath C:\Machines

# Create a Hyper-V switch, type Internal:
New-VMSwitch -Name HyperV-NAT -SwitchType Internal

# Get-NetAdapter to show the interface ID:
$switchIndex=(Get-NetAdapter -Name "vEthernet (HyperV-NAT)").ifIndex

# Configure the switch interface:
New-NetIPAddress -IPAddress 192.168.0.1 -PrefixLength 24 -InterfaceIndex $switchIndex

# Enable NAT:
New-NetNat -Name HyperV-NAT -InternalIPInterfaceAddressPrefix 192.168.0.0/24

# Install DHCP Server role:
Install-WindowsFeature -Name DHCP -IncludeManagementTools

# Create a DHCP scope for guest network:
Add-DhcpServerv4Scope -Name "Hyper-V NAT" -StartRange 192.168.0.100 -EndRange 192.168.0.199 -SubnetMask 255.255.255.0 -LeaseDuration 0.00:59:00

# Set DHCP options for guest network:
Set-DhcpServerv4OptionValue -ScopeId 192.168.0.0 -DnsServer 168.63.129.16 -Router 192.168.0.1

# Allow traffic to/from guests:
New-NetFirewallRule -DisplayName "Allow all guest traffic" -Direction Inbound -RemoteAddress 192.168.0.0/24 -Profile Any -Action Allow

# Create a new type 1 VM, connect to NAT switch and attach the copied VHD, then boot:
New-VM -Name VM1 -Generation 1 -MemoryStartupBytes 16GB -BootDevice VHD -VHDPath C:\Disks\VM1.vhd -SwitchName HyperV-NAT

# Start the VM:
Start-VM -Name VM1
