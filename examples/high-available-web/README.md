# Terraform script for 3 Tier Application deployment on G42 Cloud

This is a sample terraform script to automate following Infrastructure components on G42 Cloud to enable provisioning of 3-tier application
 - Security Group
 - ECS
 - RDS
 - ELB
 - EIP
The script contains 2 parts
 - Variable defintion: To define the varibles that will be used in the main terraform script
 - Main terraform script to automate the provisioning for each services 


How to run the script ?

Install terraform from the link - https://developer.hashicorp.com/terraform/downloads 
Copy the script to terraform working directory - https://github.com/g42cloud-terraform/terraform-provider-g42cloud/tree/main/examples/high-available-web 
Initialize the working directory using the command "terraform init"
Verify the script using "terraform plan"
Deploy the Infra using "terraform apply"

 
