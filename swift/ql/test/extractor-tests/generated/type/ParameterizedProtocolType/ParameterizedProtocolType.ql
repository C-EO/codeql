// generated by codegen/codegen.py, do not edit
import codeql.swift.elements
import TestUtils

from
  ParameterizedProtocolType x, string getName, Type getCanonicalType, ProtocolType getBase,
  int getNumberOfArgs
where
  toBeTested(x) and
  not x.isUnknown() and
  getName = x.getName() and
  getCanonicalType = x.getCanonicalType() and
  getBase = x.getBase() and
  getNumberOfArgs = x.getNumberOfArgs()
select x, "getName:", getName, "getCanonicalType:", getCanonicalType, "getBase:", getBase,
  "getNumberOfArgs:", getNumberOfArgs
