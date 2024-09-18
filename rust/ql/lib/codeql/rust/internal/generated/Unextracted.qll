// generated by codegen, do not edit
/**
 * This module provides the generated definition of `Unextracted`.
 * INTERNAL: Do not import directly.
 */

private import codeql.rust.internal.generated.Synth
private import codeql.rust.internal.generated.Raw
import codeql.rust.elements.internal.ElementImpl::Impl as ElementImpl

/**
 * INTERNAL: This module contains the fully generated definition of `Unextracted` and should not
 * be referenced directly.
 */
module Generated {
  /**
   * The base class marking everything that was not properly extracted for some reason, such as:
   * * syntax errors
   * * insufficient context information
   * * yet unimplemented parts of the extractor
   * INTERNAL: Do not reference the `Generated::Unextracted` class directly.
   * Use the subclass `Unextracted`, where the following predicates are available.
   */
  class Unextracted extends Synth::TUnextracted, ElementImpl::Element { }
}
