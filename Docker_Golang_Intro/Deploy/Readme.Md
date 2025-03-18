aws ec2 describe-images --owners amazon --filters "Name=name,Values=amzn2-ami-hvm-*-x86_64-gp2" --query "Images | sort_by(@, &CreationDate) [-1].[ImageId]" --output text


{
    "ImageId": "ami-014e2b14bdb83e8ca",
    "InstanceType": "t3a.nano",
    "KeyName": "kube-demo",
    "MinCount": 1,
    "MaxCount": 1
  }
  

  aws ec2 run-instances --cli-input-json file://instance-config.json

  aws ec2 describe-instances --query 'Reservations[*].Instances[*].[InstanceId,State.Name,PublicIpAddress]' --output table

aws ec2 describe-instances --query 'Reservations[*].Instances[*].[InstanceId,State.Name,PublicIpAddress]' --output table
aws ec2 describe-security-groups --query "SecurityGroups[*].[GroupId,GroupName]" --output table


aws ec2 authorize-security-group-ingress --group-id YOUR_SECURITY_GROUP_ID --protocol tcp --port 22 --cidr YOUR_IP/32

curl https://checkip.amazonaws.com


  chmod 400 kube-demo.pem

  ssh -i kube-demo.pem ec2-user@3.110.231.29
-------------------

sudo yum update -y
sudo yum install docker -y
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker ec2-user
exit
ssh -i kube-demo.pem ec2-user@3.110.231.29



docker save -o my-image.tar my-image-name
scp -i kube-demo.pem docker-test.tar ec2-user@65.2.48.17:/home/ec2-user/

docker run -d  -p 8080:8081 -it docker-test
docker load -i docker-test.tar



-----_FAILED ----image is ARM ARCHITECTURE EC2 is x86


CREATE AN INSTANCE OF ARM64 
Get image first
aws ec2 describe-images --owners amazon --filters "Name=name,Values=*arm64*" --region ap-south-1

DEREGISTER IMAGE ID as its for x86 not arm 64

Delete snapshot as well
aws ec2 describe-images --image-ids ami-014e2b14bdb83e8ca --query "Images[0].BlockDeviceMappings[0].Ebs.SnapshotId" --output text

aws ec2 delete-snapshot --snapshot-id snap-xxxxxxxxxxxxxx


arm 64 image id ami-02455cce11ac4bd45 

 65.2.48.17 


 CREATE A NEW SECURITY GROUP
 aws ec2 describe-vpcs --query "Vpcs[*].[VpcId]" --output table


aws ec2 create-security-group --group-name docker-sg --description "Security group for Docker on EC2" --vpc-id vpc-0138238dd5d53692f
122.174.89.107%
sg-06ef78931e2b30459
 aws ec2 describe-instances --query "Reservations[*].Instances[*].[InstanceId]" --output table i-019d5cf56e349cb0c
aws ec2 modify-instance-attribute --instance-id YOUR_INSTANCE_ID --groups NEW_SG_ID
|  i-019d5cf56e349cb0c 
 aws ec2 describe-instances --instance-id YOUR_INSTANCE_ID --query "Reservations[*].Instances[*].[SecurityGroups]"
aws ec2 authorize-security-group-ingress --group-id sg-06ef78931e2b30459 --protocol tcp --port 8080 --cidr 122.174.89.107/32
