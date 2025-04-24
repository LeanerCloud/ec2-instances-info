# ec2-instances-info

![Build Status](https://github.com/LeanerCloud/ec2-instances-info/workflows/Test/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/LeanerCloud/ec2-instances-info)](https://goreportcard.com/report/github.com/LeanerCloud/ec2-instances-info)
[![GoDoc](https://godoc.org/github.com/LeanerCloud/ec2-instances-info?status.svg)](http://godoc.org/github.com/LeanerCloud/ec2-instances-info)

Golang library providing specs and pricing information about cloud resources such as:

- AWS EC2 instances
- AWS RDS databases
- AWS ElastiCache clusters
- AWS OpenSearch clusters
- Azure VM instances

It is based on the data that is also powering the comprehensive
[www.ec2instances.info](http://www.ec2instances.info) instance comparison
website, but made easier to consume from Go software.

This code is offered under the public domain/unlicense, but we also offer an API that automates data updates, available for a monthly subscription that helps support ongoing development of this library.

Reach out to us on [Slack](https://join.slack.com/t/leanercloud/shared_invite/zt-xodcoi9j-1IcxNozXx1OW0gh_N08sjg) if you're interested in API access.

## History

This used to be a part of my other project
[AutoSpotting](https://github.com/LeanerCloud/autospotting) which uses it
intensively, but I decided to extract it into a dedicated project since it may be
useful to someone else out there.

Some data fields that were not needed in AutoSpotting may not yet be exposed but
they can be added upon demand.

## Installation or update

You will need Go 1.16 or latest, then it's a matter of installing it as usual using `go get`

```text
go get -u github.com/LeanerCloud/ec2-instances-info/...
```

## Usage

### AWS EC2 Usage

#### One-off usage, with static data

```golang
import "github.com/LeanerCloud/ec2-instances-info"

data, err := ec2instancesinfo.Data() // only needed once

// This would print all the available instance type names:
for _, i := range *data {
  fmt.Println("Instance type", i.InstanceType)
}
```

See the examples directory for a working code example.

#### One-off usage, with updated instance type data

```golang
import "github.com/LeanerCloud/ec2-instances-info"

key := "API_KEY" // API keys are available upon demand from contact@leanercloud.com, free of charge for personal use

err:= ec2instancesinfo.UpdateData(nil, &key);
if err!= nil{
   fmt.Println("Couldn't update instance type data, reverting to static compile-time data", err.Error())
}

data, err := ec2instancesinfo.Data() // needs to be called once after data updates

// This would print all the available instance type names:
for _, i := range *data {
  fmt.Println("Instance type", i.InstanceType)
}
```

#### Continuous usage, with instance type data updated every 2 days

```golang
import "github.com/LeanerCloud/ec2-instances-info"

key := "API_KEY"
go ec2instancesinfo.Updater(2, nil, &key); // use 0 or negative values for weekly updates

data, err := ec2instancesinfo.Data() // only needed once

// This would print all the available instance type names:
for _, i := range *data {
  fmt.Println("Instance type", i.InstanceType)
}
```

### Azure VM Usage

#### Basic usage, with static data

```golang
import "github.com/LeanerCloud/ec2-instances-info"

data, err := ec2instancesinfo.AzureData() // only needed once

// This would print all the available Azure VM instance type names:
for _, i := range *data {
  fmt.Println("Azure VM instance type:", i.InstanceType)
}
```

#### One-off usage, with updated Azure VM data

```golang
import "github.com/LeanerCloud/ec2-instances-info"

key := "API_KEY" // API keys are available upon demand from contact@leanercloud.com

err := ec2instancesinfo.UpdateAzureData(nil, &key)
if err != nil {
   fmt.Println("Couldn't update Azure VM data, reverting to static compile-time data", err.Error())
}

data, err := ec2instancesinfo.AzureData() // needs to be called once after data updates

// This would print all the available Azure VM instance type names:
for _, i := range *data {
  fmt.Println("Azure VM instance type:", i.InstanceType)
}
```

#### Continuous usage, with Azure VM data updated every 2 days

```golang
import "github.com/LeanerCloud/ec2-instances-info"

key := "API_KEY"
go ec2instancesinfo.AzureUpdater(2, nil, &key); // use 0 or negative values for weekly updates

data, err := ec2instancesinfo.AzureData() // only needed once

// This would print all the available Azure VM instance type names:
for _, i := range *data {
  fmt.Println("Azure VM instance type:", i.InstanceType)
}
```

## Contributing

Pull requests and feedback are welcome.

The data can be updated for new instance type coverage by running `make`.

## Try it out on GCP Cloud Shell

- [![Open in Cloud Shell](http://gstatic.com/cloudssh/images/open-btn.svg)](https://ssh.cloud.google.com/cloudshell/editor?cloudshell_git_repo=https://github.com/LeanerCloud/ec2-instances-info.git)

- Click on the `Terminal` then `New Terminal` in the top menu
- In the terminal run `cd ~/cloudshell_open/ec2-instances-info/examples/instances/`
- `go run .` will run the example EC2 code.
- For Azure VM data: `cd ~/cloudshell_open/ec2-instances-info/examples/azure/` and `go run .`