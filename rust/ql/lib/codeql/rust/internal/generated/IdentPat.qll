// generated by codegen, do not edit
/**
 * This module provides the generated definition of `IdentPat`.
 * INTERNAL: Do not import directly.
 */

private import codeql.rust.internal.generated.Synth
private import codeql.rust.internal.generated.Raw
import codeql.rust.elements.Pat
import codeql.rust.elements.internal.PatImpl::Impl as PatImpl

/**
 * INTERNAL: This module contains the fully generated definition of `IdentPat` and should not
 * be referenced directly.
 */
module Generated {
  /**
   * A binding pattern. For example:
   * ```
   * match x {
   *     Option::Some(y) => y,
   *     Option::None => 0,
   * };
   * ```
   * ```
   * match x {
   *     y@Option::Some(_) => y,
   *     Option::None => 0,
   * };
   * ```
   * INTERNAL: Do not reference the `Generated::IdentPat` class directly.
   * Use the subclass `IdentPat`, where the following predicates are available.
   */
  class IdentPat extends Synth::TIdentPat, PatImpl::Pat {
    override string getAPrimaryQlClass() { result = "IdentPat" }

    /**
     * Gets the binding of this ident pat.
     */
    string getBindingId() {
      result = Synth::convertIdentPatToRaw(this).(Raw::IdentPat).getBindingId()
    }

    /**
     * Gets the subpat of this ident pat, if it exists.
     */
    Pat getSubpat() {
      result =
        Synth::convertPatFromRaw(Synth::convertIdentPatToRaw(this).(Raw::IdentPat).getSubpat())
    }

    /**
     * Holds if `getSubpat()` exists.
     */
    final predicate hasSubpat() { exists(this.getSubpat()) }
  }
}
