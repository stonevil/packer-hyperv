package hyperv

import (
  "strings"
  "github.com/MSOpenTech/packer-hyperv/packer/powershell"
)


func GetVirtualMachineNetworkAdapterAddress(vmName string) (string, error) {

	var script = `
param([string]$vmName, [int]$addressIndex)
try {
  $adapter = Get-VMNetworkAdapter -VMName $vmName -ErrorAction SilentlyContinue
  $ip = $adapter.IPAddresses[$addressIndex]
  if($ip -eq $null) {
    return $false
  }
} catch {
  return $false
}
$ip
`

	var ps powershell.PowerShellCmd
	cmdOut, err := ps.Output(script, vmName, "0");

	return cmdOut, err;
}

func MountDvdDrive(vmName string, path string) error {

  var script = `
param([string]$vmName,[string]$path)
Set-VMDvdDrive -VMName $vmName -Path $path
`

  var ps powershell.PowerShellCmd
  err := ps.Run(script, vmName, path)
  return err
}

func UnmountDvdDrive(vmName string) error {

  var script = `
param([string]$vmName)
Set-VMDvdDrive -VMName $vmName -Path $null
`

  var ps powershell.PowerShellCmd
  err := ps.Run(script, vmName)
  return err
}

func MountFloppyDrive(vmName string, path string) error {
  var script = `
param([string]$vmName, [string]$path)
Set-VMFloppyDiskDrive -VMName $vmName -Path $path
`

  var ps powershell.PowerShellCmd
  err := ps.Run(script, vmName, path)
  return err
}

func UnmountFloppyDrive(vmName string) error {

  var script = `
param([string]$vmName)
Set-VMFloppyDiskDrive -VMName $vmName -Path $null
`

  var ps powershell.PowerShellCmd
  err := ps.Run(script, vmName)
  return err
}

func CreateVirtualMachine(vmName string, path string, ram string, diskSize string, switchName string) error {

  var script = `
param([string]$vmName, [string]$path, [long]$memoryStartupBytes, [long]$newVHDSizeBytes, [string]$switchName)
$vhdx = $vmName + '.vhdx'
$vhdPath = Join-Path -Path $path -ChildPath $vhdx
New-VM -Name $vmName -Path $path -MemoryStartupBytes $memoryStartupBytes -NewVHDPath $vhdPath -NewVHDSizeBytes $newVHDSizeBytes -SwitchName $switchName
`

  var ps powershell.PowerShellCmd
  err := ps.Run(script, vmName, path, ram, diskSize, switchName)
  return err
}


func DeleteVirtualMachine(vmName string) error {

  var script = `
param([string]$vmName)
Remove-VM -Name $vmName -Force
`

  var ps powershell.PowerShellCmd
  err := ps.Run(script, vmName)
  return err
}


func ExportVirtualMachine(vmName string, path string) error {

  var script = `
param([string]$vmName, [string]$path)
Export-VM -Name $vmName -Path $path
`

  var ps powershell.PowerShellCmd
  err := ps.Run(script, vmName, path)
  return err
}

func CopyExportedVirtualMachine(expPath string, outputPath string, vhdDir string, vmDir string) error {

  var script = `
param([string]$srcPath, [string]$dstPath, [string]$vhdDirName, [string]$vmDir)
Copy-Item -Path $srcPath/$vhdDirName -Destination $dstPath -recurse
Copy-Item -Path $srcPath/$vmDir -Destination $dstPath
Copy-Item -Path $srcPath/$vmDir/*.xml -Destination $dstPath/$vmDir
`

  var ps powershell.PowerShellCmd
  err := ps.Run(script, expPath, outputPath, vhdDir, vmDir)
  return err
}

func CreateVirtualSwitch(switchName string, switchType string) (bool,error) {

  var script = `
param([string]$switchName,[string]$switchType)
$switches = Get-VMSwitch -Name $switchName -ErrorAction SilentlyContinue
if ($switches.Count -eq 0) {
  New-VMSwitch -Name $switchName -SwitchType $switchType
  return $true
}
return $false
`

  var ps powershell.PowerShellCmd
  cmdOut, err := ps.Output(script, switchName, switchType)
  var created = strings.TrimSpace(cmdOut) == "True"
  return created, err
}

func DeleteVirtualSwitch(switchName string) error {

  var script = `
param([string]$switchName)
$switch = Get-VMSwitch -Name $switchName -ErrorAction SilentlyContinue
if ($switch -ne $null) {
    $switch | Remove-VMSwitch -Force
}
`

  var ps powershell.PowerShellCmd
  err := ps.Run(script, switchName)
  return err
}


func StartVirtualMachine(vmName string) error {

  var script = `
param([string]$vmName)
Start-VM -Name $vmName
`

  var ps powershell.PowerShellCmd
  err := ps.Run(script, vmName)
  return err
}

func RestartVirtualMachine(vmName string) error {

  var script = `
param([string]$vmName)
Restart-VM $vmName -Force
`

  var ps powershell.PowerShellCmd
  err := ps.Run(script, vmName)
  return err
}

func StopVirtualMachine(vmName string) error {

  var script = `
param([string]$vmName)
$vm = Get-VM -Name $vmName
if ($vm.State -eq [Microsoft.HyperV.PowerShell.VMState]::Running) {
    Stop-VM -VM $vm
}
`

  var ps powershell.PowerShellCmd
  err := ps.Run(script, vmName)
  return err
}


func EnableVirtualMachineIntegrationService(vmName string, integrationServiceName string) error {

  var script = `
param([string]$vmName,[string]$integrationServiceName)
Enable-VMIntegrationService -VMName $vmName -Name $integrationServiceName
`

  var ps powershell.PowerShellCmd
  err := ps.Run(script, vmName, integrationServiceName)
  return err
}


func SetNetworkAdapterVlanId(switchName string, vlanId string) error {

  var script  = `
param([string]$networkAdapterName,[string]$vlanId)
Set-VMNetworkAdapterVlan -ManagementOS -VMNetworkAdapterName $networkAdapterName -Access -VlanId $vlanId
`

  var ps powershell.PowerShellCmd
  err := ps.Run(script, switchName, vlanId)
  return err
}


func SetVirtualMachineVlanId(vmName string, vlanId string) error {

  var script  = `
param([string]$vmName,[string]$vlanId)
Set-VMNetworkAdapterVlan -VMName $vmName -Access -VlanId $vlanId
`
  var ps powershell.PowerShellCmd
  err := ps.Run(script, vmName, vlanId)
  return err
}



func CreateExternalVirtualSwitch(vmName string, switchName string) error {

  var script  = `
param([string]$vmName,[string]$switchName)
$switch = $null
$names = @('ethernet','wi-fi','lan')
$adapters = foreach ($name in $names) {
  Get-NetAdapter -Physical -Name $name -ErrorAction SilentlyContinue | where status -eq 'up'
}

foreach ($adapter in $adapters) { 
  $switch = Get-VMSwitch –SwitchType External | where { $_.NetAdapterInterfaceDescription -eq $adapter.InterfaceDescription }

  if ($switch -eq $null) { 
    $switch = New-VMSwitch -Name $switchName -NetAdapterName $adapter.Name -AllowManagementOS $true -Notes 'Parent OS, VMs, WiFi'
  }

  if ($switch -ne $null) {
    break
  }
}

if($switch -ne $null) { 
  Get-VMNetworkAdapter –VMName $vmName | Connect-VMNetworkAdapter -VMSwitch $switch 
} else { 
  Write-Error 'No internet adapters found'
}
`
  var ps powershell.PowerShellCmd
  err := ps.Run(script, vmName, switchName)
  return err
}



func GetVirtualMachineSwitchName(vmName string) (string, error) {

  var script  = `
param([string]$vmName)
(Get-VMNetworkAdapter -VMName $vmName).SwitchName
`

  var ps powershell.PowerShellCmd
  cmdOut, err := ps.Output(script, vmName)
  if err != nil {
    return "", err
  }

  return strings.TrimSpace(cmdOut), nil
}


func ConnectVirtualMachineNetworkAdapterToSwitch(vmName string, switchName string) error {

  var script  = `
param([string]$vmName,[string]$switchName)
Get-VMNetworkAdapter –VMName $vmName | Connect-VMNetworkAdapter –SwitchName $switchName
`

  var ps powershell.PowerShellCmd
  err := ps.Run(script, vmName, switchName)
  return err
}


func UntagVirtualMachineNetworkAdapterVlan(vmName string, switchName string) error {

  var script  = `
param([string]$vmName,[string]$switchName)
Set-VMNetworkAdapterVlan -VMName $vmName -Untagged
Set-VMNetworkAdapterVlan -ManagementOS -VMNetworkAdapterName $switchName -Untagged
`

  var ps powershell.PowerShellCmd
  err := ps.Run(script, vmName, switchName)
  return err
}