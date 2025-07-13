import { useEffect } from "react"
import { ToastContainer } from "react-toastify"
import { themeChange } from "theme-change"

import ErrorBoundary from "./components/ErrorBoundary/ErrorBoundary"
import MainLayout from "./layouts/MainLayout"

function App() {

  useEffect(() => {
    themeChange(false)
  }, [])
  
  return (
    <ErrorBoundary 
    >
      <MainLayout/>
      <ToastContainer
          position="top-right"
          autoClose={2000}
          newestOnTop
          draggable
          theme="dark"
          pauseOnFocusLoss={false}
      />
    </ErrorBoundary>
  )
}

export default App
