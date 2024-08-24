# lambadass-2024 : backend

## TODO
- handle errors in framework/lambda/lambda::handleRequest

## Getting started

### 1. terraform/terraform.tfvars
Create the file terraform/terraform.tfvars.

Content should be like this : 
```
region  = "eu-west-3"
profile = "my_profile"
project = "backend"
owner   = "rlapray"
```

Profile is in your ~/.aws/credentials like this :
```
[my_profile]
aws_access_key_id = **********
aws_secret_access_key = **********
```

### 2. terraform/backend.tf
Create file terraform/backend.tf. Content should be like this : 
```
terraform {
  backend "s3" {
    profile = "my_profile"
    bucket  = "terraform-lambadass-2024"
    key     = "backend/terraform.tfstate"
    region  = "eu-west-3"
    dynamodb_table = "terraform-testproject-staging"
  }
}
```
`dynamodb_table` is not required.

### 3. Build tool
The build tool of this project is [Task](https://taskfile.dev/installation/) and its install instructions are [here](https://taskfile.dev/installation/)

Task works with *Taskfile.yml*. 

You can use it *as is* if you use an Archlinux based environment. 
If it's not the case you can change INSTALL, QUERY, SAM_PKG, DOCKER_PKG and TERRAFORM_PKG variables to the ones that works with your environment. Then your system must enable execution of different multiarchitecture containers through QEMU, look at the task *enable_multiarch* if you need to adapt it.

### 3. Run locally
Simply run `go-task run` or `task run`

Some caveats : 
- request id are not working correctly : start / end request id are inconsistent, and inside the program you'll see the same request id use for each call.

### 4. Deploy
Simply run `go-task deploy` or `task deploy`