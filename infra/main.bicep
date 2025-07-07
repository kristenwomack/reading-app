// Main infrastructure template for Reading Tracker App
targetScope = 'resourceGroup'

metadata description = 'Creates Azure infrastructure for Reading Tracker with Azure Functions and Azure AI (Phi-4)'

@minLength(1)
@maxLength(64)
@description('Name of the environment which is used to generate a short unique hash for resources.')
param environmentName string

@minLength(1)
@description('Primary location for all resources')
param location string

@description('Name of the Azure AI resource')
param aiResourceName string = ''

@description('Location for Azure AI resource (if different from main location)')
param aiLocation string = location

@description('Azure AI SKU name')
param aiSkuName string = 'S0'

// Generate a unique resource token for naming
var resourceToken = toLower(uniqueString(subscription().id, environmentName, location))

// Resource names
var appServicePlanName = 'asp-reading-tracker-${resourceToken}'
var functionAppName = 'func-reading-tracker-${resourceToken}'
var storageAccountName = 'st${replace(resourceToken, '-', '')}'
var logAnalyticsWorkspaceName = 'law-reading-tracker-${resourceToken}'
var applicationInsightsName = 'ai-reading-tracker-${resourceToken}'
var aiServiceName = !empty(aiResourceName) ? aiResourceName : 'ai-reading-tracker-${resourceToken}'

// Tags for all resources
var tags = {
  'azd-env-name': environmentName
  'azd-service-name': 'reading-tracker'
  project: 'reading-tracker'
  environment: environmentName
}

// Log Analytics Workspace
resource logAnalyticsWorkspace 'Microsoft.OperationalInsights/workspaces@2023-09-01' = {
  name: logAnalyticsWorkspaceName
  location: location
  tags: tags
  properties: {
    sku: {
      name: 'PerGB2018'
    }
    retentionInDays: 30
    features: {
      searchVersion: 1
      legacy: 0
      enableLogAccessUsingOnlyResourcePermissions: true
    }
  }
}

// Application Insights
resource applicationInsights 'Microsoft.Insights/components@2020-02-02' = {
  name: applicationInsightsName
  location: location
  tags: tags
  kind: 'web'
  properties: {
    Application_Type: 'web'
    WorkspaceResourceId: logAnalyticsWorkspace.id
  }
}

// Storage Account for Azure Functions
resource storageAccount 'Microsoft.Storage/storageAccounts@2023-05-01' = {
  name: storageAccountName
  location: location
  tags: tags
  sku: {
    name: 'Standard_LRS'
  }
  kind: 'StorageV2'
  properties: {
    supportsHttpsTrafficOnly: true
    encryption: {
      services: {
        file: {
          keyType: 'Account'
          enabled: true
        }
        blob: {
          keyType: 'Account'
          enabled: true
        }
      }
      keySource: 'Microsoft.Storage'
    }
    accessTier: 'Hot'
  }
}

// App Service Plan (Consumption)
resource appServicePlan 'Microsoft.Web/serverfarms@2024-04-01' = {
  name: appServicePlanName
  location: location
  tags: tags
  sku: {
    name: 'Y1'
    tier: 'Dynamic'
    size: 'Y1'
    family: 'Y'
    capacity: 0
  }
  properties: {
    computeMode: 'Dynamic'
  }
}

// Azure AI Service (for Phi-4 and other AI models)
resource aiService 'Microsoft.CognitiveServices/accounts@2024-10-01' = {
  name: aiServiceName
  location: aiLocation
  tags: tags
  sku: {
    name: aiSkuName
  }
  kind: 'AIServices'
  properties: {
    customSubDomainName: aiServiceName
    networkAcls: {
      defaultAction: 'Allow'
      virtualNetworkRules: []
      ipRules: []
    }
    publicNetworkAccess: 'Enabled'
    apiProperties: {
      statisticsEnabled: false
    }
  }
}

// Function App
resource functionApp 'Microsoft.Web/sites@2024-04-01' = {
  name: functionAppName
  location: location
  tags: tags
  kind: 'functionapp'
  identity: {
    type: 'SystemAssigned'
  }
  properties: {
    serverFarmId: appServicePlan.id
    httpsOnly: true
    siteConfig: {
      appSettings: [
        {
          name: 'AzureWebJobsStorage'
          value: 'DefaultEndpointsProtocol=https;AccountName=${storageAccount.name};EndpointSuffix=${environment().suffixes.storage};AccountKey=${storageAccount.listKeys().keys[0].value}'
        }
        {
          name: 'WEBSITE_CONTENTAZUREFILECONNECTIONSTRING'
          value: 'DefaultEndpointsProtocol=https;AccountName=${storageAccount.name};EndpointSuffix=${environment().suffixes.storage};AccountKey=${storageAccount.listKeys().keys[0].value}'
        }
        {
          name: 'WEBSITE_CONTENTSHARE'
          value: toLower(functionAppName)
        }
        {
          name: 'FUNCTIONS_EXTENSION_VERSION'
          value: '~4'
        }
        {
          name: 'WEBSITE_NODE_DEFAULT_VERSION'
          value: '~18'
        }
        {
          name: 'FUNCTIONS_WORKER_RUNTIME'
          value: 'node'
        }
        {
          name: 'APPINSIGHTS_INSTRUMENTATIONKEY'
          value: applicationInsights.properties.InstrumentationKey
        }
        {
          name: 'APPLICATIONINSIGHTS_CONNECTION_STRING'
          value: applicationInsights.properties.ConnectionString
        }
        {
          name: 'AZURE_AI_ENDPOINT'
          value: aiService.properties.endpoint
        }
        {
          name: 'AZURE_AI_API_KEY'
          value: aiService.listKeys().key1
        }
      ]
      cors: {
        allowedOrigins: ['*']
        supportCredentials: false
      }
      ftpsState: 'FtpsOnly'
      minTlsVersion: '1.2'
    }
  }
  dependsOn: [
    applicationInsights
    storageAccount
    aiService
  ]
}

// Outputs for azd and GitHub Actions
output AZURE_LOCATION string = location
output AZURE_TENANT_ID string = tenant().tenantId
output AZURE_SUBSCRIPTION_ID string = subscription().subscriptionId

output AZURE_RESOURCE_GROUP string = resourceGroup().name

output API_BASE_URL string = 'https://${functionApp.properties.defaultHostName}/api'
output API_URI string = 'https://${functionApp.properties.defaultHostName}/api'

output AZURE_FUNCTION_APP_NAME string = functionApp.name
output AZURE_STORAGE_ACCOUNT_NAME string = storageAccount.name
output AZURE_AI_SERVICE_NAME string = aiService.name
output AZURE_AI_ENDPOINT string = aiService.properties.endpoint

output SERVICE_API_IDENTITY_PRINCIPAL_ID string = functionApp.identity.principalId
output AZURE_OPENAI_DEPLOYMENT_NAME string = openAiDeploymentName

output SERVICE_API_IDENTITY_PRINCIPAL_ID string = functionApp.identity.principalId
