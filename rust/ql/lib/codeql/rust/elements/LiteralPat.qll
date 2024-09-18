// generated by codegen, do not edit
/**
 * This module provides the public class `LiteralPat`.
 */

private import internal.LiteralPatImpl
import codeql.rust.elements.Expr
import codeql.rust.elements.Pat

/**
 * A literal pattern. For example:
 * ```
 * match x {
 *     42 => "ok",
 *     _ => "fail",
 * }
 * ```
 */
final class LiteralPat = Impl::LiteralPat;
