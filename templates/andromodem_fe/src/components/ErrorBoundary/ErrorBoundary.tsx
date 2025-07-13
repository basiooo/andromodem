import { Component, type ErrorInfo,type ReactNode } from "react"
import { MdError, MdRefresh } from "react-icons/md"

interface ErrorFallbackProps {
  error: Error
  onReset: () => void
}

const ErrorFallback = ({ error, onReset }: ErrorFallbackProps) => {
  return (
    <div className="min-h-screen flex items-center justify-center bg-base-200 p-4">
      <div className="card w-full max-w-md bg-base-100 shadow-xl">
        <div className="card-body text-center">
          <div className="flex justify-center mb-4">
            <div className="bg-error/20 rounded-full p-4">
              <MdError className="text-4xl text-error" />
            </div>
          </div>
          
          <h2 className="card-title justify-center text-error mb-2">
            Oops! Something went wrong
          </h2>
          
          <p className="text-base-content/70 mb-4">
            An unexpected error occurred. Please try refreshing the page or contact support if the problem persists.
          </p>
          
          <div className="collapse collapse-arrow bg-base-200 mb-4">
            <input type="checkbox" className="peer" />
            <div className="collapse-title text-sm font-medium">
              Show Error Details
            </div>
            <div className="collapse-content">
              <div className="bg-base-300 rounded-lg p-3 text-left">
                <p className="text-xs font-mono text-error break-all">
                  <strong>Error:</strong> {error.message}
                </p>
                {error.stack && (
                  <details className="mt-2">
                    <summary className="text-xs cursor-pointer text-base-content/70">
                      Stack Trace
                    </summary>
                    <pre className="text-xs mt-2 whitespace-pre-wrap break-all text-base-content/60">
                      {error.stack}
                    </pre>
                  </details>
                )}
              </div>
            </div>
          </div>
          
          <div className="card-actions justify-center gap-2">
            <button 
              className="btn btn-primary btn-sm"
              onClick={onReset}
            >
              <MdRefresh className="mr-1" />
              Try Again
            </button>
            <button 
              className="btn btn-outline btn-sm"
              onClick={() => window.location.reload()}
            >
              Reload Page
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}

interface ErrorBoundaryProps {
  children: ReactNode
  onError?: (error: Error, errorInfo: ErrorInfo) => void
}

interface ErrorBoundaryState {
  hasError: boolean
  error: Error | null
}

class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: ErrorBoundaryProps) {
    super(props)
    this.state = { hasError: false, error: null }
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('ErrorBoundary caught an error:', error, errorInfo)
    
    if (this.props.onError) {
      this.props.onError(error, errorInfo)
    }
  }

  handleReset = () => {
    this.setState({ hasError: false, error: null })
  }

  render() {
    if (this.state.hasError && this.state.error) {
      return <ErrorFallback error={this.state.error} onReset={this.handleReset} />
    }

    return this.props.children
  }
}

export default ErrorBoundary