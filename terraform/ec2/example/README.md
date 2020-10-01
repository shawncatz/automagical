# Automagical Example

In this directory, you can find an example of terraform code to use the
automagical module to automagically attach/reattach Elastic IP's and EBS 
volumes to hosts in autoscaling groups. 

This module will build an autoscaling group, an Elastic IP, an EBS volume and
the necessary supporting resources.

## Basics

The basic design is simple. Add tags to an autoscaling group and propagate them
to the instances that they create. Add matching tags to other resources and the
lambda function will attach them when the instance is running.

You MUST always set the `automagical` tag to `true` on any instance to which you
want to automagically attach resources.

Then for addresses, you use the tag `automagical:address` and for EBS volumes,
`automagical:volume`. You will see references in the Terraform code and
in the Go code to `automagical:record`, this is not yet implemented as there
are some deeper design considerations for DNS records.

## Variables

see [MODULE.md](MODULE.md)
