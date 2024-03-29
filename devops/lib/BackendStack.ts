import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import * as ecs from "aws-cdk-lib/aws-ecs";
import * as ecsp from "aws-cdk-lib/aws-ecs-patterns";
import * as dotenv from "dotenv";
import { StackConfig } from "./config";

dotenv.config();

enum BackendConfig {
  TaskDefinitionName = "BackendTaskDefinition",
}

export class BackendStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    if (!process.env.OPENAI_API_KEY) {
      throw new Error("OPENAI_API_KEY is required");
    }

    const taskDefinition = new cdk.aws_ecs.FargateTaskDefinition(
      this,
      BackendConfig.TaskDefinitionName
    );

    taskDefinition.addContainer(StackConfig.ContainerName, {
      containerName: StackConfig.ContainerName,
      environment: {
        OPENAI_API_KEY: process.env.OPENAI_API_KEY || "",
        AWS_REGION: process.env.AWS_REGION || "",
        AWS_ACCESS_KEY_ID: process.env.AWS_ACCESS_KEY_ID || "",
        AWS_SECRET_ACCESS_KEY: process.env.AWS_SECRET_ACCESS_KEY || "",
        ENVIRONMENT: "prod",
      },
      image: ecs.ContainerImage.fromDockerImageAsset(
        new cdk.aws_ecr_assets.DockerImageAsset(
          this,
          StackConfig.DockerImageName,
          {
            directory: "../api",
          }
        )
      ),
      logging: ecs.LogDrivers.awsLogs({
        streamPrefix: StackConfig.ContainerName,
      }),
      portMappings: [
        {
          containerPort: 1323,
          hostPort: 1323,
          protocol: ecs.Protocol.TCP,
        },
      ],
    });

    new ecsp.ApplicationLoadBalancedFargateService(
      this,
      StackConfig.ServiceName,
      {
        taskDefinition,
        serviceName: StackConfig.ServiceName,
        publicLoadBalancer: true,
      }
    );
  }
}
