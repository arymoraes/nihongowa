import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import * as ecs from "aws-cdk-lib/aws-ecs";
import * as ecrAssets from "aws-cdk-lib/aws-ecr-assets";
import { StackConfig } from "./config";

export class FoundatationStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const vpc = new cdk.aws_ec2.Vpc(this, StackConfig.VpcName, {
      maxAzs: 3,
      natGateways: 0,
      enableDnsHostnames: true,
      enableDnsSupport: true,
      subnetConfiguration: [
        {
          cidrMask: 24,
          name: StackConfig.VpcSubnetName,
          subnetType: cdk.aws_ec2.SubnetType.PUBLIC,
        },
      ],
    });

    new ecs.Cluster(this, StackConfig.ClusterName, {
      clusterName: StackConfig.ClusterName,
      vpc,
    });

    // new ecrAssets.DockerImageAsset(this, StackConfig.DockerImageName, {
    //   directory: "../api",
    //   assetName: StackConfig.DockerImageName,
    // });
  }
}
