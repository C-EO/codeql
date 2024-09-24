import go
import semmle.go.dataflow.ExternalFlow
import ModelValidation
import semmle.go.dataflow.internal.FlowSummaryImpl as FlowSummaryImpl
import TestUtilities.InlineExpectationsTest
import MakeTest<FlowTest>

module Config implements DataFlow::ConfigSig {
  predicate isSource(DataFlow::Node source) { sources(source) }

  predicate isSink(DataFlow::Node sink) { sinks(sink) }
}

module Flow = TaintTracking::Global<Config>;

module FlowTest implements TestSig {
  string getARelevantTag() { result = "S1[t]" }

  predicate hasActualResult(Location location, string element, string tag, string value) {
    tag = "S1[t]" and
    exists(DataFlow::Node sink | Flow::flowTo(sink) |
      sink.hasLocationInfo(location.getFile().getAbsolutePath(), location.getStartLine(),
        location.getStartColumn(), location.getEndLine(), location.getEndColumn()) and
      element = sink.toString() and
      value = ""
    )
  }
}

query predicate sources(DataFlow::Node source) { source instanceof RemoteFlowSource }

query predicate sinks(DataFlow::Node sink) { sink = any(FileSystemAccess fsa).getAPathArgument() }
