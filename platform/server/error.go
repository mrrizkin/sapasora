package server

import (
	"bufio"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"sapasora/platform/config"
	"sapasora/platform/satpam"
	"sapasora/platform/support/arr"
	"sapasora/platform/support/debug"
	"sapasora/platform/ui/inertia"

	"codeberg.org/mrrizkin/nihil"
	"github.com/gofiber/fiber/v2"
)

type ResError[T any] struct {
	Status  nihil.NilString `json:"status"           example:"error"`
	Message nihil.NilString `json:"message"          example:"error"`
	Detail  T               `json:"detail,omitempty"`
} //@name ResponseError

var errorTemplate = `<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8"/>
		<meta name="viewport" content="width=device-width, initial-scale=1"/>
		<title>Error {{.Code}}</title>
		<style>
			.code-line,body{line-height:1.5}:root{--purple:#7c3aed;--purple-light:#ede9fe;--gray-900:#111827;--gray-600:#4b5563;--gray-300:#d1d5db;--code-bg:#1a1a1a}body{font-family:ui-sans-serif,system-ui,-apple-system,BlinkMacSystemFont,sans-serif;background:#f7f7f7;color:var(--gray-900);margin:0;min-height:100vh;display:flex;justify-content:center;padding:1rem}.container{max-width:1000px;width:100%}.error-box{background:#fff;border:1px solid var(--gray-300);border-radius:.5rem;box-shadow:0 4px 6px -1px rgba(0,0,0,.1)}.error-header,.stack-item.active,.stack-item:hover{background:var(--purple);color:#fff}.error-header{padding:1.25rem 1.5rem;border-radius:.5rem .5rem 0 0;display:flex;align-items:center;gap:1rem}.error-header:not(:has(.error-header + .content)){border-radius:.5rem}.status-code{font-size:1.25rem;font-weight:600;padding:.25rem .75rem;background:rgba(255,255,255,.2);border-radius:.25rem}.error-message{font-weight:500;font-size:1.125rem}.code-line,.file-info,.stack-trace{font-size:.875rem}.content{padding:1.5rem}.stack-trace{list-style:none;padding:0;margin:0;font-family:ui-monospace,SFMono-Regular,Menlo,Monaco,Consolas,monospace}.stack-item{margin:.5rem 0;padding:1rem;background:var(--purple-light);border-radius:.375rem;cursor:pointer;transition:.2s}.stack-item.active .function-name,.stack-item:hover .function-name{color:#fff}.stack-item.active .file-info,.stack-item:hover .file-info{color:rgba(255,255,255,.8)}.stack-item.external{opacity:.7;background:#f3f4f6}.stack-item.external .function-name,.stack-item.external:hover .function-name{color:var(--gray-600)}.stack-item.external:hover{background:#e5e7eb;color:var(--gray-900)}.function-name{color:var(--purple);font-weight:600;display:block;margin-bottom:.25rem}.file-info{color:var(--gray-600)}.code-context{display:none;margin-top:1rem;background:var(--code-bg);border-radius:.375rem;overflow:hidden}.stack-item.active .code-context{display:block}.code-line{display:flex;min-width:100%}.line-number{color:#6b7280;padding:.125rem 1rem;text-align:right;user-select:none;border-right:1px solid #374151;background:rgba(0,0,0,.2)}.line-content{padding:.125rem 1rem;color:#e5e7eb;white-space:pre;flex:1}.current-line{background:rgba(124,58,237,.1)}.current-line .line-number{color:var(--purple-light);font-weight:700}.framework-label{text-transform:uppercase;letter-spacing:.05em;color:var(--gray-600);font-size:.75rem;margin-bottom:.75rem;font-weight:500}@media (max-width:640px){.error-header{flex-direction:column;align-items:flex-start;gap:.5rem}.content{padding:1rem}}
		</style>
	</head>
	<body>
		<div class="container">
			<div class="framework-label">Application Error</div>
			<div class="error-box">
				<div class="error-header">
					<div class="status-code">{{.Code}}</div>
					<div class="error-message">{{.Message}}</div>
				</div>
				{{if .Frames}}
				<div class="content">
					<ul class="stack-trace">
						{{range .Frames}}
						<li class="stack-item{{if not .IsInternal}} external{{end}}">
							<span class="function-name">{{.Frame.Function}}</span>
							<span class="file-info">{{.Frame.File}}:{{.Frame.Line}}</span>
							{{if and .CodeLines .IsInternal}}
							<div class="code-context">
								{{range .CodeLines}}
								<div class="code-line{{if .IsCurrent}} current-line{{end}}">
									<span class="line-number">{{.Number}}</span>
									<span class="line-content">{{.Content}}</span>
								</div>
								{{end}}
							</div>
							{{end}}
						</li>
						{{end}}
					</ul>
				</div>
				{{end}}
			</div>
		</div>
		<script>
			document.addEventListener("DOMContentLoaded",()=>{let e=document.querySelector(".stack-item:not(.external)")||document.querySelector(".stack-item");e&&e.classList.add("active");let t=document.querySelectorAll(".stack-item");for(let l=t.length-1;l>=0;l--){let a=t[l];a.addEventListener("click",function(e){for(let l=t.length-1;l>=0;l--)t[l]?.classList?.remove("active");let a=e.currentTarget;a?.classList?.add("active")})}});
		</script>
	</body>
</html>`

var tmpl = template.Must(template.New("error").Parse(errorTemplate))

// DefaultErrorHandler is a middleware that handles errors returned by the application.
func DefaultErrorHandler(c config.Config) fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
		}

		var authorizationError *satpam.ForbiddenError
		if errors.As(err, &authorizationError) {
			code = fiber.StatusForbidden
		}

		var stackFrames []debug.StackFrame
		if isDevelopment(c) {
			if stack, ok := ctx.Locals("stack_trace").([]debug.StackFrame); ok {
				stackFrames = stack
			} else if stack, err := debug.StackTrace(); err == nil {
				stackFrames = stack
			}
		}

		if (ctx.Get("X-Requested-With") != "XMLHttpRequest" &&
			ctx.Get("Accept") != "application/json") ||
			inertia.IsInertiaRequest(ctx) {
			frames := make([]debug.StackFrameContext, 0)

			if len(stackFrames) > 0 {
				cwd, _ := os.Getwd()
				for _, frame := range stackFrames {
					frames = append(frames, debug.StackFrameContext{
						Frame:      frame,
						CodeLines:  getFileContext(frame.File, frame.Line, 10),
						IsInternal: isInternalFrame(cwd, frame),
					})
				}
			}

			error := struct {
				Code    int
				Message string
				Frames  []debug.StackFrameContext
			}{
				Code:    code,
				Message: err.Error(),
				Frames:  frames,
			}
			ctx.Set("Content-Type", "text/html; charset=utf-8")
			ctx.Status(code)
			return tmpl.Execute(ctx.Response().BodyWriter(), error)
		}

		detail := arr.Map(stackFrames, func(frame debug.StackFrame) string {
			return fmt.Sprintf(
				"%s (%s:%d)",
				frame.Function,
				frame.File,
				frame.Line,
			)
		})

		errMessage := err.Error()
		response := ResError[[]string]{
			Status:  nihil.String(http.StatusText(code)),
			Message: nihil.String(errMessage),
			Detail:  detail,
		}

		return ctx.Status(code).JSON(response)
	}
}

// isInternalFrame checks if the given stack frame is internal to the application.
func isInternalFrame(cwd string, frame debug.StackFrame) bool {
	return strings.HasPrefix(frame.File, cwd)
}

// getFileContext returns the lines of code around the given line number.
func getFileContext(filename string, targetLine int, contextLines int) []debug.CodeLine {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []debug.CodeLine{}
	lineNum := 1
	start := targetLine - contextLines
	end := targetLine + contextLines

	for scanner.Scan() {
		if lineNum >= start && lineNum <= end {
			lines = append(lines, debug.CodeLine{
				Number:    lineNum,
				Content:   scanner.Text(),
				IsCurrent: lineNum == targetLine,
			})
		} else if lineNum > end {
			break
		}
		lineNum++
	}

	return lines
}

func isDevelopment(cfg config.Config) bool {
	env := cfg.GetString("app.env", "development")
	return env == "development" || env == "dev"
}
