package middleware

import (
	"bytes"
	"fmt"
	"gin-self/extend/utils/helpers"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"gin-self/extend/self_loger"
	"gin-self/extend/utils/e"

	"github.com/gin-gonic/gin"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// RecoveryFunc defines the function passable to CustomRecovery.
type RecoveryFunc func(c *gin.Context, err interface{})

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery() gin.HandlerFunc {
	return RecoveryWithWriter(gin.DefaultErrorWriter)
}

//CustomRecovery returns a middleware that recovers from any panics and calls the provided handle func to handle it.
func CustomRecovery(handle RecoveryFunc) gin.HandlerFunc {
	return RecoveryWithWriter(gin.DefaultErrorWriter, handle)
}

// RecoveryWithWriter returns a middleware for a given writer that recovers from any panics and writes a 500 if there was one.
func RecoveryWithWriter(out io.Writer, recovery ...RecoveryFunc) gin.HandlerFunc {
	if len(recovery) > 0 {
		return CustomRecoveryWithWriter(out, recovery[0])
	}
	return CustomRecoveryWithWriter(out, defaultHandleRecovery)
}

// CustomRecoveryWithWriter returns a middleware for a given writer that recovers from any panics and calls the provided handle func to handle it.
func CustomRecoveryWithWriter(out io.Writer, handle RecoveryFunc) gin.HandlerFunc {
	//var logger *log.Logger
	//if out != nil {
	//	logger = log.New(out, "\n\n\x1b[31m", log.LstdFlags)
	//}
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				stack := fmt.Sprintf("%s", err) + "\n\t" + string(stack(3))
				c.Request.Context().Value("trace").(*self_loger.TraceData).AddErrorStackLog(strings.Split(stack, "\n\t"))

				//if logger != nil {
				//	httpRequest, _ := httputil.DumpRequest(c.Request, false)
				//	headers := strings.Split(string(httpRequest), "\r\n")
				//	for idx, header := range headers {
				//		current := strings.Split(header, ":")
				//		if current[0] == "Authorization" {
				//			headers[idx] = current[0] + ": *"
				//		}
				//	}
				//	headersToStr := strings.Join(headers, "\r\n")
				//	if brokenPipe {
				//		logger.Printf("%s\n%s%s", err, headersToStr)
				//	} else if gin.IsDebugging() {
				//		logger.Printf("[Recovery] %s panic recovered:\n%s\n%s\n%s%s",
				//			timeFormat(time.Now()), headersToStr, err, stack)
				//	} else {
				//		logger.Printf("[Recovery] %s panic recovered:\n%s\n%s%s",
				//			timeFormat(time.Now()), err, stack)
				//	}
				//}

				if brokenPipe {
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
				} else {
					helpers.ApiError(c, e.ERROR, fmt.Sprintf("%s", err))
				}
			}
		}()
		c.Next()
	}
}

func defaultHandleRecovery(c *gin.Context, err interface{}) {
	c.AbortWithStatus(http.StatusInternalServerError)
}

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

func timeFormat(t time.Time) string {
	timeString := t.Format("2006/01/02 - 15:04:05")
	return timeString
}
