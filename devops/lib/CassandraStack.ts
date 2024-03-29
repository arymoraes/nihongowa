import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import { KeyspacesConfig } from "./config";

export class CassandraStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    new cdk.aws_cassandra.CfnKeyspace(this, KeyspacesConfig.CfnKeyspaceId, {
      keyspaceName: KeyspacesConfig.KeyspaceName,
    });
  }
}
