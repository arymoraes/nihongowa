export enum StackConfig {
  VpcName = "BackendVpc",
  VpcSubnetName = "public",
  ClusterName = "BackendCluster",
  TaskDefinitionName = "BackendTaskDefinition",
  DockerImageName = "NihongowaServer",
  ContainerName = "NihongowaContainer",
  ServiceName = "NihongowaService",
}

export enum KeyspacesConfig {
  KeyspaceName = "nihongowa",
  CfnKeyspaceId = "CassandraKeyspace",
}
