#pragma once

#include <fstream>
#include <iostream>
#include <regex>
#include <vector>

#include <binlog/binlog.hpp>
#include <binlog/TextOutputStream.hpp>
#include <binlog/EventFilter.hpp>
#include <binlog/adapt_stdfilesystem.hpp>
#include <binlog/adapt_stderrorcode.hpp>
#include <binlog/adapt_stdoptional.hpp>
#include <binlog/adapt_stdvariant.hpp>

#include "swift/logging/SwiftDiagnostics.h"

// Logging macros. These will call `logger()` to get a Logger instance, picking up any `logger`
// defined in the current scope. Domain-specific loggers can be added or used by either:
// * providing a class field called `logger` (as `Logger::operator()()` returns itself)
// * declaring a local `logger` variable (to be used for one-time execution like code in `main`)
// * declaring a `Logger& logger()` function returning a reference to a static local variable
// * passing a logger around using a `Logger& logger` function parameter
// They are created with a name that appears in the logs and can be used to filter debug levels (see
// `Logger`).
// If the first argument after the format is a SwiftDiagnosticSource or
// SwiftDiagnosticSourceWithLocation, a JSON diagnostic entry is emitted. In this case the
// format string **must** start with "[{}] " (which is checked at compile time), and everything
// following that is used to form the message in the diagnostics using fmt::format instead of the
// internal binlog formatting. The two are fairly compatible though.
#define LOG_CRITICAL(...) LOG_WITH_LEVEL(critical, __VA_ARGS__)
#define LOG_ERROR(...) LOG_WITH_LEVEL(error, __VA_ARGS__)
#define LOG_WARNING(...) LOG_WITH_LEVEL(warning, __VA_ARGS__)
#define LOG_INFO(...) LOG_WITH_LEVEL(info, __VA_ARGS__)
#define LOG_DEBUG(...) LOG_WITH_LEVEL(debug, __VA_ARGS__)
#define LOG_TRACE(...) LOG_WITH_LEVEL(trace, __VA_ARGS__)

#define CODEQL_GET_SECOND(...) CODEQL_GET_SECOND_I(__VA_ARGS__, 0, 0)
#define CODEQL_GET_SECOND_I(X, Y, ...) Y

// only do the actual logging if the picked up `Logger` instance is configured to handle the
// provided log level. `LEVEL` must be a compile-time constant. `logger()` is evaluated once
// TODO(C++20) replace non-standard ##__VA_ARGS__ with __VA_OPT__(,) __VA_ARGS__
#define LOG_WITH_LEVEL_AND_CATEGORY(LEVEL, CATEGORY, ...)                                      \
  do {                                                                                         \
    static_assert(::codeql::detail::checkLogArgs<decltype(CODEQL_GET_SECOND(__VA_ARGS__))>(    \
                      MSERIALIZE_FIRST(__VA_ARGS__)),                                          \
                  "diagnostics logs must have format starting with \"[{}]\"");                 \
    constexpr auto _level = ::codeql::Log::Level::LEVEL;                                       \
    ::codeql::Logger& _logger = logger();                                                      \
    if (_level >= _logger.level()) {                                                           \
      BINLOG_CREATE_SOURCE_AND_EVENT(_logger.writer(), _level, CATEGORY, ::binlog::clockNow(), \
                                     __VA_ARGS__);                                             \
    }                                                                                          \
    if (_level >= ::codeql::Log::Level::error) {                                               \
      ::codeql::Log::flush();                                                                  \
    }                                                                                          \
  } while (false)

#define LOG_WITH_LEVEL(LEVEL, ...) LOG_WITH_LEVEL_AND_CATEGORY(LEVEL, , __VA_ARGS__)

// avoid calling into binlog's original macros
#undef BINLOG_CRITICAL
#undef BINLOG_CRITICAL_W
#undef BINLOG_CRITICAL_C
#undef BINLOG_CRITICAL_WC
#undef BINLOG_ERROR
#undef BINLOG_ERROR_W
#undef BINLOG_ERROR_C
#undef BINLOG_ERROR_WC
#undef BINLOG_WARNING
#undef BINLOG_WARNING_W
#undef BINLOG_WARNING_C
#undef BINLOG_WARNING_WC
#undef BINLOG_INFO
#undef BINLOG_INFO_W
#undef BINLOG_INFO_C
#undef BINLOG_INFO_WC
#undef BINLOG_DEBUG
#undef BINLOG_DEBUG_W
#undef BINLOG_DEBUG_C
#undef BINLOG_DEBUG_WC
#undef BINLOG_TRACE
#undef BINLOG_TRACE_W
#undef BINLOG_TRACE_C
#undef BINLOG_TRACE_WC

namespace codeql {

// tools should define this to tweak the root name of all loggers
extern const std::string_view programName;

// This class is responsible for the global log state (outputs, log level rules, flushing)
// State is stored in the singleton `Log::instance()`.
// Before using logging, `Log::configure("<name>")` should be used (e.g.
// `Log::configure("extractor")`). Then, `Log::flush()` should be regularly called.
// Logging is configured upon first usage.  This consists in
//  * using environment variable `CODEQL_EXTRACTOR_SWIFT_LOG_DIR` to choose where to dump the log
//    file(s). Log files will go to a subdirectory thereof named after `programName`
//  * using environment variable `CODEQL_EXTRACTOR_SWIFT_LOG_LEVELS` to configure levels for
//    loggers and outputs. This must have the form of a comma separated `spec:level` list, where
//    `spec` is either a glob pattern (made up of alphanumeric, `/`, `*` and `.` characters) for
//    matching logger names or one of `out:binary`, `out:text`, `out:console` or `out:diagnostics`.
//    Output default levels can be seen in the corresponding initializers below. By default, all
//    loggers are configured with the lowest output level
class Log {
 public:
  using Level = binlog::Severity;

  // Internal data required to build `Logger` instances
  struct LoggerConfiguration {
    binlog::Session& session;
    std::string fullyQualifiedName;
    Level level;
  };

  // Flush logs to the designated outputs
  static void flush() {
    if (initialized) {
      instance().flushImpl();
    }
  }

  // create `Logger` configuration, used internally by `Logger`'s constructor
  static LoggerConfiguration getLoggerConfiguration(std::string_view name) {
    return instance().getLoggerConfigurationImpl(name);
  }

  template <typename Source>
  static void diagnose(const Source& source,
                       const std::chrono::system_clock::time_point& time,
                       std::string_view message) {
    instance().diagnostics.write(source, time, message);
  }

 private:
  static constexpr const char* format = "%u %S [%n] %m (%G:%L)\n";
  static bool initialized;

  Log() { configure(); }

  static Log& instance() {
    static Log ret;
    return ret;
  }

  class Logger& logger();

  void configure();
  void flushImpl();

  LoggerConfiguration getLoggerConfigurationImpl(std::string_view name);

  // make `session.consume(*this)` work, which requires access to `write`
  friend binlog::Session;
  Log& write(const char* buffer, std::streamsize size);

  // Output filtered according to a configured log level
  template <typename Output>
  struct FilteredOutput {
    binlog::Severity level;
    Output output;
    binlog::EventFilter filter;

    template <typename... Args>
    FilteredOutput(Level level, Args&&... args)
        : level{level}, output{std::forward<Args>(args)...}, filter{filterOnLevel()} {}

    FilteredOutput& write(const char* buffer, std::streamsize size) {
      filter.writeAllowed(buffer, size, output);
      return *this;
    }

    binlog::EventFilter::Predicate filterOnLevel() const {
      return [this](const binlog::EventSource& src) { return src.severity >= level; };
    }

    // if configured as `no_logs`, the output is effectively disabled
    explicit operator bool() const { return level < Level::no_logs; }
  };

  using LevelRule = std::pair<std::regex, Level>;
  using LevelRules = std::vector<LevelRule>;

  binlog::Session session;
  std::ofstream textFile;
  FilteredOutput<std::ofstream> binary{Level::no_logs};
  FilteredOutput<binlog::TextOutputStream> text{Level::info, textFile, format};
  FilteredOutput<binlog::TextOutputStream> console{Level::warning, std::cerr, format};
  SwiftDiagnosticsDumper diagnostics{};
  LevelRules sourceRules;
  std::vector<std::string> collectLevelRulesAndReturnProblems(const char* envVar);
};

// This class represent a named domain-specific logger, responsible for pushing logs using the
// underlying `binlog::SessionWriter` class. This has a configured log level, so that logs on this
// `Logger` with a level lower than the configured one are no-ops. The level is configured based
// on rules matching `<programName>/<name>` in `CODEQL_EXTRACTOR_SWIFT_LOG_LEVELS` (see above).
// `<name>` is provided in the constructor. If no rule matches the name, the log level defaults to
// the minimum level of all outputs.
class Logger {
 public:
  // configured logger based on name, as explained above
  explicit Logger(std::string_view name) : Logger(Log::getLoggerConfiguration(name)) {}

  // used internally, public to be accessible to Log for its own logger
  explicit Logger(Log::LoggerConfiguration&& configuration)
      : w{configuration.session, queueSize, /* id */ 0,
          std::move(configuration.fullyQualifiedName)},
        level_{configuration.level} {}

  binlog::SessionWriter& writer() { return w; }
  Log::Level level() const { return level_; }

  // make defining a `Logger logger` field be equivalent to providing a `Logger& logger()` function
  // in order to be picked up by logging macros
  Logger& operator()() { return *this; }

 private:
  static constexpr size_t queueSize = 1 << 20;  // default taken from binlog

  binlog::SessionWriter w;
  Log::Level level_;
};

namespace detail {
constexpr std::string_view diagnosticsFormatPrefix = "[{}] ";

template <typename T>
constexpr bool checkLogArgs(std::string_view format) {
  using Type = std::remove_cv_t<std::remove_reference_t<T>>;
  constexpr bool isDiagnostic = std::is_same_v<Type, SwiftDiagnosticsSource> ||
                                std::is_same_v<Type, SwiftDiagnosticsSourceWithLocation>;
  return !isDiagnostic ||
         format.substr(0, diagnosticsFormatPrefix.size()) == diagnosticsFormatPrefix;
}

template <typename Writer, typename Source, typename... T>
void binlogAddEventIgnoreFirstOverload(Writer& writer,
                                       std::uint64_t eventSourceId,
                                       std::uint64_t clock,
                                       const char* format,
                                       const Source& source,
                                       T&&... t) {
  std::chrono::system_clock::time_point point{
      std::chrono::duration_cast<std::chrono::system_clock::duration>(
          std::chrono::nanoseconds{clock})};
  constexpr auto offset = ::codeql::detail::diagnosticsFormatPrefix.size();
  ::codeql::Log::diagnose(source, point, fmt::format(format + offset, t...));
  writer.addEvent(eventSourceId, clock, source, std::forward<T>(t)...);
}

}  // namespace detail
}  // namespace codeql

// we intercept this binlog plumbing function providing better overload resolution matches in
// case the first non-format argument is a diagnostic source, and emit it in that case with the
// same timestamp
namespace binlog::detail {
template <typename Writer, size_t N, typename... T>
void addEventIgnoreFirst(Writer& writer,
                         std::uint64_t eventSourceId,
                         std::uint64_t clock,
                         const char (&format)[N],
                         const codeql::SwiftDiagnosticsSource& source,
                         T&&... t) {
  codeql::detail::binlogAddEventIgnoreFirstOverload(writer, eventSourceId, clock, format, source,
                                                    std::forward<T>(t)...);
}

template <typename Writer, size_t N, typename... T>
void addEventIgnoreFirst(Writer& writer,
                         std::uint64_t eventSourceId,
                         std::uint64_t clock,
                         const char (&format)[N],
                         codeql::SwiftDiagnosticsSourceWithLocation&& source,
                         T&&... t) {
  codeql::detail::binlogAddEventIgnoreFirstOverload(writer, eventSourceId, clock, format, source,
                                                    std::forward<T>(t)...);
}

template <typename Writer, size_t N, typename... T>
void addEventIgnoreFirst(Writer& writer,
                         std::uint64_t eventSourceId,
                         std::uint64_t clock,
                         const char (&format)[N],
                         const codeql::SwiftDiagnosticsSourceWithLocation& source,
                         T&&... t) {
  codeql::detail::binlogAddEventIgnoreFirstOverload(writer, eventSourceId, clock, format, source,
                                                    std::forward<T>(t)...);
}

}  // namespace binlog::detail
