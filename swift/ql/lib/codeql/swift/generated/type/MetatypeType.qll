// generated by codegen/codegen.py, do not edit
/**
 * This module provides the generated definition of `MetatypeType`.
 * INTERNAL: Do not import directly.
 */

private import codeql.swift.generated.Synth
private import codeql.swift.generated.Raw
import codeql.swift.elements.type.internal.AnyMetatypeTypeImpl::Impl as AnyMetatypeTypeImpl

/**
 * INTERNAL: This module contains the fully generated definition of `MetatypeType` and should not
 * be referenced directly.
 */
module Generated {
  /**
   * INTERNAL: Do not reference the `Generated::MetatypeType` class directly.
   * Use the subclass `MetatypeType`, where the following predicates are available.
   */
  class MetatypeType extends Synth::TMetatypeType, AnyMetatypeTypeImpl::AnyMetatypeType {
    override string getAPrimaryQlClass() { result = "MetatypeType" }
  }
}
