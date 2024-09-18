// generated by codegen, do not edit
/**
 * This module provides the generated definition of `RecordPat`.
 * INTERNAL: Do not import directly.
 */

private import codeql.rust.internal.generated.Synth
private import codeql.rust.internal.generated.Raw
import codeql.rust.elements.internal.PatImpl::Impl as PatImpl
import codeql.rust.elements.Path
import codeql.rust.elements.RecordPatField

/**
 * INTERNAL: This module contains the fully generated definition of `RecordPat` and should not
 * be referenced directly.
 */
module Generated {
  /**
   * A record pattern. For example:
   * ```
   * match x {
   *     Foo { a: 1, b: 2 } => "ok",
   *     Foo { .. } => "fail",
   * }
   * ```
   * INTERNAL: Do not reference the `Generated::RecordPat` class directly.
   * Use the subclass `RecordPat`, where the following predicates are available.
   */
  class RecordPat extends Synth::TRecordPat, PatImpl::Pat {
    override string getAPrimaryQlClass() { result = "RecordPat" }

    /**
     * Gets the path of this record pat, if it exists.
     */
    Path getPath() {
      result =
        Synth::convertPathFromRaw(Synth::convertRecordPatToRaw(this).(Raw::RecordPat).getPath())
    }

    /**
     * Holds if `getPath()` exists.
     */
    final predicate hasPath() { exists(this.getPath()) }

    /**
     * Gets the `index`th fld of this record pat (0-based).
     */
    RecordPatField getFld(int index) {
      result =
        Synth::convertRecordPatFieldFromRaw(Synth::convertRecordPatToRaw(this)
              .(Raw::RecordPat)
              .getFld(index))
    }

    /**
     * Gets any of the flds of this record pat.
     */
    final RecordPatField getAFld() { result = this.getFld(_) }

    /**
     * Gets the number of flds of this record pat.
     */
    final int getNumberOfFlds() { result = count(int i | exists(this.getFld(i))) }

    /**
     * Holds if this record pat has ellipsis.
     */
    predicate hasEllipsis() { Synth::convertRecordPatToRaw(this).(Raw::RecordPat).hasEllipsis() }
  }
}
