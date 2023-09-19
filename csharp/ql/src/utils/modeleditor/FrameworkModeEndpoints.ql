/**
 * @name Fetch endpoints for use in the model editor (framework mode)
 * @description A list of endpoints accessible (methods and attributes) for consumers of the library. Excludes test and generated code.
 * @kind problem
 * @problem.severity recommendation
 * @id csharp/utils/modeleditor/framework-mode-endpoints
 * @tags modeleditor endpoints framework-mode
 */

private import csharp
private import dotnet
private import semmle.code.csharp.frameworks.Test
private import ModelEditor

class PublicEndpointFromSource extends Endpoint {
  PublicEndpointFromSource() { this.fromSource() and not this.getFile() instanceof TestFile }
}

from PublicEndpointFromSource endpoint, string apiName, boolean supported, string type
where
  apiName = endpoint.getApiName() and
  supported = isSupported(endpoint) and
  type = supportedType(endpoint)
select endpoint, apiName, supported.toString(), "supported", endpoint.getFile().getBaseName(),
  "library", type, "type", "unknown", "classification"
