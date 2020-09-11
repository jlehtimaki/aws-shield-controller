# aws-shield-controller

`Note! This is early alpha`

Operator to enable AWS Shield on your NLB/ELB/ALB automatically.

# Prerequisites
aws-shield-controller requires certain AWS permissions to run:
```less
  "elasticloadbalancing:DescribeLoadBalancerAttributes",
  "elasticloadbalancing:DescribeSSLPolicies",
  "elasticloadbalancing:DescribeLoadBalancers",
  "elasticloadbalancing:DescribeTargetGroupAttributes",
  "elasticloadbalancing:DescribeListeners",
  "elasticloadbalancing:DescribeTags",
  "elasticloadbalancing:DescribeAccountLimits",
  "elasticloadbalancing:DescribeTargetHealth",
  "elasticloadbalancing:DescribeTargetGroups",
  "elasticloadbalancing:DescribeListenerCertificates",
  "elasticloadbalancing:DescribeRules",
  "shield:CreateProtection",
  "ec2:DescribeAddresses",
```

It's prefered to use IRSA to map your serviceaccount to your IAM role.

# Usage
## Install Controller
`kustomize build kustomize |  kubectl apply -f -`

 Or
 
 `make deploy`
## Enable AWS Shield on your ingress
Add this to your ingress annotation\
`aws.shield.controller: enable`

# Build
`make build`

# TODO
- Allow also to disable Shield
- Check protected resources before trying to enable them again
- Make code more robust
- Make tests