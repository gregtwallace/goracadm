# Parent dir is root
$scriptDir = Get-Location
$outDir = Join-Path -Path $scriptDir -ChildPath "/out"

# Windows x64
$env:GOARCH = "amd64"
$env:GOOS = "windows"
go build -o $outDir/goracadm-amd64.exe ./

# Linux x64
$env:GOARCH = "amd64"
$env:GOOS = "linux"
go build -o $outDir/goracadm-amd64-linux ./
