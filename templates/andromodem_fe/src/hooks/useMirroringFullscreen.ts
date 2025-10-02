import { useEffect,useRef, useState } from 'react'

interface useMirroringFullscreenReturn {
  isFullscreen: boolean;
  toggleFullscreen: () => Promise<void>;
  exitFullscreen: () => Promise<void>;
  enterFullscreen: () => Promise<void>;
  setFullscreenElement: (element: HTMLElement | null) => void;
}

export const useMirroringFullscreen = (): useMirroringFullscreenReturn => {
  const [isFullscreen, setIsFullscreen] = useState(false)
  const elementRef = useRef<HTMLElement | null>(null)
  useEffect(() => {
    const handleFullscreenChange = () => {
      setIsFullscreen(!!document.fullscreenElement)
    }

    document.addEventListener('fullscreenchange', handleFullscreenChange)
    return () => {
      document.removeEventListener('fullscreenchange', handleFullscreenChange)
    }
  }, [])

  const enterFullscreen = async (): Promise<void> => {
    if (!elementRef.current) return
    
    try {
      if (elementRef.current.requestFullscreen) {
        await elementRef.current.requestFullscreen()
      }
    } catch (error) {
      console.error('Error entering fullscreen:', error)
    }
  }

  const exitFullscreen = async (): Promise<void> => {
    try {
      if (document.fullscreenElement && document.exitFullscreen) {
        await document.exitFullscreen()
      }
    } catch (error) {
      console.error('Error exiting fullscreen:', error)
    }
  }

  const toggleFullscreen = async (): Promise<void> => {
    if (isFullscreen) {
      await exitFullscreen()
    } else {
      await enterFullscreen()
    }
  }

  const setFullscreenElement = (element: HTMLElement | null) => {
    elementRef.current = element
  }

  return {
    isFullscreen,
    toggleFullscreen,
    exitFullscreen,
    enterFullscreen,
    setFullscreenElement
  }
}