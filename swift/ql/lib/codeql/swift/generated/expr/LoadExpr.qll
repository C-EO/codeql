// generated by codegen/codegen.py, do not edit
/**
 * This module provides the generated definition of `LoadExpr`.
 * INTERNAL: Do not import directly.
 */

private import codeql.swift.generated.Synth
private import codeql.swift.generated.Raw
import codeql.swift.elements.expr.internal.ImplicitConversionExprImpl::Impl as ImplicitConversionExprImpl

/**
 * INTERNAL: This module contains the fully generated definition of `LoadExpr` and should not
 * be referenced directly.
 */
module Generated {
  /**
   * INTERNAL: Do not reference the `Generated::LoadExpr` class directly.
   * Use the subclass `LoadExpr`, where the following predicates are available.
   */
  class LoadExpr extends Synth::TLoadExpr, ImplicitConversionExprImpl::ImplicitConversionExpr {
    override string getAPrimaryQlClass() { result = "LoadExpr" }
  }
}
