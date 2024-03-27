import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import * as ecs from "aws-cdk-lib/aws-ecs";
import * as ecsp from "aws-cdk-lib/aws-ecs-patterns";
import * as ecrAssets from "aws-cdk-lib/aws-ecr-assets";

export class HelloEcsStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    new cdk.aws_cassandra.CfnKeyspace(this, "MyCfnKeyspace", {
      keyspaceName: "nihongowa",
    });

    const dockerImageAsset = new ecrAssets.DockerImageAsset(
      this,
      "MyDockerImage",
      {
        directory: "../api",
      }
    );

    new ecsp.ApplicationLoadBalancedFargateService(this, "MyWebServer", {
      taskImageOptions: {
        image: ecs.ContainerImage.fromDockerImageAsset(dockerImageAsset),
      },
      publicLoadBalancer: true,
    });
  }
}
