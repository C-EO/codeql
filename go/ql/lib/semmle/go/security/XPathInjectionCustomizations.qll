/**
 * Provides default sources, sinks and sanitizers for reasoning about untrusted user input used in an XPath expression,
 * as well as extension points for adding your own.
 */

import go

/**
 * Provides extension points for reasoning about untrusted user input used in an XPath expression.
 */
module XPathInjection {
  /**
   * A data flow source for untrusted user input used in an XPath expression.
   */
  abstract class Source extends DataFlow::Node { }

  /**
   * A data flow sink for untrusted user input used in an XPath expression.
   */
  abstract class Sink extends DataFlow::ExprNode { }

  /**
   * A sanitizer for untrusted user input used in an XPath expression.
   */
  abstract class Sanitizer extends DataFlow::ExprNode { }

  /** A source of untrusted data, used in an XPath expression. */
  class RemoteFlowAsSource extends Source instanceof RemoteFlowSource { }

  /** An XPath expression string, considered as a taint sink for XPath injection. */
  class XPathExpressionStringAsSink extends Sink instanceof XPath::XPathExpressionString { }
}
