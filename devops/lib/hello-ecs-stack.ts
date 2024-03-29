import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import * as ecs from "aws-cdk-lib/aws-ecs";
import * as ecsp from "aws-cdk-lib/aws-ecs-patterns";
import * as ecrAssets from "aws-cdk-lib/aws-ecr-assets";
import * as dotenv from "dotenv";
import * as cassandra from "aws-cdk-lib/aws-cassandra";

dotenv.config();

export class HelloEcsStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    if (!process.env.OPENAI_API_KEY) {
      throw new Error("OPENAI_API_KEY is required");
    }

    // new cdk.aws_cassandra.CfnKeyspace(this, "MyCfnKeyspace", {
    //   keyspaceName: "nihongowa",
    // });

    const cluster = new ecs.Cluster(this, "BolamosoCluster", {
      clusterName: "BolamosoCluster",
    });

    const dockerImageAsset = new ecrAssets.DockerImageAsset(this, "Bolamoso", {
      directory: "../api",
      buildArgs: {
        OPENAI_API_KEY: process.env.OPENAI_API_KEY || "",
        AWS_REGION: process.env.AWS_REGION || "",
        AWS_ACCESS_KEY_ID: process.env.AWS_ACCESS_KEY_ID || "",
        AWS_SECRET_ACCESS_KEY: process.env.AWS_SECRET_ACCESS_KEY || "",
        ENVIRONMENT: "prod",
      },
    });

    new ecsp.ApplicationLoadBalancedFargateService(this, "Bolamax", {
      taskImageOptions: {
        image: ecs.ContainerImage.fromDockerImageAsset(dockerImageAsset),
      },
      publicLoadBalancer: true,
      cluster,
    });
  }
}
