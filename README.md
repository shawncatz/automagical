# Automagical

This repo contains the code for a lambda function and terraform code to deploy it.

## What it can do

I've used this code to simplify the use of ASG's and the instances they manage in
a bunch of different ways:

* Resource pools
  * Autoscaling group of worker servers that have a deterministic set of IP's
    from a pool of Elastic IPs. As you scale up, the instances will grab an
    IP from the pool.
  * Guarantee a particular server is up as much as possible and has a
    static IP.
* Stateful services
  * Graphite servers that store their data on secondary EBS volumes and when they
    fall over are automatically recreated and have the volume reattached.
  * Kafka servers that do the same, and when they rejoin the cluster only have
    to worry about what they missed while they were down (minutes)
* DNS (*not yet supported*)
  * Guarantee that all hosts have a name associated to them based on Tag
    conventions (Name, Environment, etc), including reverse records.
  * Support Round Robin records that automatically update when a server is replaced.
  * Associate a CNAME record to the host's default record, for easier to 
    remember - and shorter - names.

## Basics

The basic design is simple. Add tags to an autoscaling group and propagate them
to the instances that it creates. Add matching tags to other resources and the
lambda function will attach them when the instance is running.

This works around the problem that from terraform you cannot attach an elastic IP
to an instance that is built by an autoscaling group.

In its simplest form, you can ensure that an instance comes up with an address
and/or volume reattached after replacement. You can take it further and have
the instances grab resources from pools, to guarantee as you scale up and down
that the new instances use addresses from the pool.

## Code

The code is written in Go and was originally designed and built by me while working
at a previous job. It had become best practice at multiple companies where I worked to
always build instances as part of an autoscaling group (even if you want just
a single instance). This ensured that you were always building things as cattle
not pets, by distancing you from the direct manipulation of instances. However,
it does introduce a problem where you need to ensure that the host has an Elastic
IP or EBS volume attached.

`Automagical` was born from that need.

Previously, it also managed DNS records (forward, reverse, CNAME's, and round robin),
but this will require a bit more design, to make it configurable enough to be
useful to a larger audience that may not architect their DNS infrastructure in the
same way.

## Terraform

To make it as easy as possible to use, I've included Terraform modules that you can
include into your existing Terraform configurations to manage all the related
resources needed for the Automagical lambda function to run.

## Future

I'm not sure how long people will keep using instances in AWS, as many companies
are moving to container-based architecture, but I have a feeling that instances
won't go away entirely.

Please feel free to ping me or open PR's to expand the functionality.
