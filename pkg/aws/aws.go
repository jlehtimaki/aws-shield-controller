package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/shield"
	"github.com/aws/aws-sdk-go/service/sts"
	log "github.com/sirupsen/logrus"
)

func EnableAWSShield(ingressList []string) error {
	awsSession := session.Must(session.NewSession())
	for _, ingress := range ingressList {
		lb, err := getLoadbalancer(awsSession, ingress)
		if err != nil {
			return err
		}
		if *lb.Type == "network" {
			// NLB currently does not support straight enabling by ARN so it needs to be enabled with EIPs
			err := shieldEnableNLB(awsSession, lb)
			if err != nil {
				return err
			}
		} else {
			// Empty function
			err = shieldEnableLB(awsSession, lb)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Gets all active load balancers in the AWS account
func getLoadbalancer(awsSession *session.Session, ingress string) (*elbv2.LoadBalancer, error) {
	svc := elbv2.New(awsSession)
	input := &elbv2.DescribeLoadBalancersInput{}
	result, err := svc.DescribeLoadBalancers(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case elbv2.ErrCodeLoadBalancerNotFoundException:
				log.Error(elbv2.ErrCodeLoadBalancerNotFoundException, aerr.Error())
			default:
				log.Error(aerr.Error())
			}
		} else {
			log.Error(err.Error())
		}
		return nil, err
	}

	for _, loadbalancer := range result.LoadBalancers {
		if *loadbalancer.DNSName == ingress {
			log.Infof("found matching loadbalancer: %s", *loadbalancer.DNSName)
			return loadbalancer, nil
		}
	}
	return nil, fmt.Errorf("did not find matching loadbalancer for ingress: %s", ingress)
}

// Specific LB Shield enablers
func shieldEnableLB(awsSession *session.Session, lb *elbv2.LoadBalancer) error {
	svc := shield.New(awsSession)
	err := enableShield(svc, *lb.LoadBalancerArn, *lb.LoadBalancerName)
	if err != nil {
		return err
	}
	return nil
}

// Loops through NLBs EIP addresses and enables shield to them
// NLB currently does not provide straight Shield enablement
func shieldEnableNLB(awsSession *session.Session, lb *elbv2.LoadBalancer) error {
	svc := shield.New(awsSession)
	for _, l := range lb.AvailabilityZones {
		for _, address := range l.LoadBalancerAddresses {
			eipArn, err := generateEipArn(awsSession, address)
			if err != nil {
				return err
			}
			err = enableShield(svc, eipArn, *address.AllocationId)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Enables the shield on the resources
func enableShield(svc *shield.Shield, resourceArn string, resourceName string) error {
	log.Infof("enabling shield for %s", resourceName)
	input := &shield.CreateProtectionInput{
		Name:        aws.String(resourceName),
		ResourceArn: aws.String(resourceArn),
	}
	result, err := svc.CreateProtection(input)
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case shield.ErrCodeResourceAlreadyExistsException:
			log.Infof("target: %s is already protected", resourceName)
			return nil
		default:
			log.Error(aerr.Error())
			return aerr
		}
	}

	log.Infof("Resource %s is now protected: %s", resourceName, *result.ProtectionId)
	return nil
}

// Generates the EIP ARN because it's not provided by the API
func generateEipArn(awsSession *session.Session, address *elbv2.LoadBalancerAddress) (string, error) {
	svc := sts.New(awsSession)
	input := &sts.GetCallerIdentityInput{}
	result, err := svc.GetCallerIdentity(input)
	if err != nil {
		return "", err
	}
	accountId := *result.Account
	region := *svc.Config.Region
	eipArn := fmt.Sprintf(
		"arn:aws:ec2:%s:%s:eip-allocation/%s", region, accountId, *address.AllocationId,
	)
	return eipArn, nil
}
