// generated by codegen/codegen.py, do not edit
/**
 * This module provides the generated definition of `BuiltinNativeObjectType`.
 * INTERNAL: Do not import directly.
 */

private import codeql.swift.generated.Synth
private import codeql.swift.generated.Raw
import codeql.swift.elements.type.internal.BuiltinTypeImpl::Impl as BuiltinTypeImpl

/**
 * INTERNAL: This module contains the fully generated definition of `BuiltinNativeObjectType` and should not
 * be referenced directly.
 */
module Generated {
  /**
   * INTERNAL: Do not reference the `Generated::BuiltinNativeObjectType` class directly.
   * Use the subclass `BuiltinNativeObjectType`, where the following predicates are available.
   */
  class BuiltinNativeObjectType extends Synth::TBuiltinNativeObjectType,
    BuiltinTypeImpl::BuiltinType
  {
    override string getAPrimaryQlClass() { result = "BuiltinNativeObjectType" }
  }
}
