// generated by codegen, do not edit
/**
 * This module provides the public class `ContinueExpr`.
 */

private import internal.ContinueExprImpl
import codeql.rust.elements.Expr
import codeql.rust.elements.Label

/**
 * A continue expression. For example:
 * ```
 * loop {
 *     if not_ready() {
 *         continue;
 *     }
 * }
 * ```
 * ```
 * 'label: loop {
 *     if not_ready() {
 *         continue 'label;
 *     }
 * }
 * ```
 */
final class ContinueExpr = Impl::ContinueExpr;
