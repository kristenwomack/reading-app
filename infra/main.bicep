targetScope = 'subscription'

@minLength(1)
@maxLength(64)
@description('Name of the environment (e.g., dev, staging, prod)')
param environmentName string

@minLength(1)
@description('Primary location for all resources')
param location string

@secure()
@description('Password for admin authentication')
param readingAppPassword string

var abbrs = loadJsonContent('./abbreviations.json')
var resourceToken = toLower(uniqueString(subscription().id, environmentName, location))
var tags = {
  'azd-env-name': environmentName
}

// Resource Group
resource rg 'Microsoft.Resources/resourceGroups@2022-09-01' = {
  name: 'rg-${environmentName}'
  location: location
  tags: tags
}

// Container Apps Environment and App
module app './app/reading-app.bicep' = {
  name: 'reading-app'
  scope: rg
  params: {
    name: '${abbrs.appContainerApps}reading-${resourceToken}'
    location: location
    tags: tags
    containerRegistryName: '${abbrs.containerRegistryRegistries}${resourceToken}'
    readingAppPassword: readingAppPassword
  }
}

output AZURE_CONTAINER_REGISTRY_ENDPOINT string = app.outputs.containerRegistryEndpoint
output AZURE_CONTAINER_REGISTRY_NAME string = app.outputs.containerRegistryName
output SERVICE_API_ENDPOINT_URL string = app.outputs.appUrl
output SERVICE_API_NAME string = app.outputs.appName
