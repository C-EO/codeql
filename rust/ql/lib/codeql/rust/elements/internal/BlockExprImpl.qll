// generated by codegen, remove this comment if you wish to edit this file
/**
 * This module provides a hand-modifiable wrapper around the generated class `BlockExpr`.
 *
 * INTERNAL: Do not use.
 */

private import codeql.rust.elements.internal.generated.BlockExpr

/**
 * INTERNAL: This module contains the customizable definition of `BlockExpr` and should not
 * be referenced directly.
 */
module Impl {
  /**
   * A block expression. For example:
   * ```rust
   * {
   *     let x = 42;
   * }
   * ```
   * ```rust
   * 'label: {
   *     let x = 42;
   *     x
   * }
   * ```
   */
  class BlockExpr extends Generated::BlockExpr { }
}
